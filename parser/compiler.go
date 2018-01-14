package parser

import (
	"os"
)

// Compiler 编译器
type Compiler struct {
	// 词法解析器
	lexer *Lexer

	// 语句列表
	stmts []Statement

	// 函数列表
	funcList []*FunctionDefinition

	// 声明列表
	declarationList []*DeclarationStatement

	// 当前行数
	currentLineNumber int

	// 当前块
	currentBlock *Block
}

func (c *Compiler) functionDefine(typeSpecifier *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) {

}

func newCompiler() *Compiler {
	c := &Compiler{
		stmts:           []Statement{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*DeclarationStatement{},
	}
	setCurrentCompiler(c)
	return c
}

func compileError(lineNumber int, compilerError int, message string) {
	os.Exit(1)
}

var stCurrentCompiler *Compiler

func getCurrentCompiler() *Compiler {
	return stCurrentCompiler
}

func setCurrentCompiler(c *Compiler) {
	stCurrentCompiler = c
}
