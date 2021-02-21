package compiler

import (
	"github.com/lth-go/gogo/utils"
	"github.com/lth-go/gogo/vm"
)

// TODO:
var CurrentPackageName string

type CompilerManager struct {
	doingList []*Compiler
	doneList  []*Compiler

	funcList        []*FunctionDefinition // 函数列表
	DeclarationList []*Declaration        // 声明列表
	ConstantList    []interface{}         // 常量定义

	CodeList []byte
}

func (cm *CompilerManager) AddConstantList(value interface{}) int {
	for i, v := range cm.ConstantList {
		if value == v {
			return i
		}
	}

	cm.ConstantList = append(cm.ConstantList, value)
	return len(cm.ConstantList) - 1
}

// 添加引用包函数
func (cm *CompilerManager) AddFuncList(fd *FunctionDefinition) int {
	packageName := fd.GetPackageName()
	name := fd.GetName()

	for i, f := range cm.funcList {
		if packageName == f.GetPackageName() && name == f.GetName() {
			return i
		}
	}

	return -1
}

func (cm *CompilerManager) SearchFunction(packageName string, name string) *FunctionDefinition {
	for _, f := range cm.funcList {
		if f.PackageName == packageName && f.Name == name {
			return f
		}
		if f.PackageName == "_sys" && f.Name == name {
			return f
		}
	}
	return nil
}

func (cm *CompilerManager) SearchDeclaration(packageName string, name string) *Declaration {
	for _, decl := range cm.DeclarationList {
		if decl.PackageName == packageName && decl.Name == name {
			return decl
		}
	}
	return nil
}

var compilerManager *CompilerManager

func NewCompilerManager() *CompilerManager {
	compilerManager = &CompilerManager{
		doingList: []*Compiler{},
		doneList:  []*Compiler{},
	}

	return compilerManager
}

func GetCurrentCompilerManage() *CompilerManager {
	return compilerManager
}

// GetCurrentCompiler
func GetCurrentCompiler() *Compiler {
	length := len(compilerManager.doingList)
	if length == 0 {
		// TODO: 临时处理
		if CurrentPackageName != "" {
			for _, c := range compilerManager.doneList {
				if c.GetPackageName() == CurrentPackageName {
					return c
				}
			}
		}
		return nil
	}

	return compilerManager.doingList[length-1]
}

func (cm *CompilerManager) PushCurrentCompiler(c *Compiler) {
	cm.doingList = append(cm.doingList, c)
}

func (cm *CompilerManager) PopCurrentCompiler() {
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

func (cm *CompilerManager) AddDoneCompilerList(c *Compiler) {
	cm.doneList = append(cm.doneList, c)
}

func GetDoneCompilerList() []*Compiler {
	return compilerManager.doneList
}

func (cm *CompilerManager) GetDoneCompiler(packageName string) *Compiler {
	for _, c := range GetDoneCompilerList() {
		if c.GetPackageName() == packageName {
			return c
		}
	}
	return nil
}

func (cm *CompilerManager) Parse(path string) {
	c := NewCompiler(path)

	cm.PushCurrentCompiler(c)
	cm.AddDoneCompilerList(c)

	// 生成语法树
	c.Parse()

	for _, imp := range c.importList {
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

func (cm *CompilerManager) Compile() {
	// 添加原生函数声明
	cm.AddNativeFunctionList()

	for i := len(cm.doneList) - 1; i >= 0; i-- {
		c := cm.doneList[i]

		// 添加并修正全局声明
		cm.DeclarationList = append(cm.DeclarationList, c.declarationList...)
		for index, decl := range cm.DeclarationList {
			decl.Index = index
			decl.Value = decl.Value.Fix()
			decl.Value = CreateAssignCast(decl.Value, decl.Type)
		}

		// 添加函数
		cm.funcList = append(cm.funcList, c.funcList...)
	}

	//
	// TODO: 移除CurrentPackageName依赖
	//
	for _, f := range cm.funcList {
		CurrentPackageName = f.PackageName
		f.Fix()
	}

	//
	// 函数生成字节码,并修正字节码
	//
	for _, f := range cm.funcList {
		if f.Block == nil {
			continue
		}

		ob := NewOpCodeBuf()
		generateStatementList(f.Block.statementList, ob)

		//
		// 修正Label
		//
		codeList := ob.FixLabel()

		//
		// 修正opCode
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

	cm.SetCodeList()
}

func (cm *CompilerManager) SetCodeList() {
	mainFunc := -1

	for i, f := range cm.funcList {
		if f.PackageName == "main" && f.Name == "main" {
			mainFunc = i
		}
	}

	if mainFunc == -1 {
		panic("TODO")
	}

	b := make([]byte, 2)
	utils.Set2ByteInt(b, mainFunc)
	cm.CodeList = append(cm.CodeList, b...)
	cm.CodeList = append(cm.CodeList, vm.OP_CODE_INVOKE)
}

//
// 编译文件
//
func (cm *CompilerManager) CompileFile(path string) {
	// 输出yacc错误信息
	if true {
		yyErrorVerbose = true
	}

	cm.Parse(path)

	cm.Compile()
}
