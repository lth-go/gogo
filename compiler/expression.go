package compiler

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

//
// Expression interface
//

// Expression 表达式接口
type Expression interface {
	// Pos接口
	Pos
	expr()

	fix(*Block) Expression
	generate(*Executable, *Block, *OpcodeBuf)

	typeS() *TypeSpecifier
	setType(*TypeSpecifier)
}

//
// Expression impl
//

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

func (expr *CommaExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	expr.left.generate(exe, currentBlock, ob)
	expr.right.generate(exe, currentBlock, ob)
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

func (expr *AssignExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	generateCode(ob, expr.Position(), DUPLICATE)

	identifierExpr, ok := expr.left.(*IdentifierExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}
	decl := identifierExpr.declaration

	generatePopToIdentifier(decl, expr.Position(), ob)
}

// TODO
func (expr *AssignExpression) generateEx(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	identifierExpr, ok := expr.left.(*IdentifierExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}
	decl := identifierExpr.declaration

	generatePopToIdentifier(decl, expr.Position(), ob)
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpcodeBuf) {
	var code Opcode

	if decl.isLocal {
		code = POP_STACK_INT + get_opcode_type_offset(decl.typeSpecifier.basicType)
	} else {
		code = POP_STATIC_INT + get_opcode_type_offset(decl.typeSpecifier.basicType)
	}
	generateCode(ob, pos, code, decl.variableIndex)
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
	var newExpr Expression

	switch expr.operator {
	// 数学计算
	case AddOperator, SubOperator, MulOperator, DivOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		// 能否合并计算
		newExpr = evalMathExpression(currentBlock, expr)

		switch newExpr.(type) {
		case *NumberExpression, *StringExpression:
			return newExpr
		}

		// 类型转换
		newExpr = castBinaryExpression(expr)

		newBinaryExpr, ok := newExpr.(*BinaryExpression)
		if !ok {
			compileError(expr.Position(), 0, "")
		}

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		if isNumber(newBinaryExprLeftType) && isNumber(newBinaryExprRightType) {
			newExpr.setType(&TypeSpecifier{basicType: NumberType})
		} else if isString(newBinaryExprLeftType) && isString(newBinaryExprRightType) {
			newExpr.setType(&TypeSpecifier{basicType: StringType})
		} else {
			compileError(expr.Position(), 0, "")
		}

		return newExpr

		// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		newExpr = evalCompareExpression(expr)

		switch newExpr.(type) {
		case *BooleanExpression:
			return newExpr
		}

		newExpr = castBinaryExpression(expr)

		newBinaryExpr, ok := newExpr.(*BinaryExpression)
		if !ok {
			compileError(expr.Position(), 0, "")
		}

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		if (newBinaryExprLeftType.basicType != newBinaryExprRightType.basicType) || newBinaryExprLeftType.deriveList != nil || newBinaryExprRightType.deriveList != nil {
			compileError(expr.Position(), 0, "")
		}

		newExpr.setType(&TypeSpecifier{basicType: BooleanType})

		return newExpr

	case LogicalAndOperator, LogicalOrOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		if isBoolean(expr.left.typeS()) && isBoolean(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
		} else {
			compileError(expr.Position(), 0, "")
		}
		return expr
	default:
		compileError(expr.Position(), 0, "")
	}
	return nil
}

var operatorCodeMap = map[BinaryOperatorKind]Opcode{
	EqOperator:  EQ_INT,
	NeOperator:  NE_INT,
	GtOperator:  GT_INT,
	GeOperator:  GE_INT,
	LtOperator:  LT_INT,
	LeOperator:  LE_INT,
	AddOperator: ADD_INT,
	SubOperator: SUB_INT,
	MulOperator: MUL_INT,
	DivOperator: DIV_INT,
}

func (expr *BinaryExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	switch expr.operator {
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator, AddOperator, SubOperator, MulOperator, DivOperator:

		if expr.left.typeS().basicType != expr.right.typeS().basicType {
			compileError(expr.Position(), 0, "")
		}

		code, ok := operatorCodeMap[expr.operator]
		if !ok {
			compileError(expr.Position(), 0, "")
		}

		expr.left.generate(exe, currentBlock, ob)
		expr.right.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), code+get_opcode_type_offset(expr.left.typeS().basicType))

	case LogicalAndOperator:

		falseLabel := getLabel(ob)

		expr.left.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), DUPLICATE)
		generateCode(ob, expr.Position(), JUMP_IF_FALSE, falseLabel)

		expr.right.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), LOGICAL_AND)
		setLabel(ob, falseLabel)

	case LogicalOrOperator:

		trueLabel := getLabel(ob)

		expr.left.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), DUPLICATE)
		generateCode(ob, expr.Position(), JUMP_IF_TRUE, trueLabel)

		expr.right.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), LOGICAL_OR)
		setLabel(ob, trueLabel)
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
	var newExpr Expression
	newExpr = expr.operand.fix(currentBlock)

	if !isNumber(newExpr.typeS()) {
		compileError(expr.Position(), 0, "")
	}

	kind, ok := newExpr.(*NumberExpression)
	if ok {
		kind.numberValue = -kind.numberValue
		return newExpr
	}

	expr.setType(newExpr.typeS())

	return expr
}
func (expr *MinusExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)
	generateCode(ob, expr.Position(), MINUS_INT+get_opcode_type_offset(expr.typeS().basicType))
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

	if !isBoolean(newExpr.typeS()) {
		compileError(expr.Position(), 0, "")
	}

	expr.setType(newExpr.typeS())

	return expr
}
func (expr *LogicalNotExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)
	generateCode(ob, expr.Position(), LOGICAL_NOT)

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
	argumentList []Expression
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	funcExpr := expr.function.fix(currentBlock)
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
	expr.typeSpecifier.deriveList = fd.typeS().deriveList
	return expr
}
func (expr *FunctionCallExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	for _, arg := range expr.argumentList {
		arg.generate(exe, currentBlock, ob)
	}

	expr.function.generate(exe, currentBlock, ob)

	generateCode(ob, expr.Position(), INVOKE)
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
func (expr *BooleanExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	if expr.booleanValue {
		generateCode(ob, expr.Position(), PUSH_INT_1BYTE, 1)
	} else {
		generateCode(ob, expr.Position(), PUSH_INT_1BYTE, 0)
	}

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
func (expr *NumberExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	if expr.numberValue == 0.0 {
		generateCode(ob, expr.Position(), PUSH_DOUBLE_0)

	} else if expr.numberValue == 1.0 {
		generateCode(ob, expr.Position(), PUSH_DOUBLE_1)

	} else {
		cp := &ConstantNumber{numberValue: expr.numberValue}
		cpIdx := addConstantPool(exe, cp)

		generateCode(ob, expr.Position(), PUSH_DOUBLE, cpIdx)
	}
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
func (expr *StringExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	cp := &ConstantString{stringValue: expr.stringValue}

	cpIdx := addConstantPool(exe, cp)
	generateCode(ob, expr.Position(), PUSH_STRING, cpIdx)
}

// ==============================
// IdentifierExpression
// ==============================

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionImpl

	name string

	// 是否可以去除
	isFunction bool

	// 声明要么是变量，要么是函数
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
func (expr *IdentifierExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	// 函数
	if expr.function != nil {
		generateCode(ob, expr.Position(), PUSH_FUNCTION, expr.function.index)
		return
	}

	// 变量
	var code Opcode

	if expr.declaration.isLocal {
		code = PUSH_STACK_INT + get_opcode_type_offset(expr.declaration.typeSpecifier.basicType)
	} else {
		code = PUSH_STATIC_INT + get_opcode_type_offset(expr.declaration.typeSpecifier.basicType)
	}

	generateCode(ob, expr.Position(), code, expr.declaration.variableIndex)
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

func (expr *CastExpression) fix(currentBlock *Block) Expression {
	return nil
}

func (expr *CastExpression) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	switch expr.castType {

	case BooleanToStringCast:
		generateCode(ob, expr.Position(), CAST_BOOLEAN_TO_STRING)

	case NumberToStringCast:
		generateCode(ob, expr.Position(), CAST_DOUBLE_TO_STRING)

	default:
		compileError(expr.Position(), 0, "")
	}
}

func allocCastExpression(castType CastType, expr Expression) Expression {

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case BooleanToStringCast, NumberToStringCast:
		castExpr.typeSpecifier = &TypeSpecifier{basicType: StringType}
	}

	return castExpr
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
	if src.typeS().deriveList != nil || dest.deriveList != nil {
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

	} else if isString(binaryExpr.left.typeS()) && isNumber(binaryExpr.right.typeS()) {
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
	leftNumberExpr, ok := binaryExpr.left.(*NumberExpression)
	if ok {
		rightNumberExpr, ok := binaryExpr.right.(*NumberExpression)
		if ok {
			newExpr := evalMathExpressionNumber(binaryExpr, leftNumberExpr.numberValue, rightNumberExpr.numberValue)
			return newExpr
		}
	}

	// 字符串链接
	leftStringExpr, ok := binaryExpr.left.(*StringExpression)
	if ok {
		newExpr := chainString(expr)
	}

	return expr
}

func evalMathExpressionNumber(expr Expression, left, right float64) Expression {
	var value float64

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case AddOperator:
		value = left + right
	case SubOperator:
		value = left - right
	case MulOperator:
		value = left * right
	case DivOperator:
		value = left / right
	default:
		compileError(binaryExpr.Position(), 0, "")
	}
	newExpr := &NumberExpression{numberValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: NumberType}
	return newExpr
}

func chainString(expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	rightStr := expressionToString(binaryExpr.right)
	if rightStr != "" {
		return expr
	}

	leftStringExpr, ok := binaryExpr.left.(*StringExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	newStr := leftStringExpr.stringValue + rightStr

	newExpr := &StringExpression{stringValue: newStr}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: StringType}

	return newExpr
}

func expressionToString(expr Expression) string {
	var newStr string

	switch e := expr.(type) {
	case *BooleanExpression:
		if e.booleanValue == true {
			newStr = "true"
		} else {
			newStr = "false"
		}
	case *NumberExpression:
		newStr = strconv.FormatFloat(e.numberValue, 'f', -1, 64)
	case *StringExpression:
		newStr = e.stringValue
	default:
		newStr = ""
	}

	return newStr
}

func evalCompareExpression(expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	switch leftExpr := binaryExpr.left.(type) {
	case *BooleanExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *BooleanExpression:
			newExpr := evalCompareExpressionBoolean(binaryExpr, leftExpr.booleanValue, rightExpr.booleanValue)
			return newExpr
		}
	case *NumberExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *NumberExpression:
			newExpr := evalCompareExpressionNumber(binaryExpr, leftExpr.numberValue, rightExpr.numberValue)
			return newExpr
		}
	case *StringExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *StringExpression:
			newExpr := evalCompareExpressionString(binaryExpr, leftExpr.stringValue, rightExpr.stringValue)
			return newExpr
		}
	}
	return expr
}

