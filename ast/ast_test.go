package ast_test

import (
	"testing"

	"github.com/ekediala/interpreter/ast"
	"github.com/ekediala/interpreter/token"
)

func TestString(t *testing.T) {
	source := `let myVar = anotherVar;`

	program := &ast.RootNode{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != source {
		t.Errorf("program.String() wrong; expected %q; got %q", source, program.String())
	}
}
