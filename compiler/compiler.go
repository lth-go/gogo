package compiler

import (
	"os"
	"../vm"
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
	if searchFunction(identifier) || searchDeclaration(identifier, nil) {
		compileError(nil, 0, "")
	}

	fd := &FunctionDefinition{
		typeSpecifier:    typeSpecifier,
		name:             identifier,
		fd.parameter:     parameterList,
		fd.block:         block,
		fd.index:         len(c.funcList),
		fd.localVariable: nil,
	}

	if block != nil {
		block.parent = FunctionBlockInfo{function: fd}
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

// 修正树
func (c *Compiler) fixTree() {
	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

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

	for varCount, decl := range c.declarationList {
		decl.variableIndex = varCount
	}

}

func (c *Compiler) generate() *vm.Executable {
	exe := newExecutable()

	addGlobalVariable(c, exe)
	addFunctions(c, exe)
	addTopLevel(c, exe)

	return exe
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

func searchFunction(name string) *FunctionDefinition {
	compiler := getCurrentCompiler()

	for _, pos := range compiler.funcList {
		if pos.name == name {
			return pos
		}
	}
	return nil
}

func compileError(pos Position, compilerError int, message string) {
	os.Exit(1)
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

	l := &Lexer{s: s, compiler: compiler}

	compiler.lexer = l

	if yyParse(l) != 0 {
		return nil, l.e
	}

	// 修正树
	l.compiler.fixTree()

	return l.compiler, l.e
}
