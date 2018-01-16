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

	fix(*Block) Expression

	typeS() *TypeSpecifier
	setType(*TypeSpecifier)
}

// ExpressionImpl provide commonly implementations for Expr.
type ExpressionImpl struct {
	PosImpl // ExprImpl provide Pos() function.

	// 类型
	typeSpecifier *TypeSpecifier
}

// expr provide restraint interface.
func (expr *ExpressionImpl) expr() {}

func (expr *ExpressionImpl) typeS() *TypeSpecifier {
	return expr.typeSpecifier
}

func (expr *ExpressionImpl) setType(t *TypeSpecifier) {
	expr.typeSpecifier = t
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

func (expr *CommaExpression) fix(currentBlock *Block) Expression {

	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)
	expr.typeSpecifier = expr.right.typeS()
	return expr
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

func (expr *AssignExpression) fix(currentBlock *Block) Expression {
	_, ok := expr.left.(*IdentifierExpression)
	if !ok {
		compileError(expr.left.Position(), 0, "")
	}

	expr.left = expr.left.fix(currentBlock)
	operand := expr.operand.fix(currentBlock)
	expr.operand = createAssignCast(expr.operand, expr.left.typeS())
	expr.typeSpecifier = expr.left.typeS()
	return expr

}

// ==============================
// BinaryExpression
// ==============================

// BinaryExpression 二元表达式
type BinaryExpression struct {
	ExpressionImpl

	// 操作符
	operator BinaryOperatorKind
	left     Expression
	right    Expression
}

func (expr *BinaryExpression) fix(currentBlock *Block) Expression {
	switch expr.operator {
		// 数学计算
	case AddOperator, SubOperator, MulOperator, DivOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		expr = evalMathExpression(currentBlock, expr)

		switch expr.(type) {
		case NumberExpression, StringExpression:
			return expr
		}

		expr = castBinaryExpression(expr)

		if isNumber(expr.left.typeS()) && isNumber(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: NumberType}
		} else if isString(expr.left.typeS()) && isString(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: StringType}
		} else {
			compileError(expr.Position(), 0, "")
		}

		return expr

		// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		expr = evalCompareExpression(expr)

		switch expr.(type) {
		case BooleanExpression:
			return expr
		}

		expr = castBinaryExpression(expr)

		if (expr.left.typeS().basicType != expr.right.typeS().basicType) || expr.left.typeS().derive != nil || expr.right.typeS.derive != nil {
			compileError(expr.Position(), 0, "")
		}
		expr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

		return expr

	case LogicalAndOperator, LogicalOrOperator:
		expr = expr.left.fix(currentBlock)
		expr = expr.right.fix(currentBlock)

		if isBoolean(expr.left.typeS()) && isBoolean(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
		} else {
			compileError(expr.Position(), 0, "")
		}
		return expr
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

func (expr *MinusExpression) fix(currentBlock *Block) Expression {
	newExpr := expr.operand.fix(currentBlock)

	if !isNumber(newExpr.typeS()) {
		compileError(expr.Position(), 0, "")
	}
	newExpr.setType(expr.operand.typeS())

	kind, ok := newExpr.(NumberType)
	if ok {
		kind.numberValue = -kind.numberValue
	}
	return newExpr

}

// ==============================
// LogicalNotExpression
// ==============================

// LogicalNotExpression 逻辑非表达式
type LogicalNotExpression struct {
	ExpressionImpl
	operand Expression
}

func (expr *LogicalNotExpression) fix(currentBlock *Block) Expression {
	expr.operand.fix(currentBlock)

	b, ok := expr.operand.(*BooleanExpression)
	if !ok {
		return
	}

	// TODO 增加返回值，返回*BooleanExpression
	//b.booleanValue = !b.booleanValue

	//expr.typeSpecifier = expr.operand.typeS()

}

// ==============================
// FunctionCallExpression
// ==============================

// FunctionCallExpression 函数调用表达式
type FunctionCallExpression struct {
	ExpressionImpl

	// TODO 除去函数名
	// 函数名
	name string
	// 实参列表
	argument []Expression
	// TODO
	function Expression
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	expr.function.fix(currentBlock)

	fd := searchFunction(expr.name)
	if fd == nil {
		compileError(expr.Position(), 0, "")

		checkArgument(currentBlock, fd, expr)

		expr.typeSpecifier = &TypeSpecifier{basicType: fd.typeS().basicType}
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

	booleanValue  bool
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
}

// ==============================
// NumberExpression
// ==============================

// NumberExpression 数字表达式
type NumberExpression struct {
	ExpressionImpl

	numberValue   float64
}

func (expr *NumberExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: NumberType}

}

// ==============================
// StringExpression
// ==============================

// StringExpression 字符串表达式
type StringExpression struct {
	ExpressionImpl
	stringValue   string
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: StringType}
}

// ==============================
// IdentifierExpression
// ==============================

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionImpl
	name string
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {

	decl := searchDeclaration(expr.name, currentBlock)
	if decl != nil {
		expr.typeSpecifier = decl.typeSpecifier
		//expr.u.identifier.is_function = DVM_FALSE
		//expr.u.identifier.u.declaration = decl
		return
	}
	fd := searchFunction(expr.name)
	if fd == nil {
		compileError(expr.Position(), 0, "")
	}
	expr.typeSpecifier = fd.typeS()

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
