package parser

// ==============================
// BinaryOperatorKind
// ==============================

type BinaryOperatorKind int

const (
	LogicalOrOperator BinaryOperatorKind = iota
	LogicalAndOperator
	EqOperator
	NeOperator
	GtOperator
	GeOperator
	LtOperator
	LeOperator
	AddOperator
	SubOperator
	MulOperator
	DivOperator
)

// Expression 表达式接口
type Expression interface {
	// Pos接口
	Pos
	expr()

	fix(*Block)

	typeS() *TypeSpecifier
}

// ExpressionImpl provide commonly implementations for Expr.
type ExpressionImpl struct {
	PosImpl // ExprImpl provide Pos() function.

	typeSpecifier *TypeSpecifier
	lineNumber    int
}

// expr provide restraint interface.
func (e *ExpressionImpl) expr() {}

func (e *ExpressionImpl) typeS() *TypeSpecifier {
	return e.typeSpecifier

}

// ==============================
// CommaExpression
// ==============================

// CommaExpression 逗号表达式
type CommaExpression struct {
	ExpressionImpl

	left  Expression
	right Expression
}

func (e *CommaExpression) fix(currentBlock *Block) {

	e.left.fix(currentBlock)
	e.right.fix(currentBlock)
	e.typeSpecifier = e.right.typeS()

}

// ==============================
// AssignExpression
// ==============================

// AssignExpression 赋值表达式
type AssignExpression struct {
	ExpressionImpl
	// 左值
	left Expression
	// 符号
	operator AssignmentOperator
	// 操作数
	operand Expression
}

func (e *AssignExpression) fix(currentBlock *Block) {
	t, ok := e.left.(*IdentifierExpression)
	if !ok {
		compileError(e.left.lineNumber, 0, "")
	}

	e.left.fix(currentBlock)
	e.operand.fix(currentBlock)
	createAssignCast(e.operand, e.left.typeS())
	e.typeSpecifier = e.left.typeS()

}

// ==============================
// BinaryExpression
// ==============================

// BinaryExpression 二元表达式
type BinaryExpression struct {
	ExpressionImpl
	operator BinaryOperatorKind
	left     Expression
	right    Expression
}

func (e *BinaryExpression) fix(currentBlock *Block) {
	switch e.operator {
	case AddOperator, SubOperator, MulOperator, DivOperator:
		e.left.fix(currentBlock)
		e.right.fix(currentBlock)

		evalMathExpression(currentBlock, e)
		castBinaryExpression(e)

		if isNumber(e.left.typeS()) && isNumber(e.right.typeS()) {
			e.typeSpecifier = &TypeSpecifier{basicType: NumberType}
		} else if isString(e.left.typeS()) && isString(e.right.typeS()) {
			e.typeSpecifier = &TypeSpecifier{basicType: StringType}
		} else {
			compileError(e.lineNumber, 0, "")
		}

	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		e.left.fix(currentBlock)
		e.right.fix(currentBlock)

		evalCompareExpression(e)

		castBinaryExpression(e)

		if e.left.typeS() != e.right.typeS() {
			compileError(e.lineNumber, 0, "")
		}
		e.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

	case LogicalAndOperator, LogicalOrOperator:
		e.left.fix(currentBlock)
		e.right.fix(currentBlock)

		if isBoolean(e.left.typeS()) && isBoolean(e.right.typeS()) {
			e.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

		} else {

			compileError(e.lineNumber, 0, "")
		}
	}
}

// ==============================
// MinusExpression
// ==============================

// MinusExpression 负数表达式
type MinusExpression struct {
	ExpressionImpl
	operand Expression
}

func (e *MinusExpression) fix(currentBlock *Block) {
	// TODO 是否能去掉
	e.operand.fix(currentBlock)

	n, ok := e.operand.(*NumberExpression)
	if !ok {
		compileError(e.lineNumber, 0, "")
	}

	e.typeSpecifier = e.operand.typeS()

	n.numberValue = -n.numberValue
}

