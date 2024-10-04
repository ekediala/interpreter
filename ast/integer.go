package ast

import "github.com/ekediala/interpreter/token"

type IntegerLiteral struct {
	Value int64
	Token token.Token // token.Integer
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}
