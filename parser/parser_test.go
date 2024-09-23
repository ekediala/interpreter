package parser_test

import (
	"strings"
	"testing"

	"github.com/ekediala/interpreter/ast"
	"github.com/ekediala/interpreter/lexer"
	"github.com/ekediala/interpreter/parser"
	"github.com/ekediala/interpreter/token"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements; got %d", program.Statements)
	}

	tests := []struct{ expectedIdentifier string }{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements; got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement; got %T", stmt)
		}

		expectedTokenLiteral := strings.ToLower(token.RETURN)

		if returnStmt.TokenLiteral() != expectedTokenLiteral {
			t.Errorf("returnStmt.TokenLiteral not %q; got %q", expectedTokenLiteral, returnStmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, variableName string) {
	t.Helper()

	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral() not 'let'; got: %q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement; got %T", s)
	}

	if letStmt.Name.Value != variableName {
		t.Fatalf("letStmt.Name.Value not %q; got: %q", variableName, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != variableName {
		t.Fatalf("letStmt.Name.TokenLiteral() not %q; got %q", variableName, letStmt.Name.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
