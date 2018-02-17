package compiler

import (
	"../vm"
)

//
// BinaryOperatorKind ...
//
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
type Expression interface {
	// Pos接口
	Pos

	// 用于类型修正,以及简单的类型转换
	fix(*Block) Expression
	// 生成字节码
	generate(*vm.Executable, *Block, *OpcodeBuf)

	typeS() *TypeSpecifier
	setType(*TypeSpecifier)

	show(ident int)
}

//
// Expression impl
//
type ExpressionImpl struct {
	// ExprImpl provide Pos() function
	PosImpl

	// 类型
	typeSpecifier *TypeSpecifier
}

func (expr *ExpressionImpl) fix(currentBlock *Block) Expression { return nil }

func (expr *ExpressionImpl) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {}

func (expr *ExpressionImpl) show(ident int) {}

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

func (expr *CommaExpression) show(ident int) {
	printWithIdent("CommaExpr", ident)

	subIdent := ident + 2
	expr.left.show(subIdent)
	expr.right.show(subIdent)
}

func (expr *CommaExpression) fix(currentBlock *Block) Expression {

	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	expr.setType(expr.right.typeS())

	return expr
}

func (expr *CommaExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

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
	// 操作数
	operand Expression
}

func (expr *AssignExpression) show(ident int) {
	printWithIdent("AssignExpr", ident)

	subIdent := ident + 2
	expr.left.show(subIdent)
	expr.operand.show(subIdent)
}

func (expr *AssignExpression) fix(currentBlock *Block) Expression {
	switch expr.left.(type) {
	case *IdentifierExpression, *IndexExpression:
		// pass
	default:
		compileError(expr.left.Position(), NOT_LVALUE_ERR, "")
	}

	expr.left = expr.left.fix(currentBlock)
	expr.operand.fix(currentBlock)
	expr.operand = createAssignCast(expr.operand, expr.left.typeS())

	expr.setType(expr.left.typeS())

	return expr
}

func (expr *AssignExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	// 与generateEx的区别
	ob.generateCode(expr.Position(), vm.VM_DUPLICATE)

	identifierExpr := expr.left.(*IdentifierExpression)

	decl := identifierExpr.inner.(*Declaration)

	generatePopToIdentifier(decl, expr.Position(), ob)
}

// TODO
// 顶层
func (expr *AssignExpression) generateEx(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	identifierExpr := expr.left.(*IdentifierExpression)

	decl := identifierExpr.inner.(*Declaration)

	generatePopToIdentifier(decl, expr.Position(), ob)
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpcodeBuf) {
	var code byte

	offset := getOpcodeTypeOffset(decl.typeSpecifier.basicType)
	if decl.isLocal {
		code = vm.VM_POP_STACK_INT
	} else {
		code = vm.VM_POP_STATIC_INT
	}
	ob.generateCode(pos, code+offset, decl.variableIndex)
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

func (expr *BinaryExpression) show(ident int) {
	printWithIdent("BinaryExpr", ident)

	subIdent := ident + 2
	expr.left.show(subIdent)
	expr.right.show(subIdent)
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
		case *IntExpression, *DoubleExpression, *StringExpression:
			return newExpr
		}

		// 类型转换
		newExpr = castBinaryExpression(expr)
		newBinaryExpr := newExpr.(*BinaryExpression)

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		if isInt(newBinaryExprLeftType) && isInt(newBinaryExprRightType) {
			newExpr.setType(&TypeSpecifier{basicType: vm.IntType})
		} else if isDouble(newBinaryExprLeftType) && isDouble(newBinaryExprRightType) {
			newExpr.setType(&TypeSpecifier{basicType: vm.DoubleType})
		} else if (expr.operator == AddOperator) &&
		((isString(newBinaryExprLeftType) && isString(newBinaryExprRightType)) ||
		(isString(newBinaryExprLeftType) && isNull(newBinaryExpr.left))) {
			newExpr.setType(&TypeSpecifier{basicType: vm.StringType})
		} else {
			compileError(expr.Position(), MATH_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(newBinaryExprLeftType.basicType), int(newBinaryExprRightType.basicType))
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
		newBinaryExpr := newExpr.(*BinaryExpression)

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		if !compareType(newBinaryExprLeftType, newBinaryExprRightType) ||
		(isObject(newBinaryExprLeftType) && isNull(newBinaryExpr.right) ||
		(isObject(newBinaryExprRightType) && isNull(newBinaryExpr.left))) {

			compileError(expr.Position(), COMPARE_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(newBinaryExprLeftType.basicType), int(newBinaryExprRightType.basicType))
		}

		newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})

		return newExpr

		// && ||
	case LogicalAndOperator, LogicalOrOperator:
		expr.left = expr.left.fix(currentBlock)
		expr.right = expr.right.fix(currentBlock)

		if isBoolean(expr.left.typeS()) && isBoolean(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
		} else {
			compileError(expr.Position(), LOGICAL_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(expr.left.typeS().basicType), int(expr.right.typeS().basicType))
		}
		return expr

	default:
		return nil
	}
}

