package compiler

import (
	"github.com/lth-go/gogo/vm"
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
	CastTypeIntToString CastType = iota
	CastTypeBoolToString
	CastTypeFloatToString
	CastTypeIntToFloat
	CastTypeFloatToInt
)

//
// Expression interface
//
type Expression interface {
	Pos
	fix(*Block) Expression       // 用于类型修正,以及简单的类型转换
	generate(*Block, *OpCodeBuf) // 生成字节码
	GetType() *Type
	SetType(*Type)
}

//
// Expression base
//
type ExpressionBase struct {
	PosBase
	Type *Type
}

func (expr *ExpressionBase) fix(currentBlock *Block) Expression          { return nil }
func (expr *ExpressionBase) generate(currentBlock *Block, ob *OpCodeBuf) {}
func (expr *ExpressionBase) GetType() *Type                              { return expr.Type }
func (expr *ExpressionBase) SetType(t *Type)                             { expr.Type = t }

//
// BooleanExpression
//
type BooleanExpression struct {
	ExpressionBase
	Value bool
}

func (expr *BooleanExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeBool))
	expr.GetType().fix()
	return expr
}

func (expr *BooleanExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	value := 0
	if expr.Value {
		value = 1
	}

	ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, value)
}

func createBooleanExpression(pos Position) *BooleanExpression {
	expr := &BooleanExpression{}
	expr.SetPosition(pos)

	return expr
}

//
// IntExpression
//

// IntExpression 数字表达式
type IntExpression struct {
	ExpressionBase
	Value int
}

func (expr *IntExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeInt))
	expr.GetType().fix()
	return expr
}

func (expr *IntExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	if expr.Value >= 0 && expr.Value < 256 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, expr.Value)
	} else if expr.Value >= 0 && expr.Value < 65536 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_2BYTE, expr.Value)
	} else {
		cpIdx := getCurrentCompiler().AddConstantList(expr.Value)
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT, cpIdx)
	}
}

func createIntExpression(pos Position) *IntExpression {
	expr := &IntExpression{}
	expr.SetPosition(pos)

	return expr
}

//
// DoubleExpression
//
type DoubleExpression struct {
	ExpressionBase
	Value float64
}

func (expr *DoubleExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeFloat))
	expr.GetType().fix()
	return expr
}

func (expr *DoubleExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	if expr.Value == 0.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_0)
	} else if expr.Value == 1.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_1)
	} else {
		cpIdx := getCurrentCompiler().AddConstantList(expr.Value)
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT, cpIdx)
	}
}

func createDoubleExpression(pos Position) *DoubleExpression {
	expr := &DoubleExpression{}
	expr.SetPosition(pos)

	return expr
}

//
// StringExpression
//
// StringExpression 字符串表达式
type StringExpression struct {
	ExpressionBase
	Value string
}

func NewStringExpression(pos Position, value string) *StringExpression {
	expr := &StringExpression{
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeString))
	expr.GetType().fix()

	return expr
}

func (expr *StringExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	cpIdx := getCurrentCompiler().AddConstantList(expr.Value)
	ob.generateCode(expr.Position(), vm.VM_PUSH_STRING, cpIdx)
}

func createStringExpression(pos Position) *StringExpression {
	expr := &StringExpression{}
	expr.SetPosition(pos)

	return expr
}

//
// NilExpression
//
type NilExpression struct {
	ExpressionBase
}

func (expr *NilExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeNil))
	expr.GetType().fix()
	return expr
}

func (expr *NilExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	ob.generateCode(expr.Position(), vm.VM_PUSH_NIL)
}

func createNilExpression(pos Position) *NilExpression {
	expr := &NilExpression{}
	expr.SetPosition(pos)
	return expr
}

func isNull(expr Expression) bool {
	_, ok := expr.(*NilExpression)
	return ok
}

//
// IdentifierExpression
//
type FunctionIdentifier struct {
	functionDefinition *FunctionDefinition
	Index              int
}

// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionBase
	name  string
	inner interface{} // 声明要么是变量，要么是函数, 要么是包(FunctionIdentifier Declaration Package)
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {
	// 判断是否是变量
	declaration := currentBlock.searchDeclaration(expr.name)
	if declaration != nil {
		expr.SetType(declaration.Type)
		expr.inner = declaration
		expr.GetType().fix()
		return expr
	}

	// 判断是否是函数
	fd := getCurrentCompiler().searchFunction(expr.name)
	if fd != nil {
		compiler := getCurrentCompiler()

		expr.SetType(fd.CopyType())
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			Index:              compiler.AddFuncList(fd),
		}
		expr.GetType().fix()

		return expr
	}

	// TODO 判断是否是包
	pkg := getCurrentCompiler().searchPackage(expr.name)
	if pkg != nil {
		expr.SetType(pkg.typ)
		expr.inner = pkg
		expr.GetType().fix()
		return expr
	}

	// 都不是,报错
	compileError(expr.Position(), IDENTIFIER_NOT_FOUND_ERR, expr.name)
	return nil
}

func (expr *IdentifierExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	switch inner := expr.inner.(type) {
	// 函数
	case *FunctionIdentifier:
		ob.generateCode(expr.Position(), vm.VM_PUSH_FUNCTION, inner.Index)
		// 变量
	case *Declaration:
		var code byte

		offset := getOpcodeTypeOffset(inner.Type)
		if inner.IsLocal {
			code = vm.VM_PUSH_STACK_INT
		} else {
			code = vm.VM_PUSH_STATIC_INT
		}
		ob.generateCode(expr.Position(), code+offset, inner.Index)
	}
}

func createIdentifierExpression(name string, pos Position) *IdentifierExpression {
	expr := &IdentifierExpression{name: name}
	expr.SetPosition(pos)
	return expr
}

//
// BinaryExpression
//
// BinaryExpression 二元表达式
type BinaryExpression struct {
	ExpressionBase
	operator BinaryOperatorKind // 操作符
	left     Expression
	right    Expression
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

	newExpr.GetType().fix()

	return newExpr
}

