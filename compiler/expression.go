package compiler

import (
	"github.com/lth-go/gogogogo/vm"
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
	generate(*vm.Executable, *Block, *OpCodeBuf)

	typeS() *TypeSpecifier
	setType(*TypeSpecifier)

	show(indent int)
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

func (expr *ExpressionImpl) fix(currentBlock *Block) Expression                              { return nil }
func (expr *ExpressionImpl) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {}
func (expr *ExpressionImpl) show(indent int)                                                 {}
func (expr *ExpressionImpl) typeS() *TypeSpecifier                                           { return expr.typeSpecifier }
func (expr *ExpressionImpl) setType(t *TypeSpecifier)                                        { expr.typeSpecifier = t }

// ==============================
// BooleanExpression
// ==============================

// BooleanExpression 布尔表达式
type BooleanExpression struct {
	ExpressionImpl

	booleanValue bool
}

func (expr *BooleanExpression) show(indent int) {
	printWithIndent("BoolExpr", indent)
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.setType(newTypeSpecifier(vm.BooleanType))
	expr.typeS().fix()
	return expr
}

func (expr *BooleanExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

func (expr *IntExpression) show(indent int) {
	printWithIndent("IntExpr", indent)
}

func (expr *IntExpression) fix(currentBlock *Block) Expression {
	expr.setType(newTypeSpecifier(vm.IntType))
	expr.typeS().fix()
	return expr
}
func (expr *IntExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	if expr.intValue >= 0 && expr.intValue < 256 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, expr.intValue)
	} else if expr.intValue >= 0 && expr.intValue < 65536 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_2BYTE, expr.intValue)
	} else {
		c := vm.NewConstantInt(expr.intValue)
		cpIdx := exe.AddConstantPool(c)

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

func (expr *DoubleExpression) show(indent int) {
	printWithIndent("DoubleExpr", indent)
}

func (expr *DoubleExpression) fix(currentBlock *Block) Expression {
	expr.setType(newTypeSpecifier(vm.DoubleType))
	expr.typeS().fix()
	return expr
}
func (expr *DoubleExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	if expr.doubleValue == 0.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_DOUBLE_0)

	} else if expr.doubleValue == 1.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_DOUBLE_1)

	} else {
		c := vm.NewConstantDouble(expr.doubleValue)
		cpIdx := exe.AddConstantPool(c)

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

func (expr *StringExpression) show(indent int) {
	printWithIndent("StringExpr", indent)
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.setType(newTypeSpecifier(vm.StringType))
	expr.typeS().fix()
	return expr
}

func (expr *StringExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	c := vm.NewConstantString(expr.stringValue)
	cpIdx := exe.AddConstantPool(c)

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

func (expr *NullExpression) show(indent int) {
	printWithIndent("NullExpr", indent)
}

func (expr *NullExpression) fix(currentBlock *Block) Expression {
	expr.setType(newTypeSpecifier(vm.NullType))
	expr.typeS().fix()
	return expr
}

func (expr *NullExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

	// 声明要么是变量，要么是函数, 要么是包(FunctionIdentifier Declaration Module)
	inner IdentifierInner
}

func (expr *IdentifierExpression) show(indent int) {
	printWithIndent("IdentifierExpr", indent)
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {
	// 判断是否是变量
	declaration := searchDeclaration(expr.name, currentBlock)
	if declaration != nil {
		expr.setType(declaration.typeSpecifier)
		expr.inner = declaration
		expr.typeS().fix()
		return expr
	}

	// 判断是否是函数
	fd := searchFunction(expr.name)
	if fd != nil {
		compiler := getCurrentCompiler()

		expr.setType(createFunctionDeriveType(fd))
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			functionIndex:      compiler.addToVmFunctionList(fd),
		}
		expr.typeS().fix()

		return expr
	}

	// TODO 判断是否是包
	module := searchModule(expr.name)
	if module != nil {
		expr.setType(module.typ)
		expr.inner = module
		expr.typeS().fix()
		return expr
	}

	// 都不是,报错
	compileError(expr.Position(), IDENTIFIER_NOT_FOUND_ERR, expr.name)
	return nil
}

func (expr *IdentifierExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	switch inner := expr.inner.(type) {
	// 函数
	case *FunctionIdentifier:
		ob.generateCode(expr.Position(), vm.VM_PUSH_FUNCTION, inner.functionIndex)
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

func (expr *AssignExpression) show(indent int) {
	printWithIndent("AssignExpr", indent)

	subIndent := indent + 2
	expr.left.show(subIndent)
	expr.operand.show(subIndent)
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

func (expr *AssignExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	expr.generateEx(exe, currentBlock, ob, false)
}

// 顶层
func (expr *AssignExpression) generateEx(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf, isTopLevel bool) {
	expr.operand.generate(exe, currentBlock, ob)

	if !isTopLevel {
		ob.generateCode(expr.Position(), vm.VM_DUPLICATE)
	}

	generatePopToLvalue(exe, currentBlock, expr.left, ob)
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

func (expr *BinaryExpression) show(indent int) {
	printWithIndent("BinaryExpr", indent)

	subIndent := indent + 2
	expr.left.show(subIndent)
	expr.right.show(subIndent)
}

func (expr *BinaryExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	switch expr.operator {
	// 数学计算
	case AddOperator, SubOperator, MulOperator, DivOperator:
		newExpr = fixMathBinaryExpression(expr, currentBlock)
		// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		newExpr = fixCompareBinaryExpression(expr, currentBlock)
		// && ||
	case LogicalAndOperator, LogicalOrOperator:
		newExpr = fixLogicalBinaryExpression(expr, currentBlock)
	default:
		panic("TODO")
	}

	newExpr.typeS().fix()

	return newExpr
}

func (expr *BinaryExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

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

func (expr *MinusExpression) show(indent int) {
	printWithIndent("MinusExpr", indent)

	subIndent := indent + 2
	expr.operand.show(subIndent)
}

func (expr *MinusExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.operand = expr.operand.fix(currentBlock)

	if !isInt(expr.operand.typeS()) && !isDouble(expr.operand.typeS()) {
		compileError(expr.Position(), MINUS_TYPE_MISMATCH_ERR, "")
	}

	expr.setType(expr.operand.typeS())

	switch operand := expr.operand.(type) {
	case *IntExpression:
		operand.intValue = -operand.intValue
		newExpr = operand
	case *DoubleExpression:
		operand.doubleValue = -operand.doubleValue
		newExpr = operand
	default:
		newExpr = expr
	}

	newExpr.typeS().fix()

	return newExpr
}

func (expr *MinusExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

func (expr *LogicalNotExpression) show(indent int) {
	printWithIndent("LogicalNotExpr", indent)

	subIndent := indent + 2

	expr.operand.show(subIndent)
}

func (expr *LogicalNotExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.operand = expr.operand.fix(currentBlock)

	switch operand := expr.operand.(type) {
	case *BooleanExpression:
		operand.booleanValue = !operand.booleanValue
		operand.setType(createTypeSpecifier(vm.BooleanType, expr.Position()))
		newExpr = operand
	default:
		if !isBoolean(expr.operand.typeS()) {
			compileError(expr.Position(), LOGICAL_NOT_TYPE_MISMATCH_ERR, "")
		}
		expr.setType(expr.operand.typeS())
		newExpr = expr
	}

	newExpr.typeS().fix()

	return newExpr
}

func (expr *LogicalNotExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

func (expr *FunctionCallExpression) show(indent int) {
	printWithIndent("FuncCallExpr", indent)

	subIndent := indent + 2

	expr.function.show(subIndent)
	for _, arg := range expr.argumentList {
		printWithIndent("ArgList", subIndent)
		arg.show(subIndent + 2)
	}
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	var fd *FunctionDefinition
	var arrayBase *TypeSpecifier
	var name string

	funcIfs := expr.function.fix(currentBlock)

	expr.function = funcIfs

	switch funcExpr := funcIfs.(type) {
	case *IdentifierExpression:
		fd = searchFunction(funcExpr.name)
		name = funcExpr.name
	}

	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, name)
	}

	fd.checkArgument(currentBlock, expr.argumentList, arrayBase)

	expr.setType(newTypeSpecifier(fd.typeS().basicType))

	expr.typeSpecifier.deriveType = fd.typeS().deriveType

	expr.typeS().fix()
	return expr
}

func (expr *FunctionCallExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	generatePushArgument(expr.argumentList, exe, currentBlock, ob)

	expr.function.generate(exe, currentBlock, ob)

	ob.generateCode(expr.Position(), vm.VM_INVOKE)
}

// ==============================
// MemberExpression
// ==============================
type MemberExpression struct {
	ExpressionImpl

	// 实例
	expression Expression

	// 成员名称
	memberName string

	// module func
	moduleFunc *FunctionDefinition
}

func (expr *MemberExpression) show(indent int) {
	printWithIndent("MemberExpr", indent)

	subIndent := indent + 2
	expr.expression.show(subIndent)
}

func (expr *MemberExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.expression = expr.expression.fix(currentBlock)

	typ := expr.expression.typeS()

	switch {
	case typ.isModule():
		newExpr = fixModuleMemberExpression(expr, expr.memberName)
	default:
		compileError(expr.Position(), MEMBER_EXPRESSION_TYPE_ERR)
	}

	newExpr.typeS().fix()

	return newExpr
}

func (expr *MemberExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
}

func createMemberExpression(expression Expression, memberName string) *MemberExpression {
	expr := &MemberExpression{
		expression: expression,
		memberName: memberName,
	}
	expr.SetPosition(expression.Position())

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

func (expr *CastExpression) show(indent int) {
	printWithIndent("CastExpr", indent)
}

func (expr *CastExpression) fix(currentBlock *Block) Expression { return expr }

func (expr *CastExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

func (expr *ArrayLiteralExpression) show(indent int) {
	printWithIndent("ArrayLiteralExpr", indent)
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

	expr.setType(newTypeSpecifier(elemType.basicType))

	expr.typeS().deriveType = &ArrayDerive{}
	expr.typeS().sliceType = NewSliceType(elemType)

	expr.typeS().fix()

	return expr
}

func (expr *ArrayLiteralExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	if !expr.typeS().isArrayDerive() {
		panic("TODO")
	}

	// TODO: 可以创建空
	if len(expr.arrayLiteral) == 0 {
		panic("TODO")
	}

	for _, subExpr := range expr.arrayLiteral {
		subExpr.generate(exe, currentBlock, ob)
	}

	itemType := expr.arrayLiteral[0].typeS()
	offset := getOpcodeTypeOffset(itemType)

	count := len(expr.arrayLiteral)
	ob.generateCode(expr.Position(), vm.VM_NEW_ARRAY_LITERAL_INT+offset, count)
}

// ==============================
// IndexExpression
// ==============================
type IndexExpression struct {
	ExpressionImpl

	array Expression
	index Expression
}

func (expr *IndexExpression) show(indent int) {
	printWithIndent("IndexExpr", indent)

	subIndent := indent + 2
	expr.array.show(subIndent)
	expr.index.show(subIndent)
}

func (expr *IndexExpression) fix(currentBlock *Block) Expression {

	expr.array = expr.array.fix(currentBlock)
	expr.index = expr.index.fix(currentBlock)

	if !expr.array.typeS().isArrayDerive() {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	expr.setType(cloneTypeSpecifier(expr.array.typeS()))

	expr.typeS().deriveType = nil

	if !isInt(expr.index.typeS()) {
		compileError(expr.Position(), INDEX_NOT_INT_ERR)
	}

	expr.typeS().fix()

	return expr
}

func (expr *IndexExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

//
// TODO Module
type Module struct {
	typ      *TypeSpecifier
	compiler *Compiler
}
