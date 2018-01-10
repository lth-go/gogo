package parser

// ExpressionKind 表达式类型
type ExpressionKind int

const (
	BOOLEAN_EXPRESSION ExpressionKind = iota
	NumberExpression
	STRING_EXPRESSION
	IDENTIFIER_EXPRESSION
	COMMA_EXPRESSION
	ASSIGN_EXPRESSION
	ADD_EXPRESSION
	SUB_EXPRESSION
	MUL_EXPRESSION
	DIV_EXPRESSION
	EQ_EXPRESSION
	NE_EXPRESSION
	GT_EXPRESSION
	GE_EXPRESSION
	LT_EXPRESSION
	LE_EXPRESSION
	LOGICAL_AND_EXPRESSION
	LOGICAL_OR_EXPRESSION
	LOGICAL_NOT_EXPRESSION
	FUNCTION_CALL_EXPRESSION
)

// Expression 表达式接口
type Expression interface {
	//Pos
}

// CommaExpression 逗号表达式
type CommaExpression struct {
	left  Expression
	right Expression
}

// AssignExpression 赋值表达式
type AssignExpression struct {
	// 左值
	left Expression
	// 符号
	operator AssignmentOperator
	// 操作数
	operand Expression
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
	operator ExpressionKind
	left     Expression
	right    Expression
}

// MinusExpression 负数表达式
type MinusExpression struct {
	operand Expression
}

// LogicalNotExpression 逻辑非表达式
type LogicalNotExpression struct {
	operand Expression
}

// FunctionCallExpression 函数调用表达式
type FunctionCallExpression struct {
	function Expression
	argument []Expression
}

// Boolean 布尔类型
type Boolean int

const (
	// BooleanTrue true
	BooleanTrue Boolean = iota
	// BooleanFalse false
	BooleanFalse
)

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	booleanValue Boolean
}

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	name string
}