func (expr *BinaryExpression) generate(currentBlock *Block, ob *OpCodeBuf) {

	switch operator := expr.operator; operator {
	case GtOperator, GeOperator, LtOperator, LeOperator,
		AddOperator, SubOperator, MulOperator, DivOperator,
		EqOperator, NeOperator:

		var offset byte

		leftExpr := expr.left
		rightExpr := expr.right

		leftExpr.generate(currentBlock, ob)
		rightExpr.generate(currentBlock, ob)

		code, ok := operatorCodeMap[operator]
		if !ok {
			// TODO
			panic("TODO")
		}

		// TODO 啥意思
		if (isNull(leftExpr) && !isNull(rightExpr)) ||
			(!isNull(leftExpr) && isNull(rightExpr)) {
			offset = byte(2)
		} else if (operator == EqOperator || operator == NeOperator) && leftExpr.GetType().IsString() {
			offset = byte(3)
		} else {
			offset = getOpcodeTypeOffset(expr.left.GetType())
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

		expr.left.generate(currentBlock, ob)
		ob.generateCode(expr.Position(), vm.VM_DUPLICATE)
		ob.generateCode(expr.Position(), jumpCode, label)

		expr.right.generate(currentBlock, ob)

		// 判断结果
		ob.generateCode(expr.Position(), logicalCode)

		ob.setLabel(label)
	}
}

//
// MinusExpression
//
// MinusExpression 负数表达式
type MinusExpression struct {
	ExpressionBase
	operand Expression
}

func (expr *MinusExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.operand = expr.operand.fix(currentBlock)

	if !expr.operand.GetType().IsInt() && !expr.operand.GetType().IsFloat() {
		compileError(expr.Position(), MINUS_TYPE_MISMATCH_ERR, "")
	}

	expr.SetType(expr.operand.GetType())

	switch operand := expr.operand.(type) {
	case *IntExpression:
		operand.Value = -operand.Value
		newExpr = operand
	case *DoubleExpression:
		operand.Value = -operand.Value
		newExpr = operand
	default:
		newExpr = expr
	}

	newExpr.GetType().fix()

	return newExpr
}

func (expr *MinusExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr.operand.generate(currentBlock, ob)
	code := vm.VM_MINUS_INT + getOpcodeTypeOffset(expr.GetType())
	ob.generateCode(expr.Position(), code)
}

//
// LogicalNotExpression
//
// LogicalNotExpression 逻辑非表达式
type LogicalNotExpression struct {
	ExpressionBase
	operand Expression
}

func (expr *LogicalNotExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.operand = expr.operand.fix(currentBlock)

	switch operand := expr.operand.(type) {
	case *BooleanExpression:
		operand.Value = !operand.Value
		operand.SetType(CreateType(vm.BasicTypeBool, expr.Position()))
		newExpr = operand
	default:
		if !expr.operand.GetType().IsBool() {
			compileError(expr.Position(), LOGICAL_NOT_TYPE_MISMATCH_ERR, "")
		}
		expr.SetType(expr.operand.GetType())
		newExpr = expr
	}

	newExpr.GetType().fix()

	return newExpr
}

func (expr *LogicalNotExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr.operand.generate(currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_LOGICAL_NOT)
}

//
// FunctionCallExpression 函数调用表达式
//
type FunctionCallExpression struct {
	ExpressionBase
	function     Expression   // 函数名
	argumentList []Expression // 实参列表
}

func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	var fd *FunctionDefinition
	var arrayBase *Type
	var name string

	funcIfs := expr.function.fix(currentBlock)

	expr.function = funcIfs

	switch funcExpr := funcIfs.(type) {
	case *IdentifierExpression:
		fd = funcExpr.inner.(*FunctionIdentifier).functionDefinition
		name = funcExpr.name
	}

	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, name)
	}

	fd.checkArgument(currentBlock, expr.argumentList, arrayBase)

	expr.SetType(fd.GetType().Copy())

	// 设置返回值类型
	// TODO: 兼容代码, 待移除
	if len(fd.Type.funcType.Results) == 0 {
		expr.GetType().basicType = vm.BasicTypeVoid
	} else {
		expr.GetType().basicType = fd.Type.funcType.Results[0].Type.basicType
	}

	expr.GetType().fix()
	return expr
}

func (expr *FunctionCallExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	generatePushArgument(expr.argumentList, currentBlock, ob)
	expr.function.generate(currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_INVOKE)
}

//
// MemberExpression
//
type MemberExpression struct {
	ExpressionBase
	expression Expression // 实例
	memberName string     // 成员名称
}

func (expr *MemberExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.expression = expr.expression.fix(currentBlock)
	typ := expr.expression.GetType()

	switch {
	case typ.IsPackage():
		newExpr = fixPackageMemberExpression(expr, expr.memberName)
	default:
		compileError(expr.Position(), MEMBER_EXPRESSION_TYPE_ERR)
	}

	newExpr.GetType().fix()

	return newExpr
}

func (expr *MemberExpression) generate(currentBlock *Block, ob *OpCodeBuf) {}

// 仅限函数
func fixPackageMemberExpression(expr *MemberExpression, memberName string) Expression {
	innerExpr := expr.expression

	innerExpr.GetType().fix()

	p := innerExpr.(*IdentifierExpression).inner.(*Package)

	fd := p.compiler.searchFunction(memberName)
	if fd == nil {
		panic("TODO")
	}

	// TODO 得用当前compiler来添加
	currentCompiler := getCurrentCompiler()

	newExpr := &IdentifierExpression{
		name: memberName,
		inner: &FunctionIdentifier{
			functionDefinition: fd,
			Index:              currentCompiler.AddFuncList(fd),
		},
	}

	newExpr.SetType(fd.CopyType())
	newExpr.GetType().fix()

	return newExpr
}

