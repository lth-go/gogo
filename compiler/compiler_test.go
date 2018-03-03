package compiler

import (
	"testing"
)

var testFile = "../test/shape.4g"

//func TestLexer(t *testing.T) {
//    l := newLexerByFilePath(testFile)
//    l.show()
//}

func TestParse(t *testing.T) {
	yyErrorVerbose = true

	compiler := createCompilerByPath(testFile)
	//exeList := compiler.Compile()
	compiler.Compile()

	for _, c := range stCompilerList {
		println("=======")
		c.Show()
	}
}
