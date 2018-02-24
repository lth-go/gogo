package compiler

import (
	"../vm"
)

//
// BinaryOperatorKind
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

//
// Cast
//
type CastType int

const (
	IntToStringCast CastType = iota
	BooleanToStringCast
	DoubleToStringCast
	IntToDoubleCast
	DoubleToIntCast
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

func (expr *ExpressionImpl) typeS() *TypeSpecifier { return expr.typeSpecifier }

func (expr *ExpressionImpl) setType(t *TypeSpecifier) { expr.typeSpecifier = t }

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
	expr.typeS().fix()
	return expr
}

func (expr *BooleanExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if expr.booleanValue {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, 1)
	} else {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, 0)
	}

}

func createBooleanExpression(pos Position) *BooleanExpression {
	expr := &BooleanExpression{}
	expr.SetPosition(pos)

	return expr
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
	expr.typeS().fix()
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

func createIntExpression(pos Position) *IntExpression {
	expr := &IntExpression{}
	expr.SetPosition(pos)

	return expr
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
	expr.typeS().fix()
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

func createDoubleExpression(pos Position) *DoubleExpression {
	expr := &DoubleExpression{}
	expr.SetPosition(pos)

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

func (expr *StringExpression) show(ident int) {
	printWithIdent("StringExpr", ident)
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.StringType})
	expr.typeS().fix()
	return expr
}

func (expr *StringExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	c := vm.NewConstantString(expr.stringValue)
	cpIdx := addConstantPool(exe, c)

	ob.generateCode(expr.Position(), vm.VM_PUSH_STRING, cpIdx)
}

func createStringExpression(pos Position) *StringExpression {
	expr := &StringExpression{}
	expr.SetPosition(pos)

	return expr
}

// ==============================
// NullExpression
// ==============================
type NullExpression struct {
	ExpressionImpl
}

func (expr *NullExpression) show(ident int) {
	printWithIdent("NullExpr", ident)
}

func (expr *NullExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.NullType})
	expr.typeS().fix()
	return expr
}

func (expr *NullExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	ob.generateCode(expr.Position(), vm.VM_PUSH_NULL)
}

func createNullExpression(pos Position) *NullExpression {
	expr := &NullExpression{}
	expr.SetPosition(pos)
	return expr
}

// ==============================
// IdentifierExpression
// ==============================
type IdentifierInner interface{}

type FunctionIdentifier struct {
	functionDefinition *FunctionDefinition
	functionIndex      int
}

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
		expr.typeS().fix()
		return expr
	}

	// 判断是否是函数
	fd := SearchFunction(expr.name)
	if fd != nil {
		compiler := getCurrentCompiler()

		expr.setType(create_function_derive_type(fd))
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			functionIndex:      reserve_function_index(compiler, fd),
		}

		expr.typeS().fix()

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

		offset := getOpcodeTypeOffset(inner.typeSpecifier)
		if inner.isLocal {
			code = vm.VM_PUSH_STACK_INT
		} else {
			code = vm.VM_PUSH_STATIC_INT
		}
		ob.generateCode(expr.Position(), code+offset, inner.variableIndex)
	}
}

func createIdentifierExpression(name string, pos Position) *IdentifierExpression {
	expr := &IdentifierExpression{name: name}
	expr.SetPosition(pos)
	return expr
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
	expr.typeS().fix()

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
	case *IdentifierExpression, *IndexExpression, *MemberExpression:
		// pass
	default:
		compileError(expr.left.Position(), NOT_LVALUE_ERR, "")
	}

	expr.left = expr.left.fix(currentBlock)

	expr.operand = expr.operand.fix(currentBlock)
	expr.operand = createAssignCast(expr.operand, expr.left.typeS())

	expr.setType(expr.left.typeS())
	expr.typeS().fix()

	return expr
}

func (expr *AssignExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	// 与generateEx的区别
	ob.generateCode(expr.Position(), vm.VM_DUPLICATE)

	generate_pop_to_lvalue(exe, currentBlock, expr.left, ob)
}