var operatorCodeMap = map[BinaryOperatorKind]byte{
	EqOperator:  vm.VM_EQ_INT,
	NeOperator:  vm.VM_NE_INT,
	GtOperator:  vm.VM_GT_INT,
	GeOperator:  vm.VM_GE_INT,
	LtOperator:  vm.VM_LT_INT,
	LeOperator:  vm.VM_LE_INT,
	AddOperator: vm.VM_ADD_INT,
	SubOperator: vm.VM_SUB_INT,
	MulOperator: vm.VM_MUL_INT,
	DivOperator: vm.VM_DIV_INT,
}

func (expr *BinaryExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	switch operator := expr.operator; operator {
		case GtOperator, GeOperator, LtOperator, LeOperator,
		AddOperator, SubOperator, MulOperator, DivOperator,
		EqOperator, NeOperator:

		if expr.left.typeS().basicType != expr.right.typeS().basicType {
			// TODO
			panic("TODO")
		}

		expr.left.generate(exe, currentBlock, ob)
		expr.right.generate(exe, currentBlock, ob)

		code, ok := operatorCodeMap[operator]
		if !ok {
			// TODO
			panic("TODO")
		}
		// 入栈
		codeOffset := code + getOpcodeTypeOffset(expr.left.typeS().basicType)
		ob.generateCode(expr.Position(), codeOffset)

	case LogicalAndOperator, LogicalOrOperator:
		var jumpCode, logicalCode byte

		if operator == LogicalAndOperator {
			jumpCode = vm.VM_JUMP_IF_FALSE
			logicalCode = vm.VM_LOGICAL_AND
		} else {
			jumpCode = vm.VM_JUMP_IF_TRUE
			logicalCode = vm.VM_LOGICAL_OR
		}

		label := ob.getLabel()

		expr.left.generate(exe, currentBlock, ob)
		ob.generateCode(expr.Position(), vm.VM_DUPLICATE)
		ob.generateCode(expr.Position(), jumpCode, label)

		expr.right.generate(exe, currentBlock, ob)

		// 判断结果
		ob.generateCode(expr.Position(), logicalCode)

		ob.setLabel(label)
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

func (expr *MinusExpression) show(ident int) {
	printWithIdent("MinusExpr", ident)

	subIdent := ident + 2
	expr.operand.show(subIdent)
}

func (expr *MinusExpression) fix(currentBlock *Block) Expression {
	expr.operand = expr.operand.fix(currentBlock)

	if !isInt(expr.operand.typeS()) && !isDouble(expr.operand.typeS()) {
		compileError(expr.Position(), MINUS_TYPE_MISMATCH_ERR, "")
	}

	expr.setType(expr.operand.typeS())

	switch newExpr := expr.operand.(type) {
	case *IntExpression:
		newExpr.intValue = -newExpr.intValue
		return newExpr
	case *DoubleExpression:
		newExpr.doubleValue = -newExpr.doubleValue
		return newExpr
	}

	return expr
}

func (expr *MinusExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)
	code := vm.VM_MINUS_INT + getOpcodeTypeOffset(expr.typeS().basicType)
	ob.generateCode(expr.Position(), code)
}

// ==============================
// LogicalNotExpression
// ==============================

// LogicalNotExpression 逻辑非表达式
type LogicalNotExpression struct {
	ExpressionImpl

	operand Expression
}

func (expr *LogicalNotExpression) show(ident int) {
	printWithIdent("LogicalNotExpr", ident)

	subIdent := ident + 2

	expr.operand.show(subIdent)
}

func (expr *LogicalNotExpression) fix(currentBlock *Block) Expression {
	expr.operand = expr.operand.fix(currentBlock)

	switch newExpr := expr.operand.(type) {
	case *BooleanExpression:
		newExpr.booleanValue = !newExpr.booleanValue
		return newExpr
	}

	if !isBoolean(expr.operand.typeS()) {
		compileError(expr.Position(), LOGICAL_NOT_TYPE_MISMATCH_ERR, "")
	}

	expr.setType(expr.operand.typeS())

	return expr
}

func (expr *LogicalNotExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_LOGICAL_NOT)
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

func (expr *FunctionCallExpression) show(ident int) {
	printWithIdent("FuncCallExpr", ident)

	subIdent := ident + 2

	expr.function.show(subIdent)
	for _, arg := range expr.argumentList {
		printWithIdent("ArgList", subIdent)
		arg.show(subIdent + 2)
	}
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	funcExpr := expr.function.fix(currentBlock)

	identifierExpr, ok := funcExpr.(*IdentifierExpression)
	if !ok {
		compileError(expr.Position(), FUNCTION_NOT_IDENTIFIER_ERR, "")
	}

	fd := SearchFunction(identifierExpr.name)
	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, "Function name: %s\n", identifierExpr.name)
	}

	fd.checkArgument(currentBlock, expr)

	expr.setType(&TypeSpecifier{basicType: fd.typeS().basicType})

	expr.typeSpecifier.deriveList = fd.typeS().deriveList

	return expr
}
func (expr *FunctionCallExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	for _, arg := range expr.argumentList {
		arg.generate(exe, currentBlock, ob)
	}

	expr.function.generate(exe, currentBlock, ob)

	ob.generateCode(expr.Position(), vm.VM_INVOKE)
}

