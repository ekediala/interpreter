package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// illegal is a character we do not understand
	ILLEGAL = "ILLEGAL"
	// we will use this to signal a stop to our parser
	EOF = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER" // variable names
	INT        = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	BANG     = "!"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "true"
	FALSE    = "false"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdentifier(identifier string) TokenType {
	if tokType, ok := keywords[identifier]; ok {
		return tokType
	}
	return IDENTIFIER
}
