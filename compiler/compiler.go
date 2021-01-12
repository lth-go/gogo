package compiler

import (
	"fmt"
	"log"
	"strings"

	"github.com/lth-go/gogo/vm"
)

var stCompilerList []*Compiler  // 全局compiler列表
var stCurrentCompiler *Compiler // 全局compiler

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

	// 已加载compiler列表
	importedList []*Compiler

	// 依赖的包
	importList []*ImportSpec
	// 函数列表
	funcList []*FunctionDefinition
	// 声明列表
	declarationList []*Declaration
	// 语句列表
	statementList []Statement

	ConstantList []interface{}

	// 当前块
	currentBlock *Block
}

func NewCompiler(path string) *Compiler {
	c := &Compiler{
		path:            path,
		importList:      []*ImportSpec{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		statementList:   []Statement{},
		importedList:    []*Compiler{},
		ConstantList:    []interface{}{},
	}

	c.SetLexer(path)

	return c
}

func (c *Compiler) getPackageName() string {
	return strings.Join(c.packageNameList, ".")
}

func (c *Compiler) SetLexer(path string) {
	lexer := newLexerByFilePath(path)
	lexer.compiler = c

	c.lexer = lexer
}

//
// 函数定义
//
func (c *Compiler) functionDefine(pos Position, receiver *Parameter, identifier string, typ *TypeSpecifier, block *Block) {
	var dummyBlock *Block

	// 定义重复
	if c.searchFunction(identifier) != nil || dummyBlock.searchDeclaration(identifier) != nil {
		compileError(pos, FUNCTION_MULTIPLE_DEFINE_ERR, identifier)
	}

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              identifier,
		packageNameList:   c.GetPackageNameList(),
		parameterList:     typ.funcType.Params,
		block:             block,
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
func (c *Compiler) compile(isRequired bool) []*vm.Executable {
	compilerBackup := getCurrentCompiler()
	setCurrentCompiler(c)

	// 开始解析文件
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}

	exeList := make([]*vm.Executable, 0)

	for _, import_ := range c.importList {
		// 判断是否已经被解析过
		importedC := searchCompiler(stCompilerList, import_.getPackageNameList())
		if importedC != nil {
			c.importedList = append(c.importedList, importedC)
			continue
		}

		importedC = NewCompiler(import_.getFullPath())
		importedC.packageNameList = import_.getPackageNameList()
		importedC.packageName = import_.packageName

		c.importedList = append(c.importedList, importedC)
		stCompilerList = append(stCompilerList, importedC)

		tmpExeList := importedC.compile(true)
		exeList = append(exeList, tmpExeList...)
	}

	// fix and generate
	c.fixTree()
	exe := c.Generate()

	exeList = append(exeList, exe)

	setCurrentCompiler(compilerBackup)

	return exeList
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
	// TODO: 添加原生函数
	c.AddNativeFunctions()

	// 修正函数
	for _, fd := range c.funcList {
		fd.fix()
	}

	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	// 修正全局声明
	for varCount, declaration := range c.declarationList {
		declaration.variableIndex = varCount
	}
}

func (c *Compiler) AddFuncList(fd *FunctionDefinition) int {
	packageName := fd.getPackageName()
	funcName := fd.getVmFuncName()

	for i, f := range c.funcList {
		if packageName == f.getPackageName() && funcName == f.getVmFuncName() {
			return i
		}
	}

	c.funcList = append(c.funcList, fd)

	return len(c.funcList) - 1
}

//
// 生成字节码
//
func (c *Compiler) Generate() *vm.Executable {
	exe := vm.NewExecutable()
	exe.PackageName = c.getPackageName()
	exe.FunctionList = c.GetVmFunctionList(exe)
	exe.VariableList.VariableList = c.GetVmVariableList() // 添加全局变量声明

	// 添加字节码
	opCodeBuf := newCodeBuf()
	generateStatementList(nil, c.statementList, opCodeBuf)

	exe.CodeList = opCodeBuf.fixOpcodeBuf()
	exe.LineNumberList = opCodeBuf.lineNumberList

	// TODO: remove
	exe.ConstantPool.SetPool(c.GetVmConstantList())

	return exe
}

func (c *Compiler) GetVmVariableList() []*vm.Variable {
	variableList := make([]*vm.Variable, 0)

	for _, dl := range c.declarationList {
		newValue := vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))
		variableList = append(variableList, newValue)
	}

	return variableList
}

