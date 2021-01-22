package compiler

import (
	"log"
	"strings"

	"github.com/lth-go/gogo/vm"
)

var (
	CompilerList    []*Compiler // 全局compiler列表
	CurrentCompiler *Compiler   // 全局compiler
)

func GetCurrentCompiler() *Compiler {
	return CurrentCompiler
}

func SetCurrentCompiler(c *Compiler) {
	CurrentCompiler = c
}

// Compiler 编译器
type Compiler struct {
	lexer           *Lexer                // 词法解析器
	path            string                // 源文件路径
	packageName     string                // 包名
	importList      []*Import             // 依赖的包
	funcList        []*FunctionDefinition // 函数列表
	declarationList []*Declaration        // 声明列表
	ConstantList    []interface{}         // 常量定义
	currentBlock    *Block                // 当前块
}

func NewCompiler(path string) *Compiler {
	c := &Compiler{
		path:            path,
		importList:      []*Import{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		ConstantList:    []interface{}{},
	}

	c.lexer = NewLexer(path, c)
	return c
}

//
// 函数定义
//
func createFunctionDefine(pos Position, receiver *Parameter, identifier string, typ *Type, block *Block) {
	c := GetCurrentCompiler()

	fd := &FunctionDefinition{
		Type:            typ,
		Name:            identifier,
		PackageName:     c.packageName,
		ParameterList:   typ.funcType.Params,
		Block:           block,
		DeclarationList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

//
// 编译
//
func (c *Compiler) compile() []*vm.Executable {
	compilerBackup := GetCurrentCompiler()
	SetCurrentCompiler(c)

	// 开始解析文件
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}

	exeList := make([]*vm.Executable, 0)

	for _, import_ := range c.importList {
		// 判断是否已经被解析过
		importedCompiler := searchCompiler(CompilerList, import_.packageName)
		if importedCompiler != nil {
			continue
		}

		// new compiler
		importedCompiler = NewCompiler(import_.GetPath())

		// add global
		CompilerList = append(CompilerList, importedCompiler)

		// parse
		tmpExeList := importedCompiler.compile()
		exeList = append(exeList, tmpExeList...)
	}

	// fix and generate
	c.FixTree()
	exe := c.Generate()

	exeList = append(exeList, exe)

	SetCurrentCompiler(compilerBackup)

	return exeList
}

//
// 修正树
//
func (c *Compiler) FixTree() {
	// 修正全局变量
	c.FixDeclarationList()
	// TODO: check func list, if is redifined

	// 原先函数在func fix之前添加,类型c头文件
	// 表达式fix中会添加其他包函数到本包

	// 添加原生函数
	c.AddNativeFunctions()

	// 修正函数
	for _, fd := range c.funcList {
		fd.fix()
	}
}

func (c *Compiler) FixDeclarationList() {
	for _, decl := range c.declarationList {
		if decl.InitValue != nil {
			decl.InitValue = decl.InitValue.fix(nil)
			decl.InitValue = CreateAssignCast(decl.InitValue, decl.Type)
		}
	}
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
	opCodeBuf := NewOpCodeBuf()

	exe.CodeList = opCodeBuf.fixOpcodeBuf()
	exe.LineNumberList = opCodeBuf.lineNumberList

	// TODO: remove
	exe.ConstantPool.SetPool(c.ConstantList)

	return exe
}

func (c *Compiler) GetVmVariableList() []*vm.Variable {
	variableList := make([]*vm.Variable, 0)

	for _, decl := range c.declarationList {
		newValue := vm.NewVmVariable(decl.PackageName, decl.Name, CopyToVmType(decl.Type))
		newValue.Value = GetVmVariable(decl.InitValue)
		variableList = append(variableList, newValue)
	}

	return variableList
}

func GetVmVariable(valueIFS Expression) interface{} {
	if valueIFS == nil {
		return nil
	}

	switch value := valueIFS.(type) {
	case *BoolExpression:
		return value.Value
	case *IntExpression:
		return value.Value
	case *FloatExpression:
		return value.Value
	case *StringExpression:
		return value.Value
	case *ArrayExpression:
		return TODOGetVmVariable(value)
	}

	return nil
}

func TODOGetVmVariable(valueIFS Expression) vm.Object {
	switch value := valueIFS.(type) {
	// case *BoolExpression:
	//     return vm.NewObjectInt()
	case *IntExpression:
		return vm.NewObjectInt(value.Value)
	case *FloatExpression:
		return vm.NewObjectFloat(value.Value)
	case *StringExpression:
		return vm.NewObjectString(value.Value)
	case *ArrayExpression:
		arrayValue := vm.NewObjectArray(len(value.List))
		for i, subValue := range value.List {
			arrayValue.List[i] = TODOGetVmVariable(subValue)
		}
		return arrayValue
	}

	return nil
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
	ob := NewOpCodeBuf()

	dest := &vm.Function{
		PackageName: src.GetPackageName(),
		Name:        src.Name,
		Type:        CopyToVmType(src.GetType()),
		IsMethod:    false,
	}

	if src.Block != nil && inThisExe {
		generateStatementList(src.Block, src.Block.statementList, ob)

		dest.IsImplemented = true
		dest.CodeList = ob.fixOpcodeBuf()
		dest.LineNumberList = ob.lineNumberList
		dest.LocalVariableList = copyVmVariableList(src)
	} else {
		dest.IsImplemented = false
		dest.LocalVariableList = nil
	}

	return dest
}

func (c *Compiler) getFunctionIndex(src *FunctionDefinition, exe *vm.Executable) int {
	var funcName string

	srcPackageName := src.GetPackageName()
	funcName = src.Name

	for i, vmFunc := range exe.FunctionList {
		if srcPackageName == vmFunc.PackageName && funcName == vmFunc.Name {
			return i
		}
	}

	panic("TODO")
}

func (c *Compiler) searchFunction(name string) *FunctionDefinition {
	for _, func_ := range c.funcList {
		if func_.Name == name {
			return func_
		}
	}
	return nil
}

func (c *Compiler) searchPackage(name string) *Package {
	for _, imp := range c.importList {
		if name == imp.packageName {
			for _, importCompiler := range CompilerList {
				if name == importCompiler.packageName {
					return &Package{
						compiler: importCompiler,
						Type:     NewType(vm.BasicTypePackage),
					}
				}
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

	for _, exe := range c.compile() {
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
	c.AddNativeFunctionPrint()
	c.AddNativeFunctionItoa()
}

func (c *Compiler) AddNativeFunctionPrint() {
	paramsType := []*Parameter{{Type: NewType(vm.BasicTypeString), Name: "str"}}
	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, nil),
		Name:            "print",
		PackageName:     "_sys",
		ParameterList:   paramsType,
		Block:           nil,
		DeclarationList: nil,
	}

	c.funcList = append(c.funcList, fd)
}

func (c *Compiler) AddNativeFunctionItoa() {
	paramsType := []*Parameter{{Type: NewType(vm.BasicTypeInt), Name: "int"}}
	resultsType := []*Parameter{{Type: NewType(vm.BasicTypeString)}}
	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, resultsType),
		Name:            "itoa",
		PackageName:     "_sys",
		ParameterList:   paramsType,
		Block:           nil,
		DeclarationList: nil,
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
	decl.Index = len(c.declarationList) - 1

	return decl.Index
}

func (c *Compiler) SearchDeclaration(name string) *Declaration {
	for _, declaration := range c.declarationList {
		if declaration.Name == name {
			return declaration
		}
	}
	return nil
}

func AddDeclList(decl *Declaration) {
	c := GetCurrentCompiler()
	decl.PackageName = c.packageName
	c.AddDeclarationList(decl)
}

func SetPackageName(packageName string) {
	c := GetCurrentCompiler()
	c.packageName = packageName
}

func SetImportList(importList []*Import) {
	c := GetCurrentCompiler()
	c.importList = importList
}

func PushCurrentBlock() *Block {
	c := GetCurrentCompiler()
	c.currentBlock = &Block{outerBlock: c.currentBlock}

	return c.currentBlock
}

func PopCurrentBlock() *Block {
	c := GetCurrentCompiler()

	b := c.currentBlock

	c.currentBlock = c.currentBlock.outerBlock

	return b
}
