package vm

//
// 字节码解释器
//

type Executable struct {
	// 常量池
	ConstantPool ConstantPool

	// 全局变量
	// 仅保存名称和类型
	GlobalVariableList []*VmVariable

	// 函数列表
	FunctionList []*VmFunction

	// 顶层结构代码
	CodeList []byte

	// 行号对应表
	// 保存字节码和与之对应的源代码的行号
	LineNumberList []*VmLineNumber

	TypeSpecifierList []*VmTypeSpecifier
}

func NewExecutable() *Executable {
	exe := &Executable{
		ConstantPool: NewConstantPool(),
		GlobalVariableList: []*VmVariable{},
		FunctionList: []*VmFunction{},
		CodeList: []byte{},
		LineNumberList: []*VmLineNumber{},
		TypeSpecifierList: []*VmTypeSpecifier{},
	}
	return exe
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
//
//
type VmTypeDerive interface{}

type VmFunctionDerive struct {
	ParameterList []*VmLocalVariable
}

type VmArrayDerive struct {
}

type VmTypeSpecifier struct {
	BasicType  BasicType
	DeriveList []VmTypeDerive
}

func (t *VmTypeSpecifier) AppendDerive(derive VmTypeDerive) {
	if t.DeriveList == nil {
		t.DeriveList = []VmTypeDerive{}
	}
	t.DeriveList = append(t.DeriveList, derive)
}

// ==============================
// 全局变量
// ==============================

type VmVariable struct {
	name          string
	typeSpecifier *VmTypeSpecifier
}

func NewVmVariable(name string, typeSpecifier *VmTypeSpecifier) *VmVariable {
	return &VmVariable{
		name:          name,
		typeSpecifier: typeSpecifier,
	}
}

// ==============================
// 函数
// ==============================

type VmFunction struct {
	// 类型
	TypeSpecifier *VmTypeSpecifier
	// 函数名
	Name string
	// 形参列表
	ParameterList []*VmLocalVariable
	// 是否原生函数
	IsImplemented bool
	// 局部变量列表
	LocalVariableList []*VmLocalVariable
	// 字节码类表
	CodeList []byte

	// 行号对应表
	LineNumberList []*VmLineNumber
}

type VmLocalVariable struct {
	Name          string
	TypeSpecifier *VmTypeSpecifier
}

// ==============================
// 行号对应表
// ==============================

type VmLineNumber struct {
	// 源代码行号
	LineNumber int

	// 字节码开始的位置
	StartPc int

	// 接下来有多少字节码对应相同的行号
	PcCount int
}
