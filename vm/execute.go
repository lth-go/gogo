package vm

import (
	"fmt"
)

//
// 字节码解释器
//
type Executable struct {
	PackageName    string        // 包名
	Path           string        // 源码路径
	ConstantPool   ConstantPool  // 常量池
	VariableList   *VariableList // 全局变量
	FunctionList   []*Function   // 函数列表
	CodeList       []byte        // 顶层结构代码
	LineNumberList []*LineNumber // 行号对应表,保存字节码和与之对应的源代码的行号
}

func NewExecutable() *Executable {
	exe := &Executable{
		ConstantPool:   NewConstantPool(),
		VariableList:   NewVmVariableList(),
		FunctionList:   []*Function{},
		CodeList:       []byte{},
		LineNumberList: []*LineNumber{},
	}

	return exe
}

func (exe *Executable) ShowCode() {
	for i := 0; i < len(exe.CodeList); {
		code := exe.CodeList[i]
		info := OpcodeInfo[int(code)]
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

//
// ExecutableList
//
type ExecutableList struct {
	List []*Executable
}

func NewExecutableList() *ExecutableList {
	return &ExecutableList{}
}

func (exeList *ExecutableList) AddExe(exe *Executable) bool {
	for _, itemExe := range exeList.List {
		if itemExe.PackageName == exe.PackageName {
			return false
		}
	}

	exeList.List = append(exeList.List, exe)
	return true
}

func (exeList *ExecutableList) GetTopExe() *Executable {
	return exeList.List[len(exeList.List)-1]
}

//
// 全局变量
//
type VariableList struct {
	VariableList []*Variable
}

func (vl *VariableList) Append(v *Variable) {
	vl.VariableList = append(vl.VariableList, v)
}

func (vl *VariableList) Init() {
	for _, value := range vl.VariableList {
		value.Init()
	}
}

func (vl *VariableList) getInt(index int) int {
	return vl.VariableList[index].Value.(int)
}

func (vl *VariableList) getDouble(index int) float64 {
	return vl.VariableList[index].Value.(float64)
}

func (vl *VariableList) getObject(index int) Object {
	return vl.VariableList[index].Value.(Object)
}

func (vl *VariableList) setInt(index int, value int) {
	vl.VariableList[index].Value = value
}

func (vl *VariableList) setDouble(index int, value float64) {
	vl.VariableList[index].Value = value
}

func (vl *VariableList) setObject(index int, value Object) {
	vl.VariableList[index].Value = value
}

func NewVmVariableList() *VariableList {
	return &VariableList{
		VariableList: []*Variable{},
	}
}

type Variable struct {
	Name  string
	Type  *Type
	Value interface{}
}

func (v *Variable) Init() {
	v.Value = initializeValue(v.Type)
}

func (v *Variable) IsReferenceType() bool {
	return v.Type.IsReferenceType()
}

func NewVmVariable(name string, typ *Type) *Variable {
	return &Variable{
		Name: name,
		Type: typ,
	}
}

// ==============================
// 函数
// ==============================
type Function struct {
	Type              *Type         // 类型
	PackageName       string        // 包名
	Name              string        // 函数名
	IsImplemented     bool          // 是否在当前包实现
	IsMethod          bool          // 是否是方法
	ParameterList     []*Variable   // 形参列表
	LocalVariableList []*Variable   // 局部变量列表
	CodeList          []byte        // 字节码类表
	LineNumberList    []*LineNumber // 行号对应表
	// Executable *Executable
}

func (f *Function) ShowCode() {
	for i := 0; i < len(f.CodeList); {
		code := f.CodeList[i]
		info := OpcodeInfo[int(code)]
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
