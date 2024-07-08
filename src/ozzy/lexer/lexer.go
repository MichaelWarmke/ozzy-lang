package lexer

import (
	"fmt"
	"ozzy/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.createTwoCharToken(token.EQUALS)
		} else {
			tok = token.Token{token.ASSIGN, string(l.ch)}
		}
	case ';':
		tok = token.Token{token.SEMICOLON, string(l.ch)}
	case '(':
		tok = token.Token{token.LPAREN, string(l.ch)}
	case ')':
		tok = token.Token{token.RPAREN, string(l.ch)}
	case ',':
		tok = token.Token{token.COMMA, string(l.ch)}
	case '+':
		tok = token.Token{token.PLUS, string(l.ch)}
	case '-':
		tok = token.Token{token.MINUS, string(l.ch)}
	case '!':
		if l.peekChar() == '=' {
			tok = l.createTwoCharToken(token.NOT_EQUALS)
		} else {
			tok = token.Token{token.BANG, string(l.ch)}
		}
	case '*':
		tok = token.Token{token.ASTERISK, string(l.ch)}
	case '/':
		tok = token.Token{token.SLASH, string(l.ch)}
	case '<':
		tok = token.Token{token.LT, string(l.ch)}
	case '>':
		tok = token.Token{token.GT, string(l.ch)}
	case '{':
		tok = token.Token{token.LBRACE, string(l.ch)}
	case '}':
		tok = token.Token{token.RBRACE, string(l.ch)}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readIdentifier(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = token.Token{token.ILLEGAL, string(l.ch)}
		}
	}
	l.readChar()

	return tok
}

func (l *Lexer) createTwoCharToken(t token.TokenType) token.Token {
	ch := l.ch
	l.readChar()
	return token.Token{t, string(ch) + string(l.ch)}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

type predicate func(ch byte) bool

func (l *Lexer) readIdentifier(fn predicate) string {
	position := l.position
	for fn(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) printTokenString() {
	for range l.input {
		tok := l.NextToken()
		fmt.Println("{token." + string(tok.Type) + ", \"" + tok.Literal + "\"")
	}
}