// ==============================
// BooleanExpression
// ==============================

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	ExpressionImpl

	booleanValue bool
}

func (expr *BooleanExpression) show(ident int) {
	printWithIdent("BoolExpr", ident)
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.BooleanType})
	return expr
}

func (expr *BooleanExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if expr.booleanValue {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, 1)
	} else {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, 0)
	}

}

// ==============================
// IntExpression
// ==============================

// IntExpression 数字表达式
type IntExpression struct {
	ExpressionImpl

	intValue int
}

func (expr *IntExpression) show(ident int) {
	printWithIdent("IntExpr", ident)
}

func (expr *IntExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.IntType})
	return expr
}
func (expr *IntExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	if expr.intValue >= 0 && expr.intValue < 256 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, expr.intValue)
	} else if expr.intValue >= 0 && expr.intValue < 65536 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_2BYTE, expr.intValue)
	} else {
		c := vm.NewConstantInt(expr.intValue)
		cpIdx := addConstantPool(exe, c)

		ob.generateCode(expr.Position(), vm.VM_PUSH_INT, cpIdx)
	}
}

// ==============================
// DoubleExpression
// ==============================

// DoubleExpression 数字表达式
type DoubleExpression struct {
	ExpressionImpl

	doubleValue float64
}

func (expr *DoubleExpression) show(ident int) {
	printWithIdent("DoubleExpr", ident)
}

func (expr *DoubleExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.DoubleType})
	return expr
}
func (expr *DoubleExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	if expr.doubleValue == 0.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_DOUBLE_0)

	} else if expr.doubleValue == 1.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_DOUBLE_1)

	} else {
		c := vm.NewConstantDouble(expr.doubleValue)
		cpIdx := addConstantPool(exe, c)

		ob.generateCode(expr.Position(), vm.VM_PUSH_DOUBLE, cpIdx)
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

func (expr *StringExpression) show(ident int) {
	printWithIdent("StringExpr", ident)
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.StringType})
	return expr
}

func (expr *StringExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	c := vm.NewConstantString(expr.stringValue)
	cpIdx := addConstantPool(exe, c)

	ob.generateCode(expr.Position(), vm.VM_PUSH_STRING, cpIdx)
}

// ==============================
// IdentifierExpression
// ==============================
type IdentifierInner interface{}

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionImpl

	name string

	// 声明要么是变量，要么是函数 (FunctionDefinition Declaration)
	inner IdentifierInner
}

