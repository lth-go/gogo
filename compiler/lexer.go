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
	msg := fmt.Sprintf("\nLine: %d\nMessage:%s\n", e.Pos.Line, e.Message)
	return msg
}

// ==============================
// Lexer
// ==============================

// Lexer provides inteface to parse codes.
type Lexer struct {
	s        *Scanner
	lit      string
	pos      Position
	e        error
	compiler *Compiler
}

func newLexerByFilePath(path string) *Lexer {
	return &Lexer{
		s: newScannerByFilePath(path),
	}
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

func (l *Lexer) show() {
	for {
		tok, lit, pos, err := l.s.Scan()
		if err != nil {
			panic(err)
		}
		if tok == EOF {
			fmt.Println("end")
			break
		}
		fmt.Printf("tok: '%v', lit: '%v', pos: %v\n", tok, lit, pos)
	}
}

// Error sets parse error.
// parse的错误
func (l *Lexer) Error(msg string) {
	l.e = &Error{Message: msg, Pos: l.pos, Fatal: false}
}
