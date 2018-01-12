package parser

import (
	"fmt"
	"testing"
)

var src = `number a = 1 + 1; b = a;
if (a > b) {
	c = a;
}
`

func TestLexer(t *testing.T) {


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

func TestParse(t *testing.T) {
	stmts, err := ParseSrc(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(stmts))
}
