package main

import (
	"strconv"
)

type Scanner struct {
	source  []rune // rune slice! support utf8 strings!! no utf8 identifiers tho :sob:
	Tokens  []*Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func newScanner(source string) *Scanner {
	return &Scanner{[]rune(source), make([]*Token, 0), 0, 0, 1}
}

func (s *Scanner) ScanTokens() []*Token {
	defer func() {
		recover()
	}()
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, newToken(EOF, "", nil, s.line))
	return s.Tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen, nil)
	case ')':
		s.addToken(RightParen, nil)
	case '{':
		s.addToken(LeftBrace, nil)
	case '}':
		s.addToken(RightBrace, nil)
	case ',':
		s.addToken(Comma, nil)
	case '.':
		s.addToken(Dot, nil)
	case '-':
		s.addToken(Minus, nil)
	case '+':
		s.addToken(Plus, nil)
	case ';':
		s.addToken(Semicolon, nil)
	case '*':
		s.addToken(Star, nil)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual, nil)
		} else {
			s.addToken(Bang, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual, nil)
		} else {
			s.addToken(Equal, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual, nil)
		} else {
			s.addToken(Less, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual, nil)
		} else {
			s.addToken(Greater, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash, nil)
		}
	case ' ', '\r', '\t':
		//ignore
	case '\n':
		s.line++
	case '"':
		s.getString()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			emitError(s.line, "Unexpected character.")
		}
	}
}
func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	t, ok := keywords[text]
	if ok == false {
		t = Identifier
	}
	s.addToken(t, nil)
}
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	s.addToken(Number, value)
}
func (s *Scanner) getString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		emitError(s.line, "Unterminated string.")
		return
	}
	s.advance()
	// Trim the quotes
	value := string(s.source[s.start+1 : s.current-1])
	s.addToken(String, value)
}
func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}
func (s Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}
func (s Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}
func (s Scanner) isDigit(c rune) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}
func (s Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) advance() rune {
	curr := rune(s.source[s.current])
	s.current++
	return curr
}
func (s *Scanner) addToken(t TokenType, literal any) {
	text := s.source[s.start:s.current]
	//fmt.Println(t, "(", literal, ")")
	s.Tokens = append(s.Tokens, newToken(t, string(text), literal, s.line))
}
