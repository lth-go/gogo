package compiler

import (
	"fmt"
	"log"
	"strings"

	"github.com/lth-go/gogogogo/vm"
)

// 全局compiler列表
var stCompilerList []*Compiler

// 全局compiler
var stCurrentCompiler *Compiler

func getCurrentCompiler() *Compiler {
	return stCurrentCompiler
}

func setCurrentCompiler(c *Compiler) {
	stCurrentCompiler = c
}

// Compiler 编译器
type Compiler struct {
	// 词法解析器
	lexer *Lexer

	// 包名
	packageNameList []string // TODO: remove
	packageName     string
	// 源文件路径
	path string
	// 依赖的包
	importList []*ImportSpec

	// 函数列表
	funcList []*FunctionDefinition
	// 声明列表
	declarationList []*Declaration
	// 语句列表
	statementList []Statement

	// 当前块
	currentBlock *Block

	// 已加载compiler列表
	importedList []*Compiler

	// TODO 能否去掉
	// vm函数列表
	vmFunctionList []*vm.Function
}

func newCompiler() *Compiler {
	compilerBackup := getCurrentCompiler()
	c := &Compiler{
		importList:      []*ImportSpec{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		statementList:   []Statement{},
		importedList:    []*Compiler{},
		vmFunctionList:  []*vm.Function{},
	}
	setCurrentCompiler(c)
	// TODO 添加默认函数

	setCurrentCompiler(compilerBackup)
	return c
}

func (c *Compiler) getPackageName() string {
	return strings.Join(c.packageNameList, ".")
}

func (c *Compiler) addLexer(lexer *Lexer) {
	lexer.compiler = c
	c.lexer = lexer
}

func (c *Compiler) addLexerByPath(path string) {
	lexer := newLexerByFilePath(path)
	c.addLexer(lexer)

	c.path = path
}

//
// 函数定义
//
func (c *Compiler) functionDefine(typ *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) {
	// 定义重复
	if searchFunction(identifier) != nil || searchDeclaration(identifier, nil) != nil {
		compileError(typ.Position(), FUNCTION_MULTIPLE_DEFINE_ERR, identifier)
	}

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              identifier,
		packageNameList:   c.GetPackageNameList(),
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

//
// 编译
//
func (c *Compiler) compile(exeList *vm.ExecutableList, isRequired bool) *vm.Executable {
	compilerBackup := getCurrentCompiler()
	setCurrentCompiler(c)

	// 开始解析文件
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}

	for _, import_ := range c.importList {
		// 判断是否已经被解析过
		requireCompiler := searchCompiler(stCompilerList, import_.getPackageNameList())
		if requireCompiler != nil {
			c.importedList = append(c.importedList, requireCompiler)
			continue
		}

		requireCompiler = newCompiler()

		requireCompiler.packageNameList = import_.getPackageNameList()
		requireCompiler.packageName = import_.packageName

		c.importedList = append(c.importedList, requireCompiler)
		stCompilerList = append(stCompilerList, requireCompiler)

		// 获取要导入的全路径
		foundPath := import_.getFullPath()

		// 编译导入的包
		requireCompiler.addLexerByPath(foundPath)
		requireCompiler.compile(exeList, true)
	}

	// fix and generate
	c.fixTree()
	exe := c.generate()

	exe.Path = c.path
	exe.IsRequired = isRequired

	exeList.AddExe(exe)

	setCurrentCompiler(compilerBackup)

	return exe
}

//////////////////////////////
// 打印语法树
//////////////////////////////
func (c *Compiler) Show() {
	fmt.Println("==========")
	fmt.Println("stmt list start")
	for _, stmt := range c.statementList {
		stmt.show(0)
	}
	fmt.Println("\nstmt list end")
	fmt.Println("==========")
}

//////////////////////////////
// 修正树
//////////////////////////////
func (c *Compiler) fixTree() {
	// TODO remove
	// add function
	for _, fd := range c.funcList {
		c.addToVmFunctionList(fd)
	}

	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	// 修正函数
	for _, fd := range c.funcList {
		fd.fix()
	}

	// 修正全局声明
	for varCount, declaration := range c.declarationList {
		declaration.variableIndex = varCount
	}
}

// 添加VmFunction
func (c *Compiler) addToVmFunctionList(src *FunctionDefinition) int {

	srcPackageName := src.getPackageName()
	vmFuncName := src.getVmFuncName()

	for i, vmFunction := range c.vmFunctionList {
		if srcPackageName == vmFunction.PackageName && vmFuncName == vmFunction.Name {
			return i
		}
	}

	dest := &vm.Function{
		PackageName: srcPackageName,
		Name:        vmFuncName,
	}

	c.vmFunctionList = append(c.vmFunctionList, dest)

	return len(c.vmFunctionList) - 1
}

//////////////////////////////
// 生成字节码
//////////////////////////////
func (c *Compiler) generate() *vm.Executable {
	exe := vm.NewExecutable()
	exe.PackageName = c.getPackageName()

	exe.FunctionList = c.vmFunctionList

	// 添加全局变量声明
	c.addGlobalVariable(exe)
	// 添加函数信息
	c.addFunctions(exe)
	// 添加顶层代码
	c.addTopLevel(exe)

	return exe
}

// 添加全局变量
func (c *Compiler) addGlobalVariable(exe *vm.Executable) {
	for _, dl := range c.declarationList {

		newValue := vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))
		exe.GlobalVariableList = append(exe.GlobalVariableList, newValue)
	}
}

