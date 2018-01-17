package parser

import (
	"strconv"
)

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

		newExpr := evalMathExpression(currentBlock, expr)

		switch newExpr.(type) {
		case *NumberExpression, *StringExpression:
			return newExpr
		}

		newExpr = castBinaryExpression(expr)

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
	newExpr := expr.operand.fix(currentBlock)

	boolExpr, ok := newExpr.(*BooleanExpression)
	if ok {
		boolExpr.booleanValue = !boolExpr.booleanValue
		return boolExpr
	}

	if !isBoolea(newExpr.typeS()) {
		compileError(expr.Position(), 0, "")
	}

	//newExpr.setType(expr.operand.typeS())

	return newExpr
}

// ==============================
// FunctionCallExpression
// ==============================

// FunctionCallExpression 函数调用表达式
type FunctionCallExpression struct {
	ExpressionImpl

	// 函数名
	function Expression
	// 实参列表
	// TODO 改成argumentList
	argumentList []Expression
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	funcExpr = expr.function.fix(currentBlock)
	identifierExpr, ok := funcExpr.(*IdentifierExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	fd := searchFunction(identifierExpr.name)
	if fd == nil {
		compileError(expr.Position(), 0, "")
	}

	checkArgument(currentBlock, fd, expr)

	expr.typeSpecifier = &TypeSpecifier{basicType: fd.typeS().basicType}
	// TODO
	expr.typeSpecifier.derive = fd.typeS().derive
	return expr
}

// ==============================
// BooleanExpression
// ==============================

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	ExpressionImpl

	booleanValue bool
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
	return expr
}

// ==============================
// NumberExpression
// ==============================

// NumberExpression 数字表达式
type NumberExpression struct {
	ExpressionImpl

	numberValue float64
}

func (expr *NumberExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: NumberType}

	return expr
}

// ==============================
// StringExpression
// ==============================

// StringExpression 字符串表达式
type StringExpression struct {
	ExpressionImpl
	stringValue string
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.typeSpecifier = &TypeSpecifier{basicType: StringType}

	return expr
}

// ==============================
// IdentifierExpression
// ==============================

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionImpl

	name string

	isFunction bool

	function    *FunctionDefinition
	declaration *Declaration
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {

	// 判断是否是变量
	decl := searchDeclaration(expr.name, currentBlock)
	if decl != nil {
		expr.typeSpecifier = decl.typeSpecifier
		expr.isFunction = false
		expr.declaration = decl
		return expr
	}

	// 判断是否是函数
	fd := searchFunction(expr.name)
	if fd != nil {
		expr.typeSpecifier = fd.typeS()
		expr.isFunction = true
		expr.function = fd
	}

	// 都不是,报错
	compileError(expr.Position(), 0, "")
	return nil

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

// 声明类型转换, 目前仅有number类型,不存在类型转换
// TODO: 待去除
func createAssignCast(src Expression, dest *TypeSpecifier) Expression {
	if src.typeS().derive != nil || dest.derive != nil {
		compileError(src.Position(), 0, "")
	}

	if src.typeS().basicType == dest.basicType {
		return src
	}

	compileError(src.Position(), 0, "")
	return nil
}

func castBinaryExpression(expr Expression) Expression {

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	if isString(binaryExpr.left.typeS()) && isBoolean(binaryExpr.right.typeS()) {
		newExpr := allocCastExpression(BooleanToStringCast, binaryExpr.right)
		return newExpr

	} else if isString(binaryExpr.left.typeS()) && isInt(binaryExpr.right.typeS()) {
		newExpr := allocCastExpression(NumberToStringCast, binaryExpr.right)
		return newExpr

	}
	return expr
}

func evalMathExpression(currentBlock *Block, expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	// number 运算
	leftNumberExpr, ok := binaryExpr.left.(NumberType)
	if ok {
		rightNumberExpr, ok := binaryExpr.right.(NumberType)
		if ok {
			newExpr := evalMathExpressionNumber(binaryExpr, leftNumberExpr.numberValue, rightNumberExpr.numberValue)
			return newExpr
		}
	}

	// 字符串链接
	leftNumberExpr, ok := binaryExpr.left.(StringType)
	if ok {
		newExpr := chainString(expr)
	}

	return expr
}

func evalMathExpressionNumber(expr Expression, left, right float64) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case AddOperator:
		value := left + right
	case SubOperator:
		value := left - right
	case MulOperator:
		value := left * right
	case DivOperator:
		value := left / right
	default:
		compileError(binaryExpr.Position(), 0, "")
	}
	newExpr := &NumberExpression{numberValue: value, typeSpecifier: &TypeSpecifier{basicType: NumberType}}
	return newExpr
}

