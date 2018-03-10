package compiler

import (
	"../vm"
)

//
// BinaryOperatorKind
//
type BinaryOperatorKind int

const (
	LogicalOrOperator  BinaryOperatorKind = iota
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
	IntToStringCast     CastType = iota
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

func (expr *ExpressionImpl) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {}

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
	printWithIndent("BoolExpr", ident)
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.BooleanType})
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

func (expr *IntExpression) show(ident int) {
	printWithIndent("IntExpr", ident)
}

func (expr *IntExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.IntType})
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

func (expr *DoubleExpression) show(ident int) {
	printWithIndent("DoubleExpr", ident)
}

func (expr *DoubleExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.DoubleType})
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

func (expr *StringExpression) show(ident int) {
	printWithIndent("StringExpr", ident)
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.StringType})
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

func (expr *NullExpression) show(ident int) {
	printWithIndent("NullExpr", ident)
}

func (expr *NullExpression) fix(currentBlock *Block) Expression {
	expr.setType(&TypeSpecifier{basicType: vm.NullType})
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

	// 声明要么是变量，要么是函数 (FunctionIdentifier Declaration)
	inner IdentifierInner
}

func (expr *IdentifierExpression) show(ident int) {
	printWithIndent("IdentifierExpr", ident)
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

		expr.setType(createFunctionDeriveType(fd))
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			functionIndex:      compiler.addToVmFunctionList(fd),
		}
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
// CommaExpression
// ==============================

// CommaExpression 逗号表达式
type CommaExpression struct {
	ExpressionImpl

	left  Expression
	right Expression
}

