package compiler

import (
	"fmt"
	"io/ioutil"
	"testing"
	"../vm"
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

func TestParse(t *testing.T) {
	src, err := ioutil.ReadFile("../test.4g")
	if err != nil {
		panic(err)
	}
	compiler, err := ParseSrc(string(src))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, stmt := range compiler.statementList {
		stmt.show(0)
	}
}

func TestGenerate(t *testing.T) {
	src, err := ioutil.ReadFile("../test.4g")
	if err != nil {
		panic(err)
	}
	compiler, err := ParseSrc(string(src))
	if err != nil {
		fmt.Println(err)
		return
	}
	exe := vm.NewExecutable()
	compiler.Generate(exe)

	for i:=0; i<len(exe.CodeList); {
		code := exe.CodeList[i]
		info :=vm.OpcodeInfo[int(code)]
		paramList := []byte(info.Parameter)

		fmt.Println(info.Mnemonic)
		for _, param := range paramList {
			switch param {
			case 'b':
				i += 1
			case 's', 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
		i += 1
	}
}