func chainString(expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	rightStr = expressionToString(binaryExpr.right)
	if !rightStr {
		return expr
	}

	leftStringExpr, ok := binaryExpr.(*StringExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	newStr = leftStringExpr.stringValue + rightStr

	newExpr := &StringExpression{stringValue: newStr, typeSpecifier: &TypeSpecifier{basicType: StringType}}

	return newExpr
}

func expressionToString(expr Expression) string {

	switch e := expr.(type) {
	case BooleanExpression:
		if e.booleanValue == true {
			newStr := "true"
		} else {
			newStr := "false"
		}
	case NumberExpression:
		newStr := strconv.FormatFloat(e.numberValue, 'f', 64)
	case StringExpression:
		newStr := e.stringValue
	default:
		newStr := ""
	}

	return newStr
}

func evalCompareExpression(expr Expression) {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	switch leftExpr := binaryExpr.left.(type) {
	case BooleanType:
		switch rightExpr := binaryExpr.right.(type) {
		case BooleanType:
			newExpr := evalCompareExpressionBoolean(binaryExpr, leftExpr.booleanValue, rightExpr.booleanValue)
			return newExpr
		}
	case NumberType:
		switch rightExpr := binaryExpr.right.(type) {
		case NumberType:
			newExpr := evalCompareExpressionNumber(binaryExpr, leftExpr.numberValue, rightExpr.numberValue)
			return newExpr
		}
	case StringType:
		switch rightExpr := binaryExpr.right.(type) {
		case StringType:
			newExpr := evalCompareExpressionString(binaryExpr, leftExpr.stringValue, rightExpr.stringValue)
			return newExpr
		}
	}
	return expr
}

func evalCompareExpressionBoolean(expr Expression, left, right bool) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value := (left == right)
	case NeOperator:
		value := (left != right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value, typeSpecifier: &TypeSpecifier{basicType: BooleanType}}
	return newExpr
}

func evalCompareExpressionNumber(expr Expression, left, right float64) {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value := (left == right)
	case NeOperator:
		value := (left != right)
	case GtOperator:
		value := (left > right)
	case GeOperator:
		value := (left >= right)
	case LeOperator:
		value := (left < right)
	case LeOperator:
		value := (left <= right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value, typeSpecifier: &TypeSpecifier{basicType: BooleanType}}
	return newExpr
}

func evalCompareExpressionString(expr Expression, left, right string) {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value := (left == right)
	case NeOperator:
		value := (left != right)
	case GtOperator:
		value := (left > right)
	case GeOperator:
		value := (left >= right)
	case LeOperator:
		value := (left < right)
	case LeOperator:
		value := (left <= right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value, typeSpecifier: &TypeSpecifier{basicType: BooleanType}}
	return newExpr
}

func checkArgument(currentBlock *Block, fd *FunctionDefinition, expr Expression) {
	functionCallExpr, ok := expr.(*FunctionCallExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	ParameterList := fd.parameterList
	argumentList := functionCallExpr.argumentList

	length = len(ParameterList)
	if len(argumentList) != length {
		compileError(binaryExpr.Position(), 0, "")
	}

	for i := 0; i < length; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)
		createAssignCast(argumentList[i], ParameterList[i].typeSpecifier)
	}

}

// ==============================
// CastExpression
// ==============================

type CastType int

const (
	BooleanToStringCast CastType = iota
	NumberToStringCast
)

type CastExpression struct {
	ExpressionImpl

	castType CastType

	operand Expression
}

func (expr *CastExpression) fix(currentBlock *Block) Expression {}

func allocCastExpression(castType CastType, expr Expression) Expression {

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case BooleanToStringCast, NumberToStringCast:
		castExpr.typeSpecifier = &TypeSpecifier{basicType: StringType}
	}

	return castExpr
}