// TODO
// 顶层
func (expr *AssignExpression) generateEx(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)

	generate_pop_to_lvalue(exe, currentBlock, expr.left, ob)
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

	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	switch expr.operator {
	// 数学计算
	case AddOperator, SubOperator, MulOperator, DivOperator:

		// 能否合并计算
		newExpr := evalMathExpression(currentBlock, expr)
		switch newExpr.(type) {
		case *IntExpression, *DoubleExpression, *StringExpression:
			newExpr.typeS().fix()
			return newExpr
		}

		// 类型转换
		newBinaryExpr := castBinaryExpression(expr)

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		if isInt(newBinaryExprLeftType) && isInt(newBinaryExprRightType) {
			newBinaryExpr.setType(&TypeSpecifier{basicType: vm.IntType})

		} else if isDouble(newBinaryExprLeftType) && isDouble(newBinaryExprRightType) {
			newBinaryExpr.setType(&TypeSpecifier{basicType: vm.DoubleType})

		} else if expr.operator == AddOperator {
			if (isString(newBinaryExprLeftType) && isString(newBinaryExprRightType)) ||
				(isString(newBinaryExprLeftType) && isNull(newBinaryExpr.left)) {
				newBinaryExpr.setType(&TypeSpecifier{basicType: vm.StringType})
			}
		} else {
			compileError(expr.Position(), MATH_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(newBinaryExprLeftType.basicType), int(newBinaryExprRightType.basicType))
		}

		newBinaryExpr.typeS().fix()
		return newBinaryExpr

	// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:

		newExpr := evalCompareExpression(expr)
		switch newExpr.(type) {
		case *BooleanExpression:
			newExpr.typeS().fix()
			return newExpr
		}

		newBinaryExpr := castBinaryExpression(expr)

		newBinaryExprLeftType := newBinaryExpr.left.typeS()
		newBinaryExprRightType := newBinaryExpr.right.typeS()

		// TODO 字符串是否能跟null比较
		if !(compareType(newBinaryExprLeftType, newBinaryExprRightType) ||
			(isObject(newBinaryExprLeftType) && isNull(newBinaryExpr.right) ||
				(isNull(newBinaryExpr.left) && isObject(newBinaryExprRightType)))) {
			compileError(expr.Position(), COMPARE_TYPE_MISMATCH_ERR, getTypeName(newBinaryExprLeftType), getTypeName(newBinaryExprRightType))
		}

		newBinaryExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})

		newBinaryExpr.typeS().fix()
		return newBinaryExpr

	// && ||
	case LogicalAndOperator, LogicalOrOperator:

		if isBoolean(expr.left.typeS()) && isBoolean(expr.right.typeS()) {
			expr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
			expr.typeS().fix(0)
			return expr
		}

		compileError(expr.Position(), LOGICAL_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(expr.left.typeS().basicType), int(expr.right.typeS().basicType))
		return nil

	default:
		panic("TODO")
	}
}

