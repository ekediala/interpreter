package ast

import "github.com/ekediala/interpreter/token"

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

func (i *Identifier) String() string {
	return i.Value
}