func createMemberExpression(expression Expression, memberName string) *MemberExpression {
	expr := &MemberExpression{
		expression: expression,
		memberName: memberName,
	}
	expr.SetPosition(expression.Position())

	return expr
}

//
// CastExpression
//
type CastExpression struct {
	ExpressionBase
	castType CastType
	operand  Expression
}

func (expr *CastExpression) fix(currentBlock *Block) Expression { return expr }

func (expr *CastExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr.operand.generate(currentBlock, ob)

	switch expr.castType {
	case CastTypeIntToFloat:
		ob.generateCode(expr.Position(), vm.VM_CAST_INT_TO_FLOAT)
	case CastTypeFloatToInt:
		ob.generateCode(expr.Position(), vm.VM_CAST_FLOAT_TO_INT)
	case CastTypeBoolToString:
		ob.generateCode(expr.Position(), vm.VM_CAST_BOOLEAN_TO_STRING)
	case CastTypeIntToString:
		ob.generateCode(expr.Position(), vm.VM_CAST_INT_TO_STRING)
	case CastTypeFloatToString:
		ob.generateCode(expr.Position(), vm.VM_CAST_FLOAT_TO_STRING)
	default:
		panic("TODO")
	}
}

//
// ArrayExpression
//
// 创建列表时的值, eg:{1,2,3,4}
type ArrayExpression struct {
	ExpressionBase
	List []Expression
}

func NewArrayExpression(pos Position, exprList []Expression) *ArrayExpression {
	expr := &ArrayExpression{
		List: exprList,
	}
	expr.SetPosition(pos)

	return expr
}

func (expr *ArrayExpression) fix(currentBlock *Block) Expression {
	if expr.List == nil || len(expr.List) == 0 {
		compileError(expr.Position(), ARRAY_LITERAL_EMPTY_ERR)
	}

	firstElem := expr.List[0]
	firstElem = firstElem.fix(currentBlock)

	elemType := firstElem.GetType()

	for i := 1; i < len(expr.List); i++ {
		expr.List[i] = expr.List[i].fix(currentBlock)
		expr.List[i] = CreateAssignCast(expr.List[i], elemType)
	}

	expr.SetType(NewType(vm.BasicTypeSlice))

	expr.GetType().sliceType = NewArrayType(elemType)

	expr.GetType().fix()

	return expr
}

func (expr *ArrayExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	if !expr.GetType().IsArray() {
		panic("TODO")
	}

	// TODO: 可以创建空
	if len(expr.List) == 0 {
		panic("TODO")
	}

	for _, subExpr := range expr.List {
		subExpr.generate(currentBlock, ob)
	}

	count := len(expr.List)
	ob.generateCode(expr.Position(), vm.VM_NEW_ARRAY, count)
}

// ==============================
// IndexExpression
// ==============================
type IndexExpression struct {
	ExpressionBase
	array Expression
	index Expression
}

func (expr *IndexExpression) fix(currentBlock *Block) Expression {

	expr.array = expr.array.fix(currentBlock)
	expr.index = expr.index.fix(currentBlock)

	if !expr.array.GetType().IsArray() {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	expr.SetType(expr.array.GetType().sliceType.ElementType.Copy())
	expr.GetType().sliceType = nil

	if !expr.index.GetType().IsInt() {
		compileError(expr.Position(), INDEX_NOT_INT_ERR)
	}

	expr.GetType().fix()

	return expr
}

func (expr *IndexExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr.array.generate(currentBlock, ob)
	expr.index.generate(currentBlock, ob)

	code := vm.VM_PUSH_ARRAY_OBJECT
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

type Package struct {
	typ      *Type
	compiler *Compiler
}
