package parser

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {

	src := "a = 1 + 1\n"

	scanner := &Scanner{src: []rune(src)}
	l := Lexer{s: scanner}
	for {
		tok, lit, pos, err := l.s.Scan()
		if err != nil {
			panic(err)
		}
		if tok == EOF {
			fmt.Println("end")
			break
		}
		fmt.Printf("token: %v, lit: %v, pos: %v\n", tok, lit, pos)
	}
}