func (expr *BinaryExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	switch operator := expr.operator; operator {
	case GtOperator, GeOperator, LtOperator, LeOperator,
		AddOperator, SubOperator, MulOperator, DivOperator,
		EqOperator, NeOperator:

		var offset byte

		leftExpr := expr.left
		rightExpr := expr.right

		leftExpr.generate(exe, currentBlock, ob)
		rightExpr.generate(exe, currentBlock, ob)

		code, ok := operatorCodeMap[operator]
		if !ok {
			// TODO
			panic("TODO")
		}

		// TODO 啥意思
		if (isNull(leftExpr) && !isNull(rightExpr)) ||
			(!isNull(leftExpr) && isNull(rightExpr)) {
			offset = byte(2)
		} else if (operator == EqOperator || operator == NeOperator) &&
			isString(leftExpr.typeS()) {
			offset = byte(3)
		} else {
			offset = getOpcodeTypeOffset(expr.left.typeS())
		}

		// 入栈
		ob.generateCode(expr.Position(), code+offset)

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

	expr.typeS().fix()
	return expr
}

func (expr *MinusExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.operand.generate(exe, currentBlock, ob)
	code := vm.VM_MINUS_INT + getOpcodeTypeOffset(expr.typeS())
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
		newExpr.setType(createTypeSpecifier(vm.BooleanType, expr.Position()))
		newExpr.typeS().fix()
		return newExpr
	}

	if !isBoolean(expr.operand.typeS()) {
		compileError(expr.Position(), LOGICAL_NOT_TYPE_MISMATCH_ERR, "")
	}

	expr.setType(expr.operand.typeS())
	expr.typeS().fix()

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
	var fd *FunctionDefinition
	var arrayBase *TypeSpecifier

	funcExpr := expr.function.fix(currentBlock)

	expr.function = funcExpr

	switch todo := funcExpr.(type) {
	case *IdentifierExpression:
		fd = SearchFunction(identifierExpr.name)
	case *MemberExpression:
		// TODO 暂时没有列表和字符串方法
		if isArray(todo.expression.typeS()) {
			fd = compiler.arrayMethod[todo.methodIndex]
			arrayBase = todo.expression.typeS()
			arrayBase.deriveList = todo.expression.typeS().deriveList[1:]
		} else if isString(todo.expression.typeS()) {
			fd = compiler.stringMethod[todo.methodIndex]
		} else {
			switch member := todo.declaration.(type) {
			case *FieldMember:
				compileError(expr.Position(), FIELD_CAN_NOT_CALL_ERR, member.name)
			case *MethodMember:
				fd = member.functionDefinition
			default:
				panic("TODO")
			}
		}
	}

	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, identifierExpr.name)
	}

	fd.checkArgument(currentBlock, expr, arrayBase)

	expr.setType(&TypeSpecifier{basicType: fd.typeS().basicType})

	expr.typeSpecifier.deriveList = fd.typeS().deriveList

	if expr.typeS().basicType == vm.ClassType {
		expr.typeS().classRef.identifier = fd.typeS().classRef.identifier
		expr.typeS().fix()
	}

	expr.typeS().fix()
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
// MemberExpression
// ==============================
type MemberExpression struct {
	ExpressionImpl

	expression Expression
	memberName string

	declaration MemberDeclaration
	methodIndex int
}

func (expr *MemberExpression) show(ident int) {
	printWithIdent("MemberExpr", ident)
}

func (expr *MemberExpression) fix(currentBlock *Block) Expression { return expr }

