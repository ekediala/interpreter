package ast

import "github.com/ekediala/interpreter/token"

type ReturnStatement struct {
	Token       token.Token // token.Return
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
