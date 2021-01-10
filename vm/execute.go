package vm

import (
	"fmt"
)

//
// 字节码解释器
//
type Executable struct {
	// 包名
	PackageName string
	// 是否是被导入的
	IsImported bool
	// 源码路径
	Path string

	// 常量池
	ConstantPool ConstantPool
	// 全局变量
	VariableList *VariableList
	// 函数列表
	FunctionList []*Function

	// 顶层结构代码
	CodeList []byte
	// 行号对应表,保存字节码和与之对应的源代码的行号
	LineNumberList []*LineNumber
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

func (exe *Executable) AddConstantPool(cp Constant) int {
	exe.ConstantPool.Append(cp)
	return exe.ConstantPool.Length() - 1
}

//
// ExecutableList
//
type ExecutableList struct {
	TopLevel *Executable
	List     []*Executable
}

func NewExecutableList() *ExecutableList {
	return &ExecutableList{}
}

func (exeList *ExecutableList) AddExe(exe *Executable) bool {
	for _, itemExe := range exeList.List {
		if itemExe.PackageName == exe.PackageName && itemExe.IsImported == exe.IsImported {
			return false
		}
	}

	exeList.List = append(exeList.List, exe)
	return true
}

// ==============================
// 常量池
// ==============================

type Constant interface {
	getInt() int
	getDouble() float64
	getString() string
}

type ConstantImpl struct{}

func (c *ConstantImpl) getInt() int {
	panic("error")
}

func (c *ConstantImpl) getDouble() float64 {
	panic("error")
}

func (c *ConstantImpl) getString() string {
	panic("error")
}

type ConstantInt struct {
	ConstantImpl
	intValue int
}

func NewConstantInt(value int) *ConstantInt {
	return &ConstantInt{intValue: value}
}

func (c *ConstantInt) getInt() int {
	return c.intValue
}

type ConstantDouble struct {
	ConstantImpl
	doubleValue float64
}

func NewConstantDouble(value float64) *ConstantDouble {
	return &ConstantDouble{doubleValue: value}
}

func (c *ConstantDouble) getDouble() float64 {
	return c.doubleValue
}

type ConstantString struct {
	ConstantImpl
	stringValue string
}

func NewConstantString(value string) *ConstantString {
	return &ConstantString{stringValue: value}
}

func (c *ConstantString) getString() string {
	return c.stringValue
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
	return vl.VariableList[index].Value.(*IntValue).intValue
}

func (vl *VariableList) getDouble(index int) float64 {
	return vl.VariableList[index].Value.(*DoubleValue).doubleValue
}

func (vl *VariableList) getObject(index int) *ObjectRef {
	return vl.VariableList[index].Value.(*ObjectRef)
}

func (vl *VariableList) setInt(index int, value int) {
	vl.VariableList[index].Value.(*IntValue).intValue = value
}

func (vl *VariableList) setDouble(index int, value float64) {
	vl.VariableList[index].Value.(*DoubleValue).doubleValue = value
}

func (vl *VariableList) setObject(index int, value *ObjectRef) {
	vl.VariableList[index].Value = value
}

func NewVmVariableList() *VariableList {
	return &VariableList{
		VariableList: []*Variable{},
	}
}

type Variable struct {
	Name          string
	TypeSpecifier *TypeSpecifier
	Value         Value
}

func (v *Variable) Init() {
	v.Value = initializeValue(v.TypeSpecifier)
}

func (v *Variable) IsReferenceType() bool {
	return v.TypeSpecifier.IsReferenceType()
}

func NewVmVariable(name string, typeSpecifier *TypeSpecifier) *Variable {
	return &Variable{
		Name:          name,
		TypeSpecifier: typeSpecifier,
	}
}

// ==============================
// 函数
// ==============================
type Function struct {
	// 类型
	TypeSpecifier *TypeSpecifier
	// 包名
	PackageName string
	// 函数名
	Name string
	// 是否原生函数
	IsNative bool
	// 是否是方法
	IsMethod bool
	// 形参列表
	ParameterList []*Variable
	// 局部变量列表
	LocalVariableList []*Variable
	// 字节码类表
	CodeList []byte
	// 行号对应表
	LineNumberList []*LineNumber
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
