package lexer

import (
	"github.com/Jamess-Lucass/interpreter-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	character    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readCharacter()

	return l
}

func (l *Lexer) readCharacter() {
	if l.readPosition >= len(l.input) {
		l.character = 0 // ASCII code for "NUL"
	} else {
		l.character = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.character {
	case '=':
		if l.peekCharacter() == '=' {
			character := l.character
			l.readCharacter()
			tok = token.Token{Type: token.EQ, Literal: string(character) + string(l.character)}
		} else {
			tok = token.NewToken(token.ASSIGN, l.character)
		}
	case '!':
		if l.peekCharacter() == '=' {
			character := l.character
			l.readCharacter()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(character) + string(l.character)}
		} else {
			tok = token.NewToken(token.BANG, l.character)
		}
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.character)
	case '(':
		tok = token.NewToken(token.LPAREN, l.character)
	case ')':
		tok = token.NewToken(token.RPAREN, l.character)
	case ',':
		tok = token.NewToken(token.COMMA, l.character)
	case '+':
		tok = token.NewToken(token.PLUS, l.character)
	case '{':
		tok = token.NewToken(token.LBRACE, l.character)
	case '}':
		tok = token.NewToken(token.RBRACE, l.character)
	case '-':
		tok = token.NewToken(token.MINUS, l.character)
	case '/':
		tok = token.NewToken(token.SLASH, l.character)
	case '*':
		tok = token.NewToken(token.ASTERISK, l.character)
	case '<':
		tok = token.NewToken(token.LT, l.character)
	case '>':
		tok = token.NewToken(token.GT, l.character)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.character) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.character) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, l.character)
		}
	}

	l.readCharacter()

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		l.readCharacter()
	}
}

func (l *Lexer) readIdentifier() string {
	currentPosition := l.position
	for isLetter(l.character) {
		l.readCharacter()
	}

	return l.input[currentPosition:l.position]
}

func (l *Lexer) readNumber() string {
	currentPosition := l.position
	for isDigit(l.character) {
		l.readCharacter()
	}

	return l.input[currentPosition:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekCharacter() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}