// 添加函数
func (c *Compiler) addFunctions(exe *vm.Executable) {

	inThisExes := make([]bool, len(c.vmFunctionList))

	for _, fd := range c.funcList {
		destIdx := c.getFunctionIndex(fd)
		inThisExes[destIdx] = true

		addFunction(exe, fd, c.vmFunctionList[destIdx], true)
	}

	for i, vmFunc := range c.vmFunctionList {
		if inThisExes[i] {
			continue
		}

		fd := searchFunction(vmFunc.Name)
		addFunction(exe, fd, vmFunc, false)
	}
}

func addFunction(exe *vm.Executable, src *FunctionDefinition, dest *vm.Function, inThisExe bool) {
	ob := newCodeBuf()

	dest.TypeSpecifier = copyTypeSpecifier(src.typeS())
	dest.ParameterList = copyParameterList(src.parameterList)

	if src.block != nil && inThisExe {
		generateStatementList(exe, src.block, src.block.statementList, ob)

		dest.IsImplemented = true
		dest.CodeList = ob.fixOpcodeBuf()
		dest.LineNumberList = ob.lineNumberList
		dest.LocalVariableList = copyLocalVariables(src)
	} else {
		dest.IsImplemented = false
		dest.LocalVariableList = nil
	}

	dest.IsMethod = false
}

func (c *Compiler) getFunctionIndex(src *FunctionDefinition) int {
	var funcName string

	srcPackageName := src.getPackageName()

	funcName = src.name

	for i, vmFunc := range c.vmFunctionList {
		if srcPackageName == vmFunc.PackageName && funcName == vmFunc.Name {
			return i
		}
	}

	panic("TODO")
}

// 添加字节码
func (c *Compiler) addTopLevel(exe *vm.Executable) {
	ob := newCodeBuf()
	generateStatementList(exe, nil, c.statementList, ob)

	exe.CodeList = ob.fixOpcodeBuf()
	exe.LineNumberList = ob.lineNumberList
}

func (c *Compiler) searchFunction(name string) *FunctionDefinition {

	// 当前compiler查找
	for _, pos := range c.funcList {
		if pos.name == name {
			return pos
		}
	}

	return nil
}

// ==============================
// 编译文件
// ==============================

func CompileFile(path string) *vm.ExecutableList {
	// 输出yacc错误信息
	yyErrorVerbose = true

	compiler := createCompilerByPath(path)

	exeList := compiler.Compile()

	return exeList
}

func createCompilerByPath(path string) *Compiler {
	compiler := newCompiler()
	compiler.addLexerByPath(path)

	return compiler
}

func (c *Compiler) Compile() *vm.ExecutableList {
	exeList := vm.NewExecutableList()
	exe := c.compile(exeList, false)
	exeList.TopLevel = exe

	return exeList
}

func (c *Compiler) GetPackageNameList() []string {
	return strings.Split(c.packageName, "/")
}
