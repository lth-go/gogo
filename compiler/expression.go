package compiler

import (
	"fmt"

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

type UnaryOperatorKind int

const (
	UnaryOperatorKindMinus UnaryOperatorKind = iota
	UnaryOperatorKindNot
)

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
// BoolExpression
//
type BoolExpression struct {
	ExpressionBase
	Value bool
}

func (expr *BoolExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeBool))
	expr.GetType().Fix()
	return expr
}

func (expr *BoolExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	value := 0
	if expr.Value {
		value = 1
	}

	ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, value)
}

func CreateBooleanExpression(pos Position, value bool) *BoolExpression {
	expr := &BoolExpression{
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

//
// IntExpression 数字表达式
//
type IntExpression struct {
	ExpressionBase
	Value int
}

func (expr *IntExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeInt))
	expr.GetType().Fix()
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

func CreateIntExpression(pos Position, value int) *IntExpression {
	expr := &IntExpression{
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

//
// FloatExpression
//
type FloatExpression struct {
	ExpressionBase
	Value float64
}

func (expr *FloatExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeFloat))
	expr.GetType().Fix()
	return expr
}

func (expr *FloatExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	if expr.Value == 0.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_0)
	} else if expr.Value == 1.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_1)
	} else {
		cpIdx := getCurrentCompiler().AddConstantList(expr.Value)
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT, cpIdx)
	}
}

