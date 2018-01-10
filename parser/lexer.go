// Package parser implements parser for anko.
package parser

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
	s     *Scanner
	lit   string
	pos   Position
	e     error

	stmts []Statement

	funcList []*FunctionDefinition

	currentBlock *Bolck
}

func (l *Lexer) functionDefine(typeSpecifier BasicType, identifier string, parameterList []*Parameter, block *Block) {

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

// ==============================
// parse
// ==============================

// ParseSrc provides way to parse the code from source.
func ParseSrc(src string) ([]Statement, error) {
	scanner := &Scanner{
		src: []rune(src),
	}
	return Parse(scanner)
}

// Parse provides way to parse the code using Scanner.
func Parse(s *Scanner) ([]Statement, error) {
	l := Lexer{s: s}
	if yyParse(&l) != 0 {
		return nil, l.e
	}
	return l.stmts, l.e
}
