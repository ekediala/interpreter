package parser

import (
	"fmt"

	"github.com/ekediala/interpreter/ast"
	"github.com/ekediala/interpreter/lexer"
	"github.com/ekediala/interpreter/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	nextToken    token.Token
	errors       []string
}

func New(lexer *lexer.Lexer) *Parser {
	p := Parser{
		lexer:  lexer,
		errors: make([]string, 0, 20),
	}

	// read tokens twice so that currentToken and nextToken are set correctly
	p.next()
	p.next()
	return &p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.nextToken.Type)
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
		return nil
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