func CreateFloatExpression(pos Position, value float64) *FloatExpression {
	expr := &FloatExpression{
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

//
// StringExpression 字符串表达式
//
type StringExpression struct {
	ExpressionBase
	Value string
}

func (expr *StringExpression) fix(currentBlock *Block) Expression {
	expr.SetType(NewType(vm.BasicTypeString))
	expr.GetType().Fix()

	return expr
}

func (expr *StringExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	cpIdx := getCurrentCompiler().AddConstantList(expr.Value)
	ob.generateCode(expr.Position(), vm.VM_PUSH_STRING, cpIdx)
}

func CreateStringExpression(pos Position, value string) *StringExpression {
	expr := &StringExpression{
		Value: value,
	}
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
	expr.GetType().Fix()

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

func isNilExpression(expr Expression) bool {
	_, ok := expr.(*NilExpression)
	return ok
}

//
// IdentifierExpression
//
// IdentifierExpression 变量表达式
type IdentifierExpression struct {
	ExpressionBase
	name  string
	inner interface{} // 声明要么是变量，要么是函数, 要么是包(FunctionIdentifier Declaration Package)
}

type FunctionIdentifier struct {
	functionDefinition *FunctionDefinition
	Index              int
}

func (expr *IdentifierExpression) fix(currentBlock *Block) Expression {
	// 判断是否是变量
	declaration := currentBlock.searchDeclaration(expr.name)
	if declaration != nil {
		expr.SetType(declaration.Type.Copy())
		expr.inner = declaration
		expr.GetType().Fix()
		return expr
	}

	// 判断是否是函数
	fd := getCurrentCompiler().searchFunction(expr.name)
	if fd != nil {
		expr.SetType(fd.CopyType())
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			Index:              getCurrentCompiler().AddFuncList(fd),
		}
		expr.GetType().Fix()

		return expr
	}

	// TODO 判断是否是包
	pkg := getCurrentCompiler().searchPackage(expr.name)
	if pkg != nil {
		expr.SetType(pkg.Type.Copy())
		expr.inner = pkg
		expr.GetType().Fix()
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

func CreateIdentifierExpression(pos Position, name string) *IdentifierExpression {
	expr := &IdentifierExpression{name: name}
	expr.SetPosition(pos)
	return expr
}

//
// BinaryExpression 二元表达式
//
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
		newExpr = FixMathBinaryExpression(expr, currentBlock)
	// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		newExpr = FixCompareBinaryExpression(expr, currentBlock)
	// && ||
	case LogicalAndOperator, LogicalOrOperator:
		newExpr = FixLogicalBinaryExpression(expr, currentBlock)
	default:
		panic("TODO")
	}

	newExpr.GetType().Fix()

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
		if (isNilExpression(leftExpr) && !isNilExpression(rightExpr)) ||
			(!isNilExpression(leftExpr) && isNilExpression(rightExpr)) {
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

type UnaryExpression struct {
	ExpressionBase
	Operator UnaryOperatorKind
	Value    Expression
}

func (expr *UnaryExpression) fix(currentBlock *Block) Expression {
	switch expr.Operator {
	case UnaryOperatorKindMinus:
		return expr.FixMinus(currentBlock)
	case UnaryOperatorKindNot:
		return expr.FixNot(currentBlock)
	default:
		panic("TODO")
	}
}

func (expr *UnaryExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	switch expr.Operator {
	case UnaryOperatorKindMinus:
		expr.GenerateMinus(currentBlock, ob)
	case UnaryOperatorKindNot:
		expr.GenerateNot(currentBlock, ob)
	default:
		panic("TODO")
	}
}

func (expr *UnaryExpression) FixMinus(currentBlock *Block) Expression {
	var newExpr Expression

	expr.Value = expr.Value.fix(currentBlock)

	if !expr.Value.GetType().IsInt() && !expr.Value.GetType().IsFloat() {
		compileError(expr.Position(), MINUS_TYPE_MISMATCH_ERR, "")
	}

	expr.SetType(expr.Value.GetType().Copy())

	// 如果值是常量,则直接转换
	switch operand := expr.Value.(type) {
	case *IntExpression:
		operand.Value = -operand.Value
		newExpr = operand
	case *FloatExpression:
		operand.Value = -operand.Value
		newExpr = operand
	default:
		newExpr = expr
	}

	newExpr.GetType().Fix()

	return newExpr
}

func (expr *UnaryExpression) FixNot(currentBlock *Block) Expression {
	var newExpr Expression

	expr.Value = expr.Value.fix(currentBlock)

	if !expr.Value.GetType().IsBool() {
		compileError(expr.Position(), LOGICAL_NOT_TYPE_MISMATCH_ERR, "")
	}

	expr.SetType(expr.Value.GetType().Copy())

	switch operand := expr.Value.(type) {
	case *BoolExpression:
		operand.Value = !operand.Value
		newExpr = operand
	default:
		newExpr = expr
	}

	newExpr.GetType().Fix()

	return newExpr
}

func (expr *UnaryExpression) GenerateMinus(currentBlock *Block, ob *OpCodeBuf) {
	expr.Value.generate(currentBlock, ob)
	code := vm.VM_MINUS_INT + getOpcodeTypeOffset(expr.GetType())
	ob.generateCode(expr.Position(), code)
}

func (expr *UnaryExpression) GenerateNot(currentBlock *Block, ob *OpCodeBuf) {
	expr.Value.generate(currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_LOGICAL_NOT)
}

func NewUnaryExpression(pos Position, operator UnaryOperatorKind, value Expression) *UnaryExpression {
	expr := &UnaryExpression{
		Operator: operator,
		Value:    value,
	}
	expr.SetPosition(pos)

	return expr
}

//
// FunctionCallExpression 函数调用表达式
//
type FunctionCallExpression struct {
	ExpressionBase
	funcName     Expression   // 函数名
	argumentList []Expression // 实参列表
}

// TODO: 函数调用有多返回值,如何处理
func (expr *FunctionCallExpression) fix(currentBlock *Block) Expression {
	var fd *FunctionDefinition
	var name string

	expr.funcName = expr.funcName.fix(currentBlock)

	switch funcNameExpr := expr.funcName.(type) {
	case *IdentifierExpression:
		fd = funcNameExpr.inner.(*FunctionIdentifier).functionDefinition
		name = funcNameExpr.name
	}

	if fd == nil {
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, name)
	}

	fd.FixArgument(currentBlock, expr.argumentList)

	// TODO: 不拷贝
	expr.SetType(fd.GetType().Copy())

	// 设置返回值类型
	if len(fd.Type.funcType.Results) == 0 {
		expr.GetType().SetBasicType(vm.BasicTypeVoid)
	} else if len(fd.Type.funcType.Results) == 1 {
		expr.GetType().SetBasicType(fd.Type.funcType.Results[0].Type.GetBasicType())
	} else {
		typeList := make([]*Type, len(fd.Type.funcType.Results))
		for i, resultType := range fd.Type.funcType.Results {
			typeList[i] = resultType.Type.Copy()
		}
		expr.GetType().SetBasicType(vm.BasicTypeMultipleValues)
		expr.GetType().multipleValueType = NewMultipleValueType(typeList)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *FunctionCallExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	generatePushArgument(expr.argumentList, currentBlock, ob)
	expr.funcName.generate(currentBlock, ob)
	ob.generateCode(expr.Position(), vm.VM_INVOKE)
}

func NewFunctionCallExpression(pos Position, function Expression, argumentList []Expression) *FunctionCallExpression {
	expr := &FunctionCallExpression{
		funcName:     function,
		argumentList: argumentList,
	}

	expr.SetPosition(pos)

	return expr
}

//
// SelectorExpression
//
type SelectorExpression struct {
	ExpressionBase
	expression Expression // 实例
	Field      string     // 成员名称
}

func (expr *SelectorExpression) fix(currentBlock *Block) Expression {
	var newExpr Expression

	expr.expression = expr.expression.fix(currentBlock)
	typ := expr.expression.GetType()

	switch {
	case typ.IsPackage():
		newExpr = fixPackageSelectorExpression(expr, expr.Field)
	default:
		compileError(expr.Position(), MEMBER_EXPRESSION_TYPE_ERR)
	}

	newExpr.GetType().Fix()

	return newExpr
}

func (expr *SelectorExpression) generate(currentBlock *Block, ob *OpCodeBuf) {}

// 仅限函数
func fixPackageSelectorExpression(expr *SelectorExpression, field string) Expression {
	innerExpr := expr.expression
	innerExpr.GetType().Fix()

	p := innerExpr.(*IdentifierExpression).inner.(*Package)

	fd := p.compiler.searchFunction(field)
	if fd != nil {
		newExpr := CreateIdentifierExpression(expr.Position(), field)
		newExpr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			Index:              getCurrentCompiler().AddFuncList(fd),
		}

		newExpr.SetType(fd.CopyType())
		newExpr.GetType().Fix()

		return newExpr
	}

	decl := p.compiler.SearchDeclaration(field)
	if decl != nil {
		// TODO: 初始值直接给会有问题
		newDecl := NewDeclaration(decl.Position(), decl.Type.Copy(), decl.Name, nil)
		newDecl.PackageName = p.compiler.packageName
		newDecl.Index = getCurrentCompiler().AddDeclarationList(newDecl)

		newExpr := CreateIdentifierExpression(expr.Position(), field)
		newExpr.SetType(newDecl.Type.Copy())
		newExpr.inner = newDecl
		expr.GetType().Fix()

		return newExpr
	}

	panic(fmt.Sprintf("package filed not found '%s'", field))
}

func CreateSelectorExpression(expression Expression, memberName string) *SelectorExpression {
	expr := &SelectorExpression{
		expression: expression,
		Field:      memberName,
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
// ArrayExpression 创建列表时的值, eg:{1,2,3,4}
//
type ArrayExpression struct {
	ExpressionBase
	List []Expression
}

func (expr *ArrayExpression) fix(currentBlock *Block) Expression {
	// TODO: 可以为空
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

	expr.SetType(NewType(vm.BasicTypeArray))
	expr.GetType().arrayType = NewArrayType(elemType)
	expr.GetType().Fix()

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

func CreateArrayExpression(pos Position, exprList []Expression) *ArrayExpression {
	expr := &ArrayExpression{
		List: exprList,
	}
	expr.SetPosition(pos)

	return expr
}

//
// IndexExpression
//
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

	expr.SetType(expr.array.GetType().arrayType.ElementType.Copy())
	expr.GetType().arrayType = nil

	if !expr.index.GetType().IsInt() {
		compileError(expr.Position(), INDEX_NOT_INT_ERR)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *IndexExpression) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr.array.generate(currentBlock, ob)
	expr.index.generate(currentBlock, ob)

	code := vm.VM_PUSH_ARRAY_OBJECT
	ob.generateCode(expr.Position(), code)
}

func CreateIndexExpression(pos Position, array, index Expression) *IndexExpression {
	expr := &IndexExpression{
		array: array,
		index: index,
	}
	expr.SetPosition(pos)

	return expr
}

type Package struct {
	Type     *Type
	compiler *Compiler
}
