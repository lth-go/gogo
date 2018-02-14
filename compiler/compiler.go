package compiler

import (
	//"os"
	"../vm"
	"fmt"
)

// Compiler 编译器
type Compiler struct {
	// 词法解析器
	lexer *Lexer

	// 语句列表
	statementList []Statement

	// 函数列表
	funcList []*FunctionDefinition

	// 声明列表
	declarationList []*Declaration

	// 当前行数
	currentLineNumber int

	// 当前块
	currentBlock *Block
}

func (c *Compiler) functionDefine(typeSpecifier *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) {
	// 定义重复
	if SearchFunction(identifier) != nil || searchDeclaration(identifier, nil) != nil {
		panic("TODO")
	}

	fd := &FunctionDefinition{
		typeSpecifier:    typeSpecifier,
		name:             identifier,
		parameterList:     parameterList,
		block:         block,
		index:         len(c.funcList),
		localVariableList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

func newCompiler() *Compiler {
	c := &Compiler{
		statementList:   []Statement{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
	}
	setCurrentCompiler(c)
	return c
}

func (c *Compiler) Show() {
	fmt.Println("==========")
	fmt.Println("stmt list start\n")
	for _, stmt := range c.statementList {
		stmt.show(0)
	}
	fmt.Println("\nstmt list end")
	fmt.Println("==========\n")
}

// 修正树
func (c *Compiler) fixTree() {
	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	// 修正函数
	for _, fd := range c.funcList {
		if fd.block == nil {
			continue
		}
		// 添加形参声明
		fd.addParameterAsDeclaration()

		// 修正表达式列表
		fixStatementList(fd.block, fd.block.statementList, fd)

		// 修正返回值
		fd.addReturnFunction()
	}

	// 修正全局声明
	for varCount, decl := range c.declarationList {
		decl.variableIndex = varCount
	}

}

func (c *Compiler) Generate(exe *vm.Executable){

	// 添加全局变量声明
	addGlobalVariable(c, exe)
	// 添加函数信息
	addFunctions(c, exe)
	// 添加顶层代码
	addTopLevel(c, exe)
}

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

// ==============================
// utils
// ==============================

func searchDeclaration(name string, currentBlock *Block) *Declaration {

	for b := currentBlock; b != nil; b = b.outerBlock {
		for _, d := range b.declarationList {
			if d.name == name {
				return d
			}
		}
	}

	compiler := getCurrentCompiler()

	for _, d := range compiler.declarationList {
		if d.name == name {
			return d
		}
	}

	return nil
}

func SearchFunction(name string) *FunctionDefinition {
	compiler := getCurrentCompiler()

	for _, pos := range compiler.funcList {
		if pos.name == name {
			return pos
		}
	}
	return nil
}

func compileError(pos Position, errorNumber int, a ...interface{}) {
	fmt.Println("编译错误")
	fmt.Printf("Line: %d:%d\n", pos.Line, pos.Column)
	fmt.Printf(errMessageMap[errorNumber], a...)
	fmt.Println("\n")
	panic("打印栈，看看哪里出错了")
	//os.Exit(1)
}

var stCurrentCompiler *Compiler

func getCurrentCompiler() *Compiler {
	return stCurrentCompiler
}

func setCurrentCompiler(c *Compiler) {
	stCurrentCompiler = c
}

// ==============================
// parse 词法解析，语法解析
// ==============================

// ParseSrc 解析源码
func ParseSrc(src string) (*Compiler, error) {
	// 输出yacc错误信息
	yyErrorVerbose = true
	scanner := &Scanner{src: []rune(src)}
	return parse(scanner)
}

// parse provides way to parse the code using Scanner.
func parse(s *Scanner) (*Compiler, error) {
	compiler := newCompiler()

	lexer := &Lexer{s: s, compiler: compiler}

	compiler.lexer = lexer

	if yyParse(lexer) != 0 {
		return nil, lexer.e
	}

	// 修正树
	lexer.compiler.fixTree()

	return lexer.compiler, lexer.e
}
