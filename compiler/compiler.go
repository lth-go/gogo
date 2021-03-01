package compiler

import (
	"github.com/lth-go/gogo/utils"
	"github.com/lth-go/gogo/vm"
)

type Compiler struct {
	doingList []*Package
	doneList  []*Package

	FuncList        []*FunctionDefinition // 函数列表
	DeclarationList []*Declaration        // 声明列表
	ConstantList    []interface{}         // 常量定义

	CodeList []byte
}

func (cm *Compiler) AddConstant(value interface{}) int {
	for i, v := range cm.ConstantList {
		if value == v {
			return i
		}
	}

	cm.ConstantList = append(cm.ConstantList, value)

	return len(cm.ConstantList) - 1
}

func (cm *Compiler) SearchFunction(packageName string, name string) (*FunctionDefinition, int) {
	for i, f := range cm.FuncList {
		if f.PackageName == packageName && f.Name == name {
			return f, i
		}

		if f.PackageName == "_sys" && f.Name == name {
			return f, i
		}
	}

	return nil, -1
}

func (cm *Compiler) SearchDeclaration(packageName string, name string) *Declaration {
	for _, decl := range cm.DeclarationList {
		if decl.PackageName == packageName && decl.Name == name {
			return decl
		}
	}

	return nil
}

var compilerManager *Compiler

func NewCompilerManager() *Compiler {
	compilerManager = &Compiler{
		doingList: []*Package{},
		doneList:  []*Package{},
	}

	return compilerManager
}

func GetCurrentCompiler() *Compiler {
	return compilerManager
}

// GetCurrentPackage
func GetCurrentPackage() *Package {
	length := len(compilerManager.doingList)
	if length == 0 {
		return nil
	}

	return compilerManager.doingList[length-1]
}

func (cm *Compiler) PushCurrentCompiler(c *Package) {
	cm.doingList = append(cm.doingList, c)
}

func (cm *Compiler) PopCurrentCompiler() {
	cm.doingList = cm.doingList[:len(cm.doingList)-1]
}

func IsCompiling(packageName string) bool {
	for _, c := range compilerManager.doingList {
		if c.GetPackageName() == packageName {
			return true
		}
	}

	return false
}

func (cm *Compiler) AddDoneCompilerList(c *Package) {
	cm.doneList = append(cm.doneList, c)
}

func (cm *Compiler) GetDoneCompiler(packageName string) *Package {
	for _, c := range cm.doneList {
		if c.GetPackageName() == packageName {
			return c
		}
	}
	return nil
}

func (cm *Compiler) Parse(path string) {
	c := NewPackage(path)

	cm.PushCurrentCompiler(c)
	cm.AddDoneCompilerList(c)

	// 生成语法树
	c.Parse()

	for _, imp := range c.importDeclList {
		if IsCompiling(imp.packageName) {
			panic("TODO")
		}

		// 判断是否已经被解析过
		if cm.GetDoneCompiler(imp.packageName) != nil {
			continue
		}

		cm.Parse(imp.GetPath())
	}

	cm.PopCurrentCompiler()
}

func (c *Compiler) Compile() {
	// 添加原生函数声明
	c.AddNativeFunctionList()

	for i := len(c.doneList) - 1; i >= 0; i-- {
		pkg := c.doneList[i]

		// 添加并修正全局声明
		c.DeclarationList = append(c.DeclarationList, pkg.declarationList...)
		for index, decl := range c.DeclarationList {
			decl.Index = index
			decl.Value = decl.Value.Fix()
			decl.Value = CreateAssignCast(decl.Value, decl.Type)
		}

		// 添加函数
		c.FuncList = append(c.FuncList, pkg.funcList...)
	}

	for _, f := range c.FuncList {
		// 将函数所在的compiler压栈
		for _, pkg := range c.doneList {
			if pkg.GetPackageName() == f.PackageName {
				c.PushCurrentCompiler(pkg)
			}
		}

		f.Fix()
	}

	//
	// 函数生成字节码,并修正字节码
	//
	for _, f := range c.FuncList {
		if f.Block == nil {
			continue
		}

		ob := NewOpCodeBuf()
		for _, stmt := range f.Block.statementList {
			stmt.Generate(ob)
		}

		//
		// 修正Label
		//
		codeList := ob.FixLabel()

		//
		// 修正opCode
		// TODO: 去掉代码
		//
		paramCount := len(f.GetType().funcType.Params)

		for i := 0; i < len(codeList); i++ {
			code := codeList[i]
			switch code {
			// 函数内的本地声明
			case vm.OP_CODE_PUSH_STACK, vm.OP_CODE_POP_STACK:
				// 形参
				// 返回值(新增)
				// 声明

				// 增加返回值的位置
				idx := utils.Get2ByteInt(codeList[i+1:])
				if idx >= paramCount {
					utils.Set2ByteInt(codeList[i+1:], idx+1)
				} else {
					utils.Set2ByteInt(codeList[i+1:], idx-paramCount)
				}
			}

			for _, p := range []byte(vm.OpcodeInfo[code].Parameter) {
				switch p {
				case 'b':
					i++
				case 's', 'p':
					i += 2
				default:
					panic("TODO")
				}
			}
		}

		f.CodeList = codeList
	}

	c.SetCodeList()
}

func (c *Compiler) SetCodeList() {
	mainFunc := -1

	for i, f := range c.FuncList {
		if f.PackageName == "main" && f.Name == "main" {
			mainFunc = i
		}
	}

	if mainFunc == -1 {
		panic("TODO")
	}

	b := make([]byte, 2)
	utils.Set2ByteInt(b, mainFunc)
	c.CodeList = append(c.CodeList, b...)
	c.CodeList = append(c.CodeList, vm.OP_CODE_INVOKE)
}

//
// 编译文件
//
func (c *Compiler) CompileFile(path string) {
	// 输出yacc错误信息
	if true {
		yyErrorVerbose = true
	}

	c.Parse(path)

	c.Compile()
}
