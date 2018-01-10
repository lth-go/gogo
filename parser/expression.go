package parser

type ExpressionKind int

const (
	BOOLEAN_EXPRESSION ExpressionKind = iota
	INT_EXPRESSION
	DOUBLE_EXPRESSION
	STRING_EXPRESSION
	IDENTIFIER_EXPRESSION
	COMMA_EXPRESSION
	ASSIGN_EXPRESSION
	ADD_EXPRESSION
	SUB_EXPRESSION
	MUL_EXPRESSION
	DIV_EXPRESSION
	MOD_EXPRESSION
	EQ_EXPRESSION
	NE_EXPRESSION
	GT_EXPRESSION
	GE_EXPRESSION
	LT_EXPRESSION
	LE_EXPRESSION
	LOGICAL_AND_EXPRESSION
	LOGICAL_OR_EXPRESSION
	MINUS_EXPRESSION
	LOGICAL_NOT_EXPRESSION
	FUNCTION_CALL_EXPRESSION
	EXPRESSION_KIND_COUNT_PLUS_1
)

type Expression interface {
	//Pos
}

type CommaExpression struct {
	left  Expression
	right Expression
}

type AssignExpression struct {
	// 左值
	left Expression
	// 符号
	operator AssignmentOperator
	// 操作数
	operand Expression
}

type BinaryExpression struct {
	operator ExpressionKind
	left     Expression
	right    Expression
}

type MinusExpression struct {
	operand Expression
}

type LogicalNotExpression struct {
	operand Expression
}

type FunctionCallExpression struct {
	function Expression
	argument []Expression
}

type Boolean int

const (
	BOOLEAN_TRUE Boolean = iota
	BOOLEAN_FALSE
)

type BooleanExpression struct {
	boolean_value Boolean
}