// ==============================
// LogicalNotExpression
// ==============================

// LogicalNotExpression 逻辑非表达式
type LogicalNotExpression struct {
	ExpressionImpl
	operand Expression
}

func (e *LogicalNotExpression) fix(currentBlock *Block) {
	e.operand.fix(currentBlock)

	b, ok := e.operand.(*BooleanExpression)
	if !ok {
		compileError(e.lineNumber, 0, "")
	}

	b.booleanValue = !b.booleanValue

	e.typeSpecifier = e.operand.typeS()

}

// ==============================
// FunctionCallExpression
// ==============================

// FunctionCallExpression 函数调用表达式
type FunctionCallExpression struct {
	ExpressionImpl
	function Expression
	argument []Expression

	name string
}

func (e *FunctionCallExpression) fix(currentBlock *Block) {
	e.function.fix(currentBlock)

	fd := searchFunction(e.name)
	if fd == nil {
		compileError(e.lineNumber, 0, "")

		checkArgument(currentBlock, fd, e)

		e.typeSpecifier = &TypeSpecifier{basicType: fd.typeS()}
		// TODO
		//expr->type->derive = fd->type->derive;
	}
}

// Boolean 布尔类型
type Boolean int

const (
	// BooleanTrue true
	BooleanTrue Boolean = iota
	// BooleanFalse false
	BooleanFalse
)

// ==============================
// BooleanExpression
// ==============================

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	ExpressionImpl
	typeSpecifier *TypeSpecifier
	booleanValue  Boolean
}

func (e *BooleanExpression) fix(currentBlock *Block) {
	e.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
}

// ==============================
// NumberExpression
// ==============================

// NumberExpression 数字表达式
type NumberExpression struct {
	ExpressionImpl
	typeSpecifier *TypeSpecifier
	numberValue   float64
}

func (e *NumberExpression) fix(currentBlock *Block) {
	e.typeSpecifier = &TypeSpecifier{basicType: NumberType}

}

// ==============================
// StringExpression
// ==============================

// StringExpression 字符串表达式
type StringExpression struct {
	ExpressionImpl
	typeSpecifier *TypeSpecifier
	stringValue   string
}

func (e *StringExpression) fix(currentBlock *Block) {
	e.typeSpecifier = &TypeSpecifier{basicType: StringType}
}

// ==============================
// IdentifierExpression
// ==============================

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionImpl
	name string
}

func (e *IdentifierExpression) fix(currentBlock *Block) {

	decl := searchDeclaration(e.name, currentBlock)
	if decl != nil {
		e.typeSpecifier = decl.typeSpecifier
		//expr->u.identifier.is_function = DVM_FALSE;
		//expr->u.identifier.u.declaration = decl;
		return
	}
	fd := searchFunction(e.name)
	if fd == nil {
		compileError(e.lineNumber, 0, "")
	}
	e.typeSpecifier = fd.typeS()

	//expr->u.identifier.is_function = DVM_TRUE;
	//expr->u.identifier.u.function = fd;

}

// ==============================
// utils
// ==============================

func isNumber(t *TypeSpecifier) bool {
	return t.basicType == NumberType
}
func isBoolean(t *TypeSpecifier) bool {
	return t.basicType == BooleanType
}
func isString(t *TypeSpecifier) bool {
	return t.basicType == StringType
}

func createAssignCast(e Expression, t *TypeSpecifier) {

}

func evalMathExpression(currentBlock *Block, e Expression) {

}

func castBinaryExpression(e *Expression) {

}

func evalCompareExpression(e *Expression) {}

func searchFunction(name string) *FunctionDefinition {}

func searchDeclaration(name string, currentBlock *Block) *Declaration {

}

func checkArgument(currentBlock *Block, fd *FunctionDefinition, e Expression) {}
