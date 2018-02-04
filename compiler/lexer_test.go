package compiler

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLexer(t *testing.T) {
	src, err := ioutil.ReadFile("../test.4g")
	if err != nil {
		panic(err)
	}
	scanner := &Scanner{src: []rune(string(src))}
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
		fmt.Printf("tok: '%v', lit: '%v', pos: %v\n", tok, lit, pos)
	}
}

