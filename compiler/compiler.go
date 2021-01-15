package compiler

import (
	"fmt"
	"log"
	"strings"

	"github.com/lth-go/gogo/vm"
)

var (
	stCompilerList    []*Compiler // 全局compiler列表
	stCurrentCompiler *Compiler   // 全局compiler
)

func getCurrentCompiler() *Compiler {
	return stCurrentCompiler
}

func setCurrentCompiler(c *Compiler) {
	stCurrentCompiler = c
}

// Compiler 编译器
type Compiler struct {
	lexer           *Lexer                // 词法解析器
	path            string                // 源文件路径
	packageName     string                // 包名
	importedList    []*Compiler           // 已加载compiler列表
	importList      []*Import         // 依赖的包
	funcList        []*FunctionDefinition // 函数列表
	declarationList []*Declaration        // 声明列表
	statementList   []Statement           // 语句列表
	ConstantList    []interface{}         // constant
	currentBlock    *Block                // 当前块
}

func NewCompiler(path string) *Compiler {
	c := &Compiler{
		path:            path,
		importList:      []*Import{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		statementList:   []Statement{},
		importedList:    []*Compiler{},
		ConstantList:    []interface{}{},
	}

	c.lexer = NewLexer(path, c)
	return c
}

//
// 函数定义
//
func createFunctionDefine(pos Position, receiver *Parameter, identifier string, typ *Type, block *Block) {
	c := getCurrentCompiler()

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              identifier,
		packageName:       c.packageName,
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
		importedCompiler := searchCompiler(stCompilerList, import_.packageName)
		if importedCompiler != nil {
			c.importedList = append(c.importedList, importedCompiler)
			continue
		}

		// new compiler
		importedCompiler = NewCompiler(import_.GetPath())
		importedCompiler.packageName = import_.packageName

		// add global
		c.importedList = append(c.importedList, importedCompiler)
		stCompilerList = append(stCompilerList, importedCompiler)

		// parse
		tmpExeList := importedCompiler.compile(true)
		exeList = append(exeList, tmpExeList...)
	}

	// fix and generate
	c.FixTree()
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

//
// 修正树
//
func (c *Compiler) FixTree() {
	// TODO: check func list, if is redifined

	// 原先函数在func fix之前添加,类型c头文件
	// 表达式fix中会添加其他包函数到本包

	// 添加原生函数
	c.AddNativeFunctions()

	// 修正函数
	for _, fd := range c.funcList {
		fd.fix()
	}

	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)
}

func (c *Compiler) AddFuncList(fd *FunctionDefinition) int {
	packageName := fd.GetPackageName()
	name := fd.GetName()

	for i, f := range c.funcList {
		if packageName == f.GetPackageName() && name == f.GetName() {
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
	exe.PackageName = c.packageName
	exe.FunctionList = c.GetVmFunctionList(exe)
	exe.VariableList.VariableList = c.GetVmVariableList() // 添加全局变量声明

	// 添加字节码
	opCodeBuf := newCodeBuf()
	generateStatementList(nil, c.statementList, opCodeBuf)

	exe.CodeList = opCodeBuf.fixOpcodeBuf()
	exe.LineNumberList = opCodeBuf.lineNumberList

	// TODO: remove
	exe.ConstantPool.SetPool(c.ConstantList)

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
		vmFunc := c.GetVmFunction(exe, fd, fd.GetPackageName() == c.packageName)
		vmFuncList = append(vmFuncList, vmFunc)
	}

	return vmFuncList
}

func (c *Compiler) GetVmFunction(exe *vm.Executable, src *FunctionDefinition, inThisExe bool) *vm.Function {
	ob := newCodeBuf()

	dest := &vm.Function{
		PackageName:   src.GetPackageName(),
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

	srcPackageName := src.GetPackageName()
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
		if name == importedC.packageName {
			return &Package{
				compiler: importedC,
				typ:      NewType(vm.BasicTypePackage),
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

	compiler := NewCompiler(path)

	return compiler.Compile()
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

func searchCompiler(list []*Compiler, packageName string) *Compiler {
	for _, c := range list {
		if c.packageName == packageName {
			return c
		}
	}
	return nil
}

func (c *Compiler) AddNativeFunctions() {
	paramsType := []*Parameter{{typeSpecifier: NewType(vm.BasicTypeString), name: "str"}}
	typ := CreateFuncType(paramsType, nil)

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              "print",
		packageName:       "_sys",
		parameterList:     paramsType,
		block:             nil,
		localVariableList: nil,
	}

	c.funcList = append(c.funcList, fd)
}

func (c *Compiler) AddConstantList(value interface{}) int {
	for i, v := range c.ConstantList {
		if value == v {
			return i
		}
	}

	c.ConstantList = append(c.ConstantList, value)
	return len(c.ConstantList) - 1
}

// 添加声明
func (c *Compiler) AddDeclarationList(decl *Declaration) int {
	c.declarationList = append(c.declarationList, decl)
	decl.variableIndex = len(c.declarationList) - 1

	return decl.variableIndex
}
