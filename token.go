package main

import (
	"fmt"
)

type TokenType int

func (t TokenType) String() string {
	switch t {
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Semicolon:
		return "Semicolon"
	case Slash:
		return "Slash"
	case Star:
		return "Star"

		// One or two character tokens.
	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case Equal:
		return "Equal"
	case EqualEqual:
		return "EqualEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Less:
		return "Less"
	case LessEqual:
		return "LessEqual"

		// Literals.
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"

		// Keywords.
	case And:
		return "And"
	case Class:
		return "Class"
	case Else:
		return "Else"
	case False:
		return "False"
	case Fun:
		return "Fun"
	case For:
		return "For"
	case If:
		return "If"
	case Nil:
		return "Nil"
	case Or:
		return "Or"
	case Print:
		return "Print"
	case Return:
		return "Return"
	case Super:
		return "Super"
	case This:
		return "This"
	case True:
		return "True"
	case Var:
		return "Var"
	case While:
		return "While"
	case EOF:
		return "EOF"
	}
	return "Unknown Token"
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
