package ast

import "github.com/ekediala/interpreter/token"

type LetStatement struct {
	Token token.Token // token.let
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) statementNode() {}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	Token token.Token // token.identifier
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// sometimes identifiers do produce values
// say we have let x = 5; let y = x; the identifier x here does produce a value
func (i *Identifier) expressionNode() {}