func evalCompareExpressionBoolean(expr Expression, left, right bool) Expression {
	var value bool

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
	return newExpr
}

func evalCompareExpressionNumber(expr Expression, left, right float64) Expression {
	var value bool

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	case GtOperator:
		value = (left > right)
	case GeOperator:
		value = (left >= right)
	case LtOperator:
		value = (left < right)
	case LeOperator:
		value = (left <= right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}
	return newExpr
}

func evalCompareExpressionString(expr Expression, left, right string) Expression {
	var value bool
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	case GtOperator:
		value = (left > right)
	case GeOperator:
		value = (left >= right)
	case LtOperator:
		value = (left < right)
	case LeOperator:
		value = (left <= right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: BooleanType}

	return newExpr
}

func checkArgument(currentBlock *Block, fd *FunctionDefinition, expr Expression) {
	functionCallExpr, ok := expr.(*FunctionCallExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	ParameterList := fd.parameterList
	argumentList := functionCallExpr.argumentList

	length := len(ParameterList)
	if len(argumentList) != length {
		compileError(expr.Position(), 0, "")
	}

	for i := 0; i < length; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)
		createAssignCast(argumentList[i], ParameterList[i].typeSpecifier)
	}

}

func get_opcode_type_offset(basicType BasicType) Opcode {
	switch basicType {
	case BooleanType:
		return Opcode(0)
	case NumberType:
		return Opcode(0)
	case StringType:
		return Opcode(1)
	default:
		panic("basic type")
	}
}

func getLabel(ob *OpcodeBuf) int {

	ret := len(ob.labelTableList)

	return ret
}

func setLabel(ob *OpcodeBuf, label int) {
	ob.labelTableList[label].labelAddress = len(ob.labelTableList)
}

func addConstantPool(exe *Executable, cp Constant) int {
	ret := len(exe.constantPool)

	exe.constantPool = append(exe.constantPool, cp)

	return ret
}