func (expr *CommaExpression) show(ident int) {
	printWithIndent("CommaExpr", ident)

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

func (expr *CommaExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

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
	printWithIndent("AssignExpr", ident)

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

func (expr *BinaryExpression) show(ident int) {
	printWithIndent("BinaryExpr", ident)

	subIdent := ident + 2
	expr.left.show(subIdent)
	expr.right.show(subIdent)
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

func (expr *MinusExpression) show(ident int) {
	printWithIndent("MinusExpr", ident)

	subIdent := ident + 2
	expr.operand.show(subIdent)
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

func (expr *LogicalNotExpression) show(ident int) {
	printWithIndent("LogicalNotExpr", ident)

	subIdent := ident + 2

	expr.operand.show(subIdent)
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

func (expr *FunctionCallExpression) show(ident int) {
	printWithIndent("FuncCallExpr", ident)

	subIdent := ident + 2

	expr.function.show(subIdent)
	for _, arg := range expr.argumentList {
		printWithIndent("ArgList", subIdent)
		arg.show(subIdent + 2)
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
		fd = SearchFunction(funcExpr.name)
		name = funcExpr.name
	case *MemberExpression:
		switch member := funcExpr.declaration.(type) {
		case *FieldMember:
			compileError(expr.Position(), FIELD_CAN_NOT_CALL_ERR, member.name)
		case *MethodMember:
			fd = member.functionDefinition
			name = fd.name
		default:
			panic("TODO")
		}
	}

	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, name)
	}

	fd.checkArgument(currentBlock, expr.argumentList, arrayBase)

	expr.setType(&TypeSpecifier{basicType: fd.typeS().basicType})

	expr.typeSpecifier.deriveList = fd.typeS().deriveList

	if expr.typeS().basicType == vm.ClassType {
		expr.typeS().classRef.identifier = fd.typeS().classRef.identifier
		expr.typeS().fix()
	}

	expr.typeS().fix()
	return expr
}
func (expr *FunctionCallExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	switch memberExpr := expr.function.(type) {
	case *MemberExpression:
		_, ok := memberExpr.declaration.(*MethodMember)
		if ok {
			generateMethodCallExpression(expr, exe, currentBlock, ob)
			return
		}
	}

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

	declaration MemberDeclaration
	methodIndex int
}

func (expr *MemberExpression) show(ident int) {
	printWithIndent("MemberExpr", ident)

	subIdent := ident + 2
	expr.expression.show(subIdent)
}

func (expr *MemberExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.expression = expr.expression.fix(currentBlock)
	obj := expr.expression

	if isClass(obj.typeS()) {
		newExpr = fixClassMemberExpression(expr, obj, expr.memberName)
	} else {
		compileError(expr.Position(), MEMBER_EXPRESSION_TYPE_ERR)
	}

	newExpr.typeS().fix()

	return newExpr
}

func (expr *MemberExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	switch member := expr.declaration.(type) {
	case *FieldMember:
		expr.expression.generate(exe, currentBlock, ob)
		ob.generateCode(expr.Position(), vm.VM_PUSH_FIELD_INT+getOpcodeTypeOffset(expr.typeS()), member.fieldIndex)
	case *MethodMember:
		compileError(expr.Position(), METHOD_IS_NOT_CALLED_ERR, member.functionDefinition.name)
	}
}

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

func (expr *ThisExpression) show(ident int) {
	printWithIndent("ThisExpr", ident)
}

func (expr *ThisExpression) fix(currentBlock *Block) Expression {

	cd := getCurrentCompiler().currentClassDefinition

	if cd == nil {
		compileError(expr.Position(), THIS_OUT_OF_CLASS_ERR)
	}

	typ := &TypeSpecifier{basicType: vm.ClassType}
	typ.classRef = classRef{
		identifier:      cd.name,
		classDefinition: cd,
	}
	expr.setType(typ)

	expr.typeS().fix()

	return expr
}

func (expr *ThisExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	fd := currentBlock.getCurrentFunction()
	paramCount := len(fd.parameterList)
	ob.generateCode(expr.Position(), vm.VM_PUSH_STACK_OBJECT, paramCount)
}

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

func (expr *ArrayLiteralExpression) show(ident int) {
	printWithIndent("ArrayLiteralExpr", ident)
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

	expr.typeS().fix()

	return expr
}

func (expr *ArrayLiteralExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	if !expr.typeS().isArrayDerive() {
		panic("TODO")
	}

	if expr.arrayLiteral == nil {
		panic("TODO")
	}

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
	printWithIndent("ArrayCreationExpr", ident)
}

func (expr *ArrayCreation) fix(currentBlock *Block) Expression {
	expr.typeS().fix()

	deriveList := []TypeDerive{}

	for _, dim := range expr.dimensionList {
		if dim.expression != nil {
			dim.expression = dim.expression.fix(currentBlock)

			if !isInt(dim.expression.typeS()) {
				compileError(expr.Position(), ARRAY_SIZE_NOT_INT_ERR)
			}
		}
		deriveList = append([]TypeDerive{&ArrayDerive{}}, deriveList...)
	}

	expr.setType(cloneTypeSpecifier(expr.typeS()))
	expr.typeS().deriveList = deriveList

	expr.typeS().fix()

	return expr
}

func (expr *ArrayCreation) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
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

func createBasicArrayCreation(typ *TypeSpecifier, dimExprLists, dimList []*ArrayDimension, pos Position) Expression {
	expr := createClassArrayCreation(typ, dimExprLists, dimList, pos)

	return expr
}

func createClassArrayCreation(typ *TypeSpecifier, dimExprList, dimList []*ArrayDimension, pos Position) Expression {

	expr := &ArrayCreation{
		dimensionList: dimExprList,
	}

	expr.SetPosition(pos)

	expr.setType(typ)

	expr.dimensionList = append(expr.dimensionList, dimList...)

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
	printWithIndent("IndexExpr", ident)

	subIdent := ident + 2
	expr.array.show(subIdent)
	expr.index.show(subIdent)
}

func (expr *IndexExpression) fix(currentBlock *Block) Expression {

	expr.array = expr.array.fix(currentBlock)
	expr.index = expr.index.fix(currentBlock)

	if !expr.array.typeS().isArrayDerive() {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	expr.setType(cloneTypeSpecifier(expr.array.typeS()))

	expr.typeS().deriveList = expr.array.typeS().deriveList[1:]

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

// ==============================
// NewExpression
// ==============================
type NewExpression struct {
	ExpressionImpl

	// 类名
	className       string
	classDefinition *ClassDefinition
	classIndex      int

	// 类初始化方便名
	methodName        string
	// 类初始化方法
	methodDeclaration MemberDeclaration
	// 参数
	argumentList      []Expression
}

func (expr *NewExpression) fix(currentBlock *Block) Expression {
	expr.classDefinition = searchClassAndAdd(expr.Position(), expr.className, &expr.classIndex)

	if expr.methodName == "" {
		expr.methodName = defaultConstructorName
	}

	member := expr.classDefinition.searchMember(expr.methodName)

	if member == nil {
		compileError(expr.Position(), MEMBER_NOT_FOUND_ERR, expr.className, expr.methodName)
	}

	methodMember, ok := member.(*MethodMember)
	if !ok {
		compileError(expr.Position(), CONSTRUCTOR_IS_FIELD_ERR, expr.methodName)
	}

	//check_member_accessibility(expr.Position(), expr.classDefinition, member, expr.methodName)

	if !(methodMember.functionDefinition.typeS().deriveList == nil && methodMember.functionDefinition.typeS().basicType == vm.VoidType) {
		panic("TODO")
	}

	methodMember.functionDefinition.checkArgument(currentBlock, expr.argumentList, nil)

	expr.methodDeclaration = member
	typ := &TypeSpecifier{
		basicType: vm.ClassType,
		classRef: classRef{
			identifier:      expr.classDefinition.name,
			classDefinition: expr.classDefinition,
		},
	}
	expr.setType(typ)

	expr.typeS().fix()

	return expr
}

func (expr *NewExpression) generate(exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {

	methodMember := expr.methodDeclaration.(*MethodMember)
	paramCount := len(methodMember.functionDefinition.parameterList)

	ob.generateCode(expr.Position(), vm.VM_NEW, expr.classIndex)
	generatePushArgument(expr.argumentList, exe, currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_DUPLICATE_OFFSET, paramCount)

	ob.generateCode(expr.Position(), vm.VM_PUSH_METHOD, methodMember.methodIndex)
	ob.generateCode(expr.Position(), vm.VM_INVOKE)
	ob.generateCode(expr.Position(), vm.VM_POP)
}

func createNewExpression(className, memthodName string, argumentList []Expression, pos Position) *NewExpression {
	expr := &NewExpression{
		className:    className,
		methodName:   memthodName,
		argumentList: argumentList,
	}

	expr.SetPosition(pos)

	return expr
}
