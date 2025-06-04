package main

import (
	"fmt"
)

type TokenType int

func (t TokenType) String() string {
	switch t {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case Comma:
		return ","
	case Dot:
		return "."
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Semicolon:
		return ";"
	case Slash:
		return "/"
	case Star:
		return "*"

		// One or two character tokens.
	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case EqualEqual:
		return "=="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Less:
		return "<"
	case LessEqual:
		return "<="

		// Literals.
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"

		// Keywords.
	case And:
		return "and"
	case Class:
		return "class"
	case Else:
		return "else"
	case False:
		return "false"
	case Fun:
		return "fun"
	case For:
		return "for"
	case If:
		return "if"
	case Nil:
		return "nil"
	case Or:
		return "or"
	case Print:
		return "print"
	case Return:
		return "return"
	case Super:
		return "super"
	case This:
		return "this"
	case True:
		return "true"
	case Var:
		return "var"
	case While:
		return "while"
	case EOF:
		return "EOF"
	}
	panic("Unknown Token")
}

const (
	// Single-character tokens.
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func newToken(t TokenType, lexeme string, literal any, line int) *Token {
	return &Token{t, lexeme, literal, line}
}
func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}
