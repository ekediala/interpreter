package ast

import "strings"

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// This node  is going to be the root node of every AST our parser produces.
type RootNode struct {
	Statements []Statement
}

func (r *RootNode) TokenLiteral() string {
	if len(r.Statements) > 0 {
		return r.Statements[0].TokenLiteral()
	}

	return ""
}

func (r *RootNode) String() string {
	var out strings.Builder

	for _, stmt := range r.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}
