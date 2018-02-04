package vm

//
// 字节码解释器
//

type Executable struct {
	// 常量池
	ConstantPool []Constant

	// 全局变量
	GlobalVariableList []*VmVariable

	// 函数列表
	FunctionList []*VmFunction

	// 顶层结构代码
	CodeList []byte

	// 行号对应表
	// 保存字节码和与之对应的源代码的行号
	LineNumberList []*VmLineNumber
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

func (c *ConstantInt) getInt() int {
	return c.intValue
}

func (c *ConstantInt) SetInt(i int) {
	c.intValue = i
}

type ConstantDouble struct {
	ConstantImpl
	doubleValue float64
}

func (c *ConstantDouble) getDouble() float64 {
	return c.doubleValue
}

func (c *ConstantDouble) SetDouble(value float64) {
	c.doubleValue = value
}

type ConstantString struct {
	ConstantImpl
	stringValue string
}

func (c *ConstantString) getString() string {
	return c.stringValue
}

func (c *ConstantString) SetString(value string) {
	c.stringValue = value
}

//
//
//
type VmTypeDerive interface{}

type VmFunctionDerive struct {
	ParameterList []*VmLocalVariable
}
type VmTypeSpecifier struct {
	BasicType BasicType
	DeriveList    []VmTypeDerive
}

// ==============================
// 全局变量
// ==============================

type VmVariable struct {
	name          string
	typeSpecifier *VmTypeSpecifier
}

func NewVmVariable(name string, typeSpecifier *VmTypeSpecifier) *VmVariable{
	return &VmVariable{
		name: name,
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

func NewExecutable() *Executable {
	exe := &Executable{}
	return exe
}
