package compiler

import (
	"fmt"
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

	// 当前块
	currentBlock *Block
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

func (c *Compiler) Generate(exe *vm.Executable) {

	// 添加全局变量声明
	c.addGlobalVariable(exe)
	// 添加函数信息
	c.addFunctions(exe)
	// 添加顶层代码
	c.addTopLevel(exe)
}

func (compiler *Compiler) addGlobalVariable(exe *vm.Executable) {
	for _, dl := range compiler.declarationList {
		v := vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))

		exe.GlobalVariableList = append(exe.GlobalVariableList, v)
	}
}

// 为每个函数生成所需的信息
func (compiler *Compiler) addFunctions(exe *vm.Executable) {
	for _, srcFd := range compiler.funcList {
		destFd := &vm.VmFunction{}
		copyFunction(srcFd, destFd)

		exe.FunctionList = append(exe.FunctionList, destFd)

		if srcFd.block == nil {
			// 原生函数
			destFd.IsImplemented = false
			continue
		}

		ob := newCodeBuf()
		generateStatementList(exe, srcFd.block, srcFd.block.statementList, ob)

		destFd.IsImplemented = true
		destFd.CodeList = ob.fixOpcodeBuf()
		destFd.LineNumberList = ob.lineNumberList
	}
}

// 生成解释器所需的信息
func (compiler *Compiler) addTopLevel(exe *vm.Executable) {
	ob := newCodeBuf()
	generateStatementList(exe, nil, compiler.statementList, ob)

	exe.CodeList = ob.fixOpcodeBuf()
	exe.LineNumberList = ob.lineNumberList
}

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

func (c *Compiler) functionDefine(typ *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) {
	// 定义重复
	if SearchFunction(identifier) != nil || searchDeclaration(identifier, nil) != nil {
		compileError(typ.Position(), FUNCTION_MULTIPLE_DEFINE_ERR, identifier)
	}

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              identifier,
		parameterList:     parameterList,
		block:             block,
		index:             len(c.funcList),
		localVariableList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

// ==============================
// utils
// ==============================

func searchDeclaration(name string, currentBlock *Block) *Declaration {

	// 从局部作用域查找
	for b := currentBlock; b != nil; b = b.outerBlock {
		for _, decl := range b.declarationList {
			if decl.name == name {
				return decl
			}
		}
	}

	// 从全局作用域查找
	compiler := getCurrentCompiler()
	for _, decl := range compiler.declarationList {
		if decl.name == name {
			return decl
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
	os.Exit(1)
}

// 设置全局compiler
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
