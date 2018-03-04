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

	//for _, todo := range compiler.funcList {
	//    if todo.block != nil {
	//        println("HHHHHHHHH")
	//        println(todo.name)
	//        todo.block.show(0)
	//    }
	//}
	//for _, c := range stCompilerList {
	//    println("=======")
	//    c.Show()
	//}
}