func (c *Compiler) GetVmFunctionList(exe *vm.Executable) []*vm.Function {
	vmFuncList := make([]*vm.Function, 0)

	for _, fd := range c.funcList {
		vmFunc := c.GetVmFunction(exe, fd, fd.getPackageName() == c.packageName)
		vmFuncList = append(vmFuncList, vmFunc)
	}

	return vmFuncList
}

func (c *Compiler) GetVmFunction(exe *vm.Executable, src *FunctionDefinition, inThisExe bool) *vm.Function {
	ob := newCodeBuf()

	dest := &vm.Function{
		PackageName:   src.getPackageName(),
		Name:          src.name,
		TypeSpecifier: copyTypeSpecifier(src.typeS()),
		ParameterList: copyParameterList(src.parameterList),
		IsMethod:      false,
	}

	if src.block != nil && inThisExe {
		generateStatementList(src.block, src.block.statementList, ob)

		dest.IsImplemented = true
		dest.CodeList = ob.fixOpcodeBuf()
		dest.LineNumberList = ob.lineNumberList
		dest.LocalVariableList = copyLocalVariables(src)
	} else {
		dest.IsImplemented = false
		dest.LocalVariableList = nil
	}

	return dest
}

func (c *Compiler) getFunctionIndex(src *FunctionDefinition, exe *vm.Executable) int {
	var funcName string

	srcPackageName := src.getPackageName()
	funcName = src.name

	for i, vmFunc := range exe.FunctionList {
		if srcPackageName == vmFunc.PackageName && funcName == vmFunc.Name {
			return i
		}
	}

	panic("TODO")
}

func (c *Compiler) searchFunction(name string) *FunctionDefinition {
	for _, func_ := range c.funcList {
		if func_.name == name {
			return func_
		}
	}
	return nil
}

func (c *Compiler) searchPackage(name string) *Package {
	for _, importedC := range c.importedList {
		// TODO: 暂无处理重名
		lastName := importedC.packageNameList[len(importedC.packageNameList)-1]
		if name == lastName {
			return &Package{
				compiler: importedC,
				typ:      newTypeSpecifier(vm.BasicTypePackage),
			}
		}

	}
	return nil
}

//
// 编译文件
//

func CompileFile(path string) *vm.ExecutableList {
	// 输出yacc错误信息
	if true {
		yyErrorVerbose = true
	}

	compiler := createCompiler(path)

	return compiler.Compile()
}

func createCompiler(path string) *Compiler {
	compiler := NewCompiler(path)
	compiler.SetLexer(path)

	return compiler
}

func (c *Compiler) Compile() *vm.ExecutableList {
	exeList := vm.NewExecutableList()

	for _, exe := range c.compile(false) {
		exeList.AddExe(exe)
	}

	return exeList
}

func (c *Compiler) GetPackageNameList() []string {
	return strings.Split(c.packageName, "/")
}

func searchCompiler(list []*Compiler, packageName []string) *Compiler {
	for _, c := range list {
		if comparePackageName(c.packageNameList, packageName) {
			return c
		}
	}
	return nil
}

func comparePackageName(packageNameList1, packageNameList2 []string) bool {
	if packageNameList1 == nil {
		if packageNameList2 == nil {
			return true
		}
		return false
	}

	length1 := len(packageNameList1)
	length2 := len(packageNameList2)

	if length1 != length2 {
		return false
	}

	for i := 0; i < length1; i++ {
		if packageNameList1[i] != packageNameList2[i] {
			return false
		}
	}

	return true
}

func (c *Compiler) AddNativeFunctions() {
	typ := createFuncTypeSpecifier([]*Parameter{{typeSpecifier: newTypeSpecifier(vm.BasicTypeString)}}, nil)

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              "print",
		packageNameList:   c.GetPackageNameList(),
		parameterList:     []*Parameter{{typeSpecifier: newTypeSpecifier(vm.BasicTypeString), name: "str"}},
		block:             nil,
		localVariableList: nil,
	}

	c.funcList = append(c.funcList, fd)
}

func (c *Compiler) AddConstantList(value interface{}) int {
	c.ConstantList = append(c.ConstantList, value)
	return len(c.ConstantList) - 1
}

func (c *Compiler) GetVmConstantList() []vm.Constant {
	var constantValue vm.Constant

	constantList := make([]vm.Constant, 0)

	for _, valueIFS := range c.ConstantList {
		switch value := valueIFS.(type) {
		case int:
			constantValue = vm.NewConstantInt(value)
		case float64:
			constantValue = vm.NewConstantDouble(value)
		case string:
			constantValue = vm.NewConstantString(value)
		default:
			panic("TODO")
		}

		constantList = append(constantList, constantValue)
	}

	return constantList
}
