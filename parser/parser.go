package parser

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ekediala/interpreter/ast"
	"github.com/ekediala/interpreter/lexer"
	"github.com/ekediala/interpreter/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // < or >
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
	errors       []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := Parser{
		lexer:          lexer,
		errors:         make([]string, 0, 20),
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefixFn(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)

	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)

	// read tokens twice so that currentToken and nextToken are set correctly
	p.next()
	p.next()
	return &p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %q, got %s instead", t, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) next() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.ReadAndAdvanceToken()
}

func (p *Parser) ParseProgram() *ast.RootNode {
	program := ast.RootNode{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.next()
	}
	return &program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.ReturnStatement{Token: p.currentToken}
	p.next()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.next()
	}

	return &stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := ast.LetStatement{
		Token: p.currentToken,
	}

	// let statements must be followed by an identifier
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	// advance tokens to parse the identifier
	p.next()
	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// advance token to the assignment operator
	p.next()
	// we do not need to parse the assignment operator, so advance tokens again
	p.next()

	// skip parsing expressions
	for !p.currentTokenIs(token.SEMICOLON) {
		p.next()
	}

	return &stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))

	stmt := ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextTokenIs(token.SEMICOLON) {
		p.next()
	}

	return &stmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))

	prefixExp := ast.PrefixExpression{
		Operator: p.currentToken.Literal,
		Token:    p.currentToken,
	}

	p.next()

	prefixExp.Right = p.parseExpression(PREFIX)
	return &prefixExp
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.addNoPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.nextTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.nextToken.Type]
		if infixFn == nil {
			return leftExp
		}

		p.next()

		leftExp = infixFn(leftExp)
	}

	return leftExp
}

// tells us if the next token is the same type as t
func (p *Parser) expectPeek(t token.TokenType) bool {
	res := p.nextToken.Type == t
	if !res {
		p.peekError(t)
	}
	return res
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) parseIdentifier() ast.Expression {
	exp := ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	return &exp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("%q could not be parsed into int64", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	exp := ast.IntegerLiteral{
		Value: value,
		Token: p.currentToken,
	}

	return &exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))

	exp := ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.next()
	exp.Right = p.parseExpression(precedence)

	return &exp
}

func (p *Parser) parseBoolean() ast.Expression {
	exp := ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
	return &exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	slog.Info("parse grouped exp", "curr token", p.currentToken)
	p.next()

	exp := p.parseExpression(LOWEST)

	slog.Info("parse grouped exp", "exp", exp, "expect peek", p.nextToken)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	p.next()

	return exp
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) addNoPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}