func (expr *MemberExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {}

func createMemberExpression(expression Expression, memberName string, pos Position) *MemberExpression {
	expr := &MemberExpression{
		expression: expression,
		memberName: memberName,
	}
	expr.SetPosition(pos)

	return expr
}

// ==============================
// ThisExpression
// ==============================
type ThisExpression struct {
	ExpressionImpl
}

func (expr *ThisExpression) fix(currentBlock *Block) Expression { return expr }

func (expr *ThisExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {}

func createThisExpression(pos Position) *ThisExpression {
	expr := &ThisExpression{}

	expr.SetPosition(pos)

	return expr
}

// ==============================
// CastExpression
// ==============================

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
// ArrayLiteralExpression
// ==============================
// 创建列表时的值, eg,{1,2,3,4}
type ArrayLiteralExpression struct {
	ExpressionImpl

	arrayLiteral []Expression
}

func (expr *ArrayLiteralExpression) show(ident int) {
	printWithIdent("ArrayLiteralExpr", ident)
}

func (expr *ArrayLiteralExpression) fix(currentBlock *Block) Expression {
	if expr.arrayLiteral == nil || len(expr.arrayLiteral) == 0 {
		compileError(expr.Position(), ARRAY_LITERAL_EMPTY_ERR)
	}

	firstElem := expr.arrayLiteral[0]
	firstElem = firstElem.fix(currentBlock)

	elemType := firstElem.typeS()

	for i := 1; i < len(expr.arrayLiteral); i++ {
		expr.arrayLiteral[i] = expr.arrayLiteral[i].fix(currentBlock)
		expr.arrayLiteral[i] = createAssignCast(expr.arrayLiteral[i], elemType)
	}

	expr.setType(&TypeSpecifier{basicType: elemType.basicType})

	expr.typeS().deriveList = []TypeDerive{&ArrayDerive{}}
	expr.typeS().deriveList = append(expr.typeS().deriveList, elemType.deriveList...)

	return expr
}

func (expr *ArrayLiteralExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if !expr.typeS().isArrayDerive() {
		panic("TODO")
	}

	if expr.arrayLiteral == nil {
		panic("TODO")
	}

	count := len(expr.arrayLiteral)
	if count == 0 {
		panic("TODO")
	}

	for _, subExpr := range expr.arrayLiteral {
		subExpr.generate(exe, currentBlock, ob)
	}

	itemType := expr.arrayLiteral[0].typeS()
	offset := getOpcodeTypeOffset(itemType)

	ob.generateCode(expr.Position(), vm.VM_NEW_ARRAY_LITERAL_INT+offset, count)
}

// ==============================
// ArrayCreation
// ==============================

//
// ArrayDimension
//
// 列表后面的括号
type ArrayDimension struct {
	expression Expression
}

// 列表创建
type ArrayCreation struct {
	ExpressionImpl

	dimensionList []*ArrayDimension
}

func (expr *ArrayCreation) show(ident int) {
	printWithIdent("ArrayCreationExpr", ident)
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

func (expr *ArrayCreation) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	index := AddTypeSpecifier(expr.typeS(), exe)

	if !expr.typeS().isArrayDerive() {
		panic("TODO")
	}

	dimCount := 0
	for _, dim := range expr.dimensionList {
		if dim.expression == nil {
			break
		}

		dim.expression.generate(exe, currentBlock, ob)
		dimCount++
	}

	ob.generateCode(expr.Position(), vm.VM_NEW_ARRAY, dimCount, index)
}

func createBasicArrayCreation(typ *TypeSpecifier, dim_expr_list, dim_list []*ArrayDimension, pos Position) {
	expr := createClassArrayCreation(typ, dim_expr_list, dim_list, pos)

	return expr
}

func createClassArrayCreation(typ *TypeSpecifier, dim_expr_list, dim_list []*ArrayDimension, pos Position) {

	expr := &ArrayCreation{
		dimensionList: dim_expr_list,
	}

	expr.SetPosition(pos)

	expr.setType(typ)

	expr.dimensionList = append(expr.dimensionList, dim_list...)

	return expr
}

// ==============================
// IndexExpression
// ==============================
type IndexExpression struct {
	ExpressionImpl

	array Expression
	index Expression
}

func (expr *IndexExpression) show(ident int) {
	printWithIdent("IndexExpr", ident)

	subIdent := ident + 2
	expr.array.show(subIdent)
	expr.index.show(subIdent)
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

func (expr *IndexExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr.array.generate(exe, currentBlock, ob)
	expr.index.generate(exe, currentBlock, ob)

	code := vm.VM_PUSH_ARRAY_INT + getOpcodeTypeOffset(expr.typeS())
	ob.generateCode(expr.Position(), code)
}

func createIndexExpression(array, index Expression, pos Position) *IndexExpression {
	expr := &IndexExpression{
		array: array,
		index: index,
	}
	expr.SetPosition(pos)

	return expr
}

// ==============================
// NewExpression
// ==============================
type NewExpression struct {
	PosImpl

	className       string
	classDefinition ClassDefinition
	classIndex      int

	methodName        string
	methodDeclaration *MemberDeclaration
	argumentList      []Expression
}

func (expr *NewExpression) fix(currentBlock *Block) Expression                              {}
func (expr *NewExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {}

func createNewExpression(className, memthodName string, argumentList []Expression, pos Position) *NewExpression {
	expr := &NewExpression{
		className:    className,
		methodName:   memthodName,
		argumentList: argumentList,
	}

	expr.SetPosition(pos)

	return expr
}
