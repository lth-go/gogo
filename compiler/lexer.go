package compiler

import (
	"fmt"
)

// ==============================
// Token
// ==============================

// Token token
type Token struct {
	PosImpl // StmtImpl provide Pos() function.
	Tok     int
	Lit     string
}

// ==============================
// Error
// ==============================

// Error provides a convenient interface for handling runtime error.
// It can be Error inteface with type cast which can call Pos().
type Error struct {
	Message  string
	Pos      Position
	Filename string
	Fatal    bool
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// ==============================
// Lexer
// ==============================

// Lexer provides inteface to parse codes.
type Lexer struct {
	s   *Scanner
	lit string
	pos Position
	e   error

	compiler *Compiler
}


// Lex scans the token and literals.
func (l *Lexer) Lex(lval *yySymType) int {
	tok, lit, pos, err := l.s.Scan()
	if err != nil {
		l.e = &Error{Message: fmt.Sprintf("%s", err.Error()), Pos: pos, Fatal: true}
	}
	lval.tok = Token{Tok: tok, Lit: lit}
	lval.tok.SetPosition(pos)
	l.lit = lit
	l.pos = pos
	return tok
}

// Error sets parse error.
func (l *Lexer) Error(msg string) {
	l.e = &Error{Message: msg, Pos: l.pos, Fatal: false}
}

