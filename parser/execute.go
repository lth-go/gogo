package parser

type Executable struct {
	// 常量池
	constantPool []Constant

	// 全局变量
	globalVariableList []*Variable

	// 函数列表
	functionList []*Function

	// 顶层结构代码
	codeList []byte

	// 行号对应表
	// 保存字节码和与之对应的源代码的行号
	lineNumberList []*LineNumber
}

// ==============================
// 常量池
// ==============================

type Constant interface{}

type ConstantType int

const (
	ConstantNumberType ConstantType = iota
	ConstantStringType
)

type ConstantNumber struct {
	numberValue float64
}

type ConstantString struct {
	stringValue string
}

// ==============================
// 全局变量
// ==============================

type Variable struct {
	name          string
	typeSpecifier *TypeSpecifier
}

// ==============================
// 函数
// ==============================

type Function struct {
	// 类型
	typeSpecifier *TypeSpecifier
	// 函数名
	name string
	// 形参列表
	parameterList []*LocalVariable
	// 是否原生函数
	isImplemented bool
	// 局部变量列表
	localVariable []*LocalVariable
	// 字节码类表
	codeList []byte

	// 行号对应表
	lineNumberList []*LineNumber
}

type LocalVariable struct {
	name          string
	typeSpecifier TypeSpecifier
}

// ==============================
// 行号对应表
// ==============================

type LineNumber struct {
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
type Opcode int

const (
	PUSH_INT_1BYTE Opcode = iota
	PUSH_INT_2BYTE
	PUSH_INT
	PUSH_DOUBLE_0
	PUSH_DOUBLE_1
	PUSH_DOUBLE
	PUSH_STRING
	/**********/
	PUSH_STACK_INT
	PUSH_STACK_DOUBLE
	PUSH_STACK_STRING
	POP_STACK_INT
	POP_STACK_DOUBLE
	POP_STACK_STRING
	/**********/
	PUSH_STATIC_INT
	PUSH_STATIC_DOUBLE
	PUSH_STATIC_STRING
	POP_STATIC_INT
	POP_STATIC_DOUBLE
	POP_STATIC_STRING
	/**********/
	ADD_INT
	ADD_DOUBLE
	ADD_STRING
	SUB_INT
	SUB_DOUBLE
	MUL_INT
	MUL_DOUBLE
	DIV_INT
	DIV_DOUBLE
	MOD_INT
	MOD_DOUBLE
	MINUS_INT
	MINUS_DOUBLE
	INCREMENT
	DECREMENT
	CAST_INT_TO_DOUBLE
	CAST_DOUBLE_TO_INT
	CAST_BOOLEAN_TO_STRING
	CAST_INT_TO_STRING
	CAST_DOUBLE_TO_STRING
	EQ_INT
	EQ_DOUBLE
	EQ_STRING
	GT_INT
	GT_DOUBLE
	GT_STRING
	GE_INT
	GE_DOUBLE
	GE_STRING
	LT_INT
	LT_DOUBLE
	LT_STRING
	LE_INT
	LE_DOUBLE
	LE_STRING
	NE_INT
	NE_DOUBLE
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
