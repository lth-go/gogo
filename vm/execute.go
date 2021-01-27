package vm

import (
	"fmt"
)

//
// 字节码解释器
//
type Executable struct {
	PackageName  string       // 包名
	ConstantPool ConstantPool // 常量池
	VariableList []*Variable  // 全局变量
	FunctionList []*Function  // 函数列表
	CodeList     []byte       // 顶层结构代码
}

func NewExecutable() *Executable {
	exe := &Executable{
		ConstantPool: NewConstantPool(),
		VariableList: []*Variable{},
		FunctionList: []*Function{},
		CodeList:     []byte{},
	}

	return exe
}

//
// ExecutableList
//
type ExecutableList struct {
	List []*Executable
}

func NewExecutableList(exeList []*Executable) *ExecutableList {
	l := &ExecutableList{}

	for _, exe := range exeList {
		l.Add(exe)
	}

	return l
}

func (exeList *ExecutableList) Add(exe *Executable) {
	for _, itemExe := range exeList.List {
		if itemExe.PackageName == exe.PackageName {
			return
		}
	}

	exeList.List = append(exeList.List, exe)
}

func (exeList *ExecutableList) Top() *Executable {
	return exeList.List[len(exeList.List)-1]
}

//
// Variable 全局变量
//
type Variable struct {
	PackageName string
	Name        string
	Type        *Type
	Value       interface{}
}

func (v *Variable) Init() {
	if v.Value != nil {
		return
	}

	var value interface{}

	if v.Type.IsReferenceType() {
		value = NilObject
	} else {
		switch v.Type.GetBasicType() {
		case BasicTypeBool, BasicTypeInt:
			value = 0
		case BasicTypeFloat:
			value = 0.0
		case BasicTypeString:
			value = ""
		case BasicTypeNil:
			fallthrough
		default:
			panic("TODO")
		}
	}

	v.Value = value
}

func NewVmVariable(packageName string, name string, typ *Type) *Variable {
	return &Variable{
		PackageName: packageName,
		Name:        name,
		Type:        typ,
	}
}

//
// Function 函数
//
type Function struct {
	Type           *Type         // 类型
	PackageName    string        // 包名
	Name           string        // 函数名
	IsImplemented  bool          // 是否在当前包实现
	IsMethod       bool          // 是否是方法
	VariableList   []*Variable   // 局部变量列表
	CodeList       []byte        // 字节码类表
	LineNumberList []*LineNumber // 行号对应表
}

func (f *Function) ShowCode() {
	for i := 0; i < len(f.CodeList); {
		code := f.CodeList[i]
		info := OpcodeInfo[code]
		paramList := []byte(info.Parameter)

		fmt.Println(info.Mnemonic)
		for _, param := range paramList {
			switch param {
			case 'b':
				i += 1
			case 's', 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
		i += 1
	}
}

func (f *Function) GetParamCount() int {
	return len(f.Type.FuncType.ParamTypeList)
}

func (f *Function) GetResultCount() int {
	return len(f.Type.FuncType.ResultTypeList)
}

//
// 行号对应表
//
type LineNumber struct {
	// 源代码行号
	LineNumber int
	// 字节码开始的位置
	StartPc int
	// 接下来有多少字节码对应相同的行号
	PcCount int
}
