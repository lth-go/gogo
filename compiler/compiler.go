package compiler

import (
	"log"

	"github.com/lth-go/gogo/vm"
)

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

func (c *Compiler) GetPackageName() string {
	return c.packageName
}

func (c *Compiler) SetPackageName(packageName string) {
	c.packageName = packageName
}

func (c *Compiler) Parse() {
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}
}

func NewCompiler(path string) *Compiler {
	c := &Compiler{
		lexer:           NewLexer(path),
		path:            path,
		importList:      []*Import{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		ConstantList:    []interface{}{},
	}

	return c
}

//
// 修正树
//
func (c *Compiler) FixTree() {
	// 修正全局变量
	c.FixDeclarationList()

	// 添加原生函数
	c.AddNativeFunctionList()

	// 修正函数
	for _, fd := range c.funcList {
		fd.Fix()
	}
}

//
// 生成字节码
//
func (c *Compiler) Generate() *vm.Executable {
	exe := vm.NewExecutable()
	exe.PackageName = c.GetPackageName()
	exe.FunctionList = c.GetVmFunctionList()
	exe.VariableList.SetVariableList(c.GetVmVariableList()) // 添加全局变量声明
	exe.ConstantPool.SetPool(c.ConstantList)

	return exe
}

func (c *Compiler) FixDeclarationList() {
	for _, decl := range c.declarationList {
		if decl.Value != nil {
			decl.Value = decl.Value.Fix()
			decl.Value = CreateAssignCast(decl.Value, decl.Type)
		}
	}
}

// 添加引用包函数
func (c *Compiler) GetFuncIndex(fd *FunctionDefinition) int {
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

func (c *Compiler) SearchFunction(name string) *FunctionDefinition {
	for _, func_ := range c.funcList {
		if func_.Name == name {
			return func_
		}
	}
	return nil
}

// TODO: remove
func (c *Compiler) SearchPackage(packageName string) *Package {
	for _, imp := range c.importList {
		if packageName != imp.packageName {
			continue
		}

		impCompiler := GetDoneCompiler(packageName)
		if impCompiler == nil {
			panic("TODO")
		}

		return &Package{
			compiler: impCompiler,
			Type:     NewType(vm.BasicTypePackage),
		}
	}

	return nil
}
