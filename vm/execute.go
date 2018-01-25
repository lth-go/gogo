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

func (c *ConstantImpl) getString() float64 {
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

func (c *ConstantString) getString() float64 {
	return c.stringValue
}

//
//
//
type VmTypeDerive interface{}

type VmFunctionDerive struct {
	parameter []*VmLocalVariable
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

// ==============================
// 字节码
// ==============================

// 字节码
const (
	PUSH_NUMBER_0 byte = iota
	PUSH_NUMBER_1
	PUSH_NUMBER
	PUSH_STRING
	/**********/
	PUSH_STACK_NUMBER
	PUSH_STACK_STRING
	POP_STACK_NUMBER
	POP_STACK_STRING
	/**********/
	PUSH_STATIC_NUMBER
	PUSH_STATIC_STRING
	POP_STATIC_NUMBER
	POP_STATIC_STRING
	/**********/
	ADD_NUMBER
	ADD_STRING
	SUB_NUMBER
	MUL_NUMBER
	DIV_NUMBER
	MOD_NUMBER
	MINUS_NUMBER
	INCREMENT
	DECREMENT
	CAST_BOOLEAN_TO_STRING
	CAST_NUMBER_TO_STRING
	EQ_NUMBER
	EQ_STRING
	GT_NUMBER
	GT_STRING
	GE_NUMBER
	GE_STRING
	LT_NUMBER
	LT_STRING
	LE_NUMBER
	LE_STRING
	NE_NUMBER
	NE_STRING
	LOGICAL_AND
	LOGICAL_OR
	LOGICAL_NOT
	POP
	DUPLICATE
	JUMP
	JUMP_IF_TRUE
	JUMP_IF_FALSE
	/**********/
	PUSH_FUNCTION
	INVOKE
	RETURN
)

type OpcodeInfo struct {
	// 注记符
	mnemonic string

	// 参数类型，
	// `b` 一个字节整数
	// `s` 两个字节整数
	// `p` 常量池索引值
	parameter       string
	stack_increment int
}
