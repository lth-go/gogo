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
	_, ok := e.left.(*IdentifierExpression)
	if !ok {
		compileError(e.left.Position(), 0, "")
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
			compileError(e.Position(), 0, "")
		}

	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		e.left.fix(currentBlock)
		e.right.fix(currentBlock)

		evalCompareExpression(e)

		castBinaryExpression(e)

		if e.left.typeS() != e.right.typeS() {
			compileError(e.Position(), 0, "")
		}
		e.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

	case LogicalAndOperator, LogicalOrOperator:
		e.left.fix(currentBlock)
		e.right.fix(currentBlock)

		if isBoolean(e.left.typeS()) && isBoolean(e.right.typeS()) {
			e.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

		} else {

			compileError(e.Position(), 0, "")
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
		compileError(e.Position(), 0, "")
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
		compileError(e.Position(), 0, "")
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
		compileError(e.Position(), 0, "")

		checkArgument(currentBlock, fd, e)

		e.typeSpecifier = &TypeSpecifier{basicType: fd.typeS().basicType}
		// TODO
		//expr.type.derive = fd.type.derive
	}
}

// ==============================
// BooleanExpression
// ==============================

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	ExpressionImpl
	typeSpecifier *TypeSpecifier
	booleanValue  bool
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
		//expr.u.identifier.is_function = DVM_FALSE
		//expr.u.identifier.u.declaration = decl
		return
	}
	fd := searchFunction(e.name)
	if fd == nil {
		compileError(e.Position(), 0, "")
	}
	e.typeSpecifier = fd.typeS()

	//expr.u.identifier.is_function = DVM_TRUE
	//expr.u.identifier.u.function = fd

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

func createAssignCast(src Expression, dest *TypeSpecifier) Expression {
	if src.typeS().derive != nil || dest.derive != nil {
		compileError(src.Position(), 0, "")
	}

	if src.typeS().basicType == dest.basicType {
		return src
	}

	compileError(src.Position(), 0, "")
}


func castBinaryExpression(expr Expression) {

	e, ok := expr.(BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	if isString(e.left.typeS()) && isBoolean(e.right.typeS()) {
		allocCastExpression(BOOLEAN_TO_STRING_CAST, e.right)

	} else if isString(e.left.typeS()) && isInt(e.right.typeS()) {
		allocCastExpression(INT_TO_STRING_CAST, e.right)

	} else if isString(e.left.typeS()) && isDouble(e.right.typeS()) {
		allocCastExpression(DOUBLE_TO_STRING_CAST, e.right)
	}
}

func evalMathExpression(currentBlock *Block, e Expression) {
	// TODO !!!类型有问题
	expr, ok := e.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	if expr.operator == AddOperator {
		expr.u.numberValue = left + right
	} else if expr.operator == SubOperator {
		expr.u.numberValue = left - right
	} else if expr.operator == MulOperator {
		expr.u.numberValue = left * right
	} else if expr.operator == DivOperator {
		expr.u.numberValue = left / right
	} else {
		compileError(expr.Position(), 0, "")
	}

	expr.typeSpecifier = &TypeSpecifier{basicType: NumberType}
}

func evalCompareExpression(e Expression) {
	expr, ok := e.(BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	if expr.left.kind == BOOLEAN_EXPRESSION && e.right.kind == BOOLEAN_EXPRESSION {
		eval_compare_expression_boolean( expr, expr.left.u.boolean_value, expr.right.u.boolean_value)

	} else if expr.left.kind == DOUBLE_EXPRESSION && expr.right.kind == DOUBLE_EXPRESSION {
		eval_compare_expression_double( expr, expr.left.u.double_value, expr.right.u.double_value)

	} else if expr.left.kind == STRING_EXPRESSION && expr.right.kind == STRING_EXPRESSION {
		eval_compare_expression_string( expr, expr.left.u.string_value, expr.right.u.string_value)
	}
	return expr

}
func eval_compare_expression_boolean(expr Expression, left, right bool) {
	if (expr->kind == EQ_EXPRESSION) {
		expr->u.boolean_value = (left == right);
	} else if (expr->kind == NE_EXPRESSION) {
		expr->u.boolean_value = (left != right);
	} else {
		DBG_assert(0, ("expr->kind..%d\n", expr->kind));
	}

	expr->kind = BOOLEAN_EXPRESSION;
	expr->type = dkc_alloc_type_specifier(DVM_BOOLEAN_TYPE);

	return expr;
}
func eval_compare_expression_double(expr Expression , left,  right float64) {
	if (expr->kind == EQ_EXPRESSION) {
		expr->u.boolean_value = (left == right);
	} else if (expr->kind == NE_EXPRESSION) {
		expr->u.boolean_value = (left != right);
	} else if (expr->kind == GT_EXPRESSION) {
		expr->u.boolean_value = (left > right);
	} else if (expr->kind == GE_EXPRESSION) {
		expr->u.boolean_value = (left >= right);
	} else if (expr->kind == LT_EXPRESSION) {
		expr->u.boolean_value = (left < right);
	} else if (expr->kind == LE_EXPRESSION) {
		expr->u.boolean_value = (left <= right);
	} else {
		DBG_assert(0, ("expr->kind..%d\n", expr->kind));
	}

	expr->kind = BOOLEAN_EXPRESSION;
	expr->type = dkc_alloc_type_specifier(DVM_BOOLEAN_TYPE);

	return expr;
}

func eval_compare_expression_string(expr Expression, left,right string ) {
	int cmp;

	cmp = dvm_wcscmp(left, right);

	if (expr->kind == EQ_EXPRESSION) {
		expr->u.boolean_value = (cmp == 0);
	} else if (expr->kind == NE_EXPRESSION) {
		expr->u.boolean_value = (cmp != 0);
	} else if (expr->kind == GT_EXPRESSION) {
		expr->u.boolean_value = (cmp > 0);
	} else if (expr->kind == GE_EXPRESSION) {
		expr->u.boolean_value = (cmp >= 0);
	} else if (expr->kind == LT_EXPRESSION) {
		expr->u.boolean_value = (cmp < 0);
	} else if (expr->kind == LE_EXPRESSION) {
		expr->u.boolean_value = (cmp <= 0);
	} else {
		DBG_assert(0, ("expr->kind..%d\n", expr->kind));
	}

	MEM_free(left);
	MEM_free(right);

	expr->kind = BOOLEAN_EXPRESSION;
	expr->type = dkc_alloc_type_specifier(DVM_BOOLEAN_TYPE);

	return expr;
}


func checkArgument(currentBlock *Block, fd *FunctionDefinition, e Expression) {
	ParameterList *param
	ArgumentList *arg;

	for param = fd.parameter, arg = expr.u.function_call_expression.argument; param && arg; param = param.next, arg = arg.next {
		arg.expression = fix_expression(current_block, arg.expression);
		createAssignCast(arg.expression, param.typeSpecifier)
	}

	if (param || arg) {
		compileError(expr.line_number, ARGUMENT_COUNT_MISMATCH_ERR, MESSAGE_ARGUMENT_END);
	}
}


type CastType int
const (
	BOOLEAN_TO_STRING_CAST  CastType = iota
	NUMBER_TO_STRING_CAST
)
func allocCastExpression(t int, e Expression) {}
