package lexer

import "github.com/ekediala/interpreter/token"

type Lexer struct {
	input        string
	position     int  // current read position in input. should point to the current character under evaluation.
	nextPosition int  // next position after position to be read and lexed
	ch           byte // current character under evaluation
}

// Reads the next character into l.ch and advances our cursor in the input
func (l *Lexer) ReadNextChar() {
	// advance cursors. why advance cursors with defer and not at the end of the function while using an else block to handle nextPosition being less than lenth of input? I just don't like else statements. I find that they make my code harder for me to follow. I prefer early returns.
	defer l.advancePositions()

	// prevent indexing out of array. If we are at the end of the input, set ch to zero [ASCII for "NUL"] so we can identify that as the end of lexing
	if l.nextPosition >= len(l.input) {
		l.ch = 0
		return
	}

	l.ch = l.input[l.nextPosition]
}

// Returns current token and advances the cursor
func (l *Lexer) ReadToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadNextChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Type = token.EQ
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadNextChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Type = token.NOT_EQ
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			// we have already advanced past the last character of the identifier, so return here to avoid calling ReadNextChar again
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}

	}

	l.ReadNextChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	currentPosition := l.position
	for isLetter(l.ch) {
		l.ReadNextChar()
	}

	return l.input[currentPosition:l.position]
}

func (l *Lexer) readNumber() string {
	currentPosition := l.position
	for isDigit(l.ch) {
		l.ReadNextChar()
	}

	return l.input[currentPosition:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}

func (l *Lexer) advancePositions() {
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadNextChar()
	}
}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    t,
		Literal: string(ch),
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func New(source string) *Lexer {
	l := Lexer{
		input: source,
	}

	// initialise position, nextPosition, and ch to their start points. will set ch to l.input[0] as zero value of nextPosition is 0, set position to 0[unnecessary as it is already 0] and then set nextPosition to 1
	l.ReadNextChar()

	// we can eliminate the call above and initialise the lexer with default values like below:
	// 	l := Lexer{
	// 	input:        source,
	// 	nextPosition: 1,
	// 	ch:           source[0],
	// }
	// will save us a function call, no checks for length of string, no setting position to 0 when it is already zero. I know which should be faster but I don't know which is better. I find hardcoding values weird.
	return &l
}
