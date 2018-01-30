package vm

//
// 字节码解释器
//

type Executable struct {
	// 常量池
	constantPool []Constant

	// 全局变量
	globalVariableList []*VmVariable

	// 函数列表
	functionList []*VmFunction

	// 顶层结构代码
	codeList []byte

	// 行号对应表
	// 保存字节码和与之对应的源代码的行号
	lineNumberList []*VmLineNumber
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
	intValue int
}

func (c *ConstantInt) getInt() int {
	return c.intValue
}

type ConstantDouble struct {
	doubleValue float64
}

func (c *ConstantDouble) getDouble() float64 {
	return c.doubleValue
}

type ConstantString struct {
	stringValue string
}

func (c *ConstantString) getString() string {
	return c.stringValue
}

//
//
//
type VmTypeDerive interface{}

type VmFunctionDerive struct {
	parameterList []*VmLocalVariable
}
type VmTypeSpecifier struct {
	basicType BasicType
	derive    []VmTypeDerive
}

// ==============================
// 全局变量
// ==============================

type VmVariable struct {
	name          string
	typeSpecifier *VmTypeSpecifier
}

// ==============================
// 函数
// ==============================

type VmFunction struct {
	// 类型
	typeSpecifier *VmTypeSpecifier
	// 函数名
	name string
	// 形参列表
	parameterList []*VmLocalVariable
	// 是否原生函数
	isImplemented bool
	// 局部变量列表
	localVariableList []*VmLocalVariable
	// 字节码类表
	codeList []byte

	// 行号对应表
	lineNumberList []*VmLineNumber
}

type VmLocalVariable struct {
	name          string
	typeSpecifier *VmTypeSpecifier
}

// ==============================
// 行号对应表
// ==============================

type VmLineNumber struct {
	// 源代码行号
	lineNumber int

	// 字节码开始的位置
	startPc int

	// 接下来有多少字节码对应相同的行号
	pcCount int
}

func newExecutable() *Executable {
	exe := &Executable{}
	return exe
}
