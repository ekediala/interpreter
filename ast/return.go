package ast

import (
	"strings"

	"github.com/ekediala/interpreter/token"
)

type ReturnStatement struct {
	Token       token.Token // token.Return
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