func (expr *IdentifierExpression) show(ident int) {
	printWithIdent("IdentifierExpr", ident)
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {
	// 判断是否是变量
	decl := searchDeclaration(expr.name, currentBlock)
	if decl != nil {
		expr.setType(decl.typeSpecifier)
		expr.inner = decl
		return expr
	}

	// 判断是否是函数
	fd := SearchFunction(expr.name)
	if fd != nil {
		expr.setType(fd.typeSpecifier)
		expr.inner = fd
		return expr
	}

	// 都不是,报错
	compileError(expr.Position(), IDENTIFIER_NOT_FOUND_ERR, expr.name)
	return nil
}

func (expr *IdentifierExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	switch inner := expr.inner.(type) {
		// 函数
	case *FunctionDefinition:
		ob.generateCode(expr.Position(), vm.VM_PUSH_FUNCTION, inner.index)
		// 变量
	case *Declaration:
		var code byte

		offset := getOpcodeTypeOffset(inner.typeSpecifier.basicType)
		if inner.isLocal {
			code = vm.VM_PUSH_STACK_INT
		} else {
			code = vm.VM_PUSH_STATIC_INT
		}
		ob.generateCode(expr.Position(), code+offset, inner.variableIndex)
	}
}

// ==============================
// CastExpression
// ==============================

type CastType int

const (
	IntToStringCast CastType = iota
	BooleanToStringCast
	DoubleToStringCast
	IntToDoubleCast
	DoubleToIntCast
)

type CastExpression struct {
	ExpressionImpl

	castType CastType

	operand Expression
}

func (expr *CastExpression) show(ident int) {
	printWithIdent("CastExpr", ident)
}

func (expr *CastExpression) fix(currentBlock *Block) Expression { return expr }

func (expr *CastExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	switch expr.castType {
	case IntToDoubleCast:
		ob.generateCode(expr.Position(), vm.VM_CAST_INT_TO_DOUBLE)
	case DoubleToIntCast:
		ob.generateCode(expr.Position(), vm.VM_CAST_DOUBLE_TO_INT)
	case BooleanToStringCast:
		ob.generateCode(expr.Position(), vm.VM_CAST_BOOLEAN_TO_STRING)
	case IntToStringCast:
		ob.generateCode(expr.Position(), vm.VM_CAST_INT_TO_STRING)
	case DoubleToStringCast:
		ob.generateCode(expr.Position(), vm.VM_CAST_DOUBLE_TO_STRING)
	default:
		panic("TODO")
	}
}

// ==============================
// IndexExpression
// ==============================
type IndexExpression struct {
	ExpressionImpl

	array Expression
	index Expression
}

func (expr *IndexExpression) fix(currentBlock *Block) Expression {

	expr.array = expr.array.fix(currentBlock)
	expr.index = expr.index.fix(currentBlock)

	if expr.array.typeS().deriveList == nil {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	_, ok := expr.array.typeS().deriveList[0].(*ArrayDerive)
	if !ok {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	if !isInt(expr.index.typeS()) {
		compileError(expr.Position(), INDEX_NOT_INT_ERR)
	}

	expr.setType(&TypeSpecifier{basicType: expr.array.typeS().basicType})

	expr.typeS().deriveList = expr.array.typeS().deriveList[1:]

	return expr
}

// ==============================
// MemberExpression
// ==============================
type MemberExpression struct {
	ExpressionImpl

	expression Expression
	memberName string
}

func (expr *MemberExpression) fix(currentBlock *Block) Expression { return expr }

// ==============================
// ArrayLiteralExpression
// ==============================
type ArrayLiteralExpression struct {
	ExpressionImpl

	arrayLiteral []Expression
}

func (expr *ArrayLiteralExpression) fix(currentBlock *Block) Expression {
	if expr.arrayLiteral == nil {
		compileError(expr.Position(), ARRAY_LITERAL_EMPTY_ERR)
	}

	firstElem := expr.arrayLiteral[0]
	firstElem = firstElem.fix(currentBlock)

	elemType := firstElem.typeS()

	for _, elem := range expr.arrayLiteral[1:] {
		elem = elem.fix(currentBlock)
		elem = createAssignCast(elem, elemType)
	}

	expr.setType(&TypeSpecifier{basicType: elemType.basicType})
	expr.typeS().deriveList = []TypeDerive{&ArrayDerive{}}
	expr.typeS().deriveList = append(expr.typeS().deriveList, elemType.deriveList...)

	return expr
}

// ==============================
// ArrayCreation
// ==============================
type ArrayCreation struct {
	ExpressionImpl

	dimensionList []*ArrayDimension
}

func (expr *ArrayCreation) fix(currentBlock *Block) Expression {

	deriveList := []TypeDerive{}

	for _, dim := range expr.dimensionList {
		if dim.expression != nil {
			dim.expression = dim.expression.fix(currentBlock)

			if !isInt(dim.expression.typeS()) {
				compileError(expr.Position(), ARRAY_SIZE_NOT_INT_ERR)
			}
		}

		deriveList = append(deriveList, &ArrayDerive{})
	}
	expr.setType(&TypeSpecifier{basicType: expr.typeSpecifier.basicType})
	expr.typeS().deriveList = deriveList

	return expr
}

// ==============================
// ArrayDimension
// ==============================
type ArrayDimension struct {
	expression Expression
}

// ==============================
// NullDimension
// ==============================
type NullExpression struct {
	ExpressionImpl
}

func (expr *NullExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.NullType})
	return expr
}

// ==============================
// utils
// ==============================

func addConstantPool(exe *vm.Executable, cp vm.Constant) int {
	exe.ConstantPool.Append(cp)

	return exe.ConstantPool.Length() - 1
}
