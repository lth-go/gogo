package compiler

import (
	"fmt"

	"github.com/lth-go/gogo/vm"
)

//
// Expression interface
//
type Expression interface {
	Pos
	Fix() Expression     // 用于类型修正,以及简单的类型转换
	Generate(*OpCodeBuf) // 生成字节码
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

func (expr *ExpressionBase) Fix() Expression        { return nil }
func (expr *ExpressionBase) Generate(ob *OpCodeBuf) {}
func (expr *ExpressionBase) GetType() *Type         { return expr.Type }
func (expr *ExpressionBase) SetType(t *Type)        { expr.Type = t }

//
// BoolExpression
//
type BoolExpression struct {
	ExpressionBase
	Value bool
}

func (expr *BoolExpression) Fix() Expression {
	expr.SetType(NewType(vm.BasicTypeBool))
	expr.GetType().Fix()
	return expr
}

func (expr *BoolExpression) Generate(ob *OpCodeBuf) {
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
	Index int
}

func (expr *IntExpression) Fix() Expression {
	expr.SetType(NewType(vm.BasicTypeInt))
	expr.GetType().Fix()

	if expr.Value > 65535 || expr.Value < 0 {
		expr.Index = GetCurrentCompiler().AddConstantList(expr.Value)
	}

	return expr
}

func (expr *IntExpression) Generate(ob *OpCodeBuf) {
	if expr.Value >= 0 && expr.Value < 256 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_1BYTE, expr.Value)
	} else if expr.Value >= 0 && expr.Value < 65536 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT_2BYTE, expr.Value)
	} else {
		ob.generateCode(expr.Position(), vm.VM_PUSH_INT, expr.Index)
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
	Index int
}

func (expr *FloatExpression) Fix() Expression {
	expr.SetType(NewType(vm.BasicTypeFloat))
	expr.GetType().Fix()

	if expr.Value != 0.0 && expr.Value != 1.0 {
		expr.Index = GetCurrentCompiler().AddConstantList(expr.Value)
	}
	return expr
}

func (expr *FloatExpression) Generate(ob *OpCodeBuf) {
	if expr.Value == 0.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_0)
	} else if expr.Value == 1.0 {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT_1)
	} else {
		ob.generateCode(expr.Position(), vm.VM_PUSH_FLOAT, expr.Index)
	}
}

func CreateFloatExpression(pos Position, value float64) *FloatExpression {
	expr := &FloatExpression{
		Value: value,
		Index: -1,
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
	Index int
}

func (expr *StringExpression) Fix() Expression {
	expr.SetType(NewType(vm.BasicTypeString))
	expr.GetType().Fix()

	expr.Index = GetCurrentCompiler().AddConstantList(expr.Value)
	return expr
}

func (expr *StringExpression) Generate(ob *OpCodeBuf) {
	ob.generateCode(expr.Position(), vm.VM_PUSH_STRING, expr.Index)
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

func (expr *NilExpression) Fix() Expression {
	expr.SetType(NewType(vm.BasicTypeNil))
	expr.GetType().Fix()

	return expr
}

func (expr *NilExpression) Generate(ob *OpCodeBuf) {
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
	Block *Block
}

type FunctionIdentifier struct {
	functionDefinition *FunctionDefinition
	Index              int
}

func (expr *IdentifierExpression) Fix() Expression {
	// 判断是否是变量
	declaration := expr.Block.SearchDeclaration(expr.name)
	if declaration != nil {
		expr.SetType(declaration.Type.Copy())
		expr.inner = declaration
		expr.GetType().Fix()
		return expr
	}

	// 判断是否是函数
	fd := GetCurrentCompiler().SearchFunction(expr.name)
	if fd != nil {
		expr.SetType(fd.CopyType())
		expr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			Index:              GetCurrentCompiler().GetFuncIndex(fd),
		}
		expr.GetType().Fix()

		return expr
	}

	// TODO 判断是否是包
	pkg := GetCurrentCompiler().SearchPackage(expr.name)
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

func (expr *IdentifierExpression) Generate(ob *OpCodeBuf) {
	switch inner := expr.inner.(type) {
	// 函数
	case *FunctionIdentifier:
		ob.generateCode(expr.Position(), vm.VM_PUSH_FUNCTION, inner.Index)
		// 变量
	case *Declaration:
		var code byte

		if inner.IsLocal {
			code = vm.VM_PUSH_STACK
		} else {
			code = vm.VM_PUSH_STATIC
		}
		ob.generateCode(expr.Position(), code, inner.Index)
	}
}

func CreateIdentifierExpression(pos Position, name string) *IdentifierExpression {
	expr := &IdentifierExpression{name: name}
	expr.SetPosition(pos)
	expr.Block = GetCurrentCompiler().currentBlock
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

func (expr *BinaryExpression) Fix() Expression {
	var newExpr Expression

	switch expr.operator {
	// 数学计算
	case AddOperator, SubOperator, MulOperator, DivOperator:
		newExpr = FixMathBinaryExpression(expr)
	// 比较
	case EqOperator, NeOperator, GtOperator, GeOperator, LtOperator, LeOperator:
		newExpr = FixCompareBinaryExpression(expr)
	// && ||
	case LogicalAndOperator, LogicalOrOperator:
		newExpr = FixLogicalBinaryExpression(expr)
	default:
		panic("TODO")
	}

	newExpr.GetType().Fix()

	return newExpr
}

// TODO: 确定golang如何比较
func (expr *BinaryExpression) Generate(ob *OpCodeBuf) {

	switch operator := expr.operator; operator {
	case GtOperator, GeOperator, LtOperator, LeOperator,
		AddOperator, SubOperator, MulOperator, DivOperator,
		EqOperator, NeOperator:

		var offset byte

		leftExpr := expr.left
		rightExpr := expr.right

		leftExpr.Generate(ob)
		rightExpr.Generate(ob)

		code, ok := operatorCodeMap[operator]
		if !ok {
			panic("TODO")
		}

		if leftExpr.GetType().IsInt() || leftExpr.GetType().IsBool() {
			offset = byte(0)
		} else if leftExpr.GetType().IsFloat() {
			offset = byte(1)
		} else if leftExpr.GetType().IsString() {
			offset = byte(2)
		} else if leftExpr.GetType().IsNil() || rightExpr.GetType().IsNil() {
			offset = byte(3)
		} else if leftExpr.GetType().IsComposite() {
			offset = byte(3)
		} else {
			panic("TODO")
		}

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

		expr.left.Generate(ob)
		ob.generateCode(expr.Position(), vm.VM_DUPLICATE)
		ob.generateCode(expr.Position(), jumpCode, label)

		expr.right.Generate(ob)

		ob.generateCode(expr.Position(), logicalCode)

		ob.setLabel(label)
	}
}

type UnaryExpression struct {
	ExpressionBase
	Operator UnaryOperatorKind
	Value    Expression
}

func (expr *UnaryExpression) Fix() Expression {
	switch expr.Operator {
	case UnaryOperatorKindMinus:
		return expr.FixMinus()
	case UnaryOperatorKindNot:
		return expr.FixNot()
	default:
		panic("TODO")
	}
}

func (expr *UnaryExpression) Generate(ob *OpCodeBuf) {
	switch expr.Operator {
	case UnaryOperatorKindMinus:
		expr.GenerateMinus(ob)
	case UnaryOperatorKindNot:
		expr.GenerateNot(ob)
	default:
		panic("TODO")
	}
}

func (expr *UnaryExpression) FixMinus() Expression {
	var newExpr Expression

	expr.Value = expr.Value.Fix()

	if !expr.Value.GetType().IsInt() && !expr.Value.GetType().IsFloat() {
		compileError(expr.Position(), MINUS_TYPE_MISMATCH_ERR, "")
	}

	expr.SetType(expr.Value.GetType().Copy())

	// 如果值是常量,则直接转换
	switch operand := expr.Value.(type) {
	case *IntExpression:
		operand.Value = -operand.Value
		newExpr = operand
		newExpr = newExpr.Fix()
	case *FloatExpression:
		operand.Value = -operand.Value
		newExpr = operand
		newExpr = newExpr.Fix()
	default:
		newExpr = expr
	}

	newExpr.GetType().Fix()

	return newExpr
}

func (expr *UnaryExpression) FixNot() Expression {
	var newExpr Expression

	expr.Value = expr.Value.Fix()

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

func (expr *UnaryExpression) GenerateMinus(ob *OpCodeBuf) {
	expr.Value.Generate(ob)
	code := vm.VM_MINUS_INT
	if expr.GetType().IsFloat() {
		code = vm.VM_MINUS_FLOAT
	}
	ob.generateCode(expr.Position(), code)
}

func (expr *UnaryExpression) GenerateNot(ob *OpCodeBuf) {
	expr.Value.Generate(ob)
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

func (expr *FunctionCallExpression) Fix() Expression {
	var fd *FunctionDefinition
	var name string

	expr.funcName = expr.funcName.Fix()

	switch funcNameExpr := expr.funcName.(type) {
	case *IdentifierExpression:
		fd = funcNameExpr.inner.(*FunctionIdentifier).functionDefinition
		name = funcNameExpr.name
	default:
		compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, name)
	}

	fd.FixArgument(expr.argumentList)

	// 设置返回值类型
	resultCount := len(fd.Type.funcType.Results)

	if resultCount == 0 {
		expr.GetType().SetBasicType(vm.BasicTypeVoid)
	} else if resultCount == 1 {
		expr.GetType().SetBasicType(fd.Type.funcType.Results[0].Type.GetBasicType())
	} else {
		typeList := make([]*Type, resultCount)
		for i, resultType := range fd.Type.funcType.Results {
			typeList[i] = resultType.Type.Copy()
		}
		expr.GetType().SetBasicType(vm.BasicTypeMultipleValues)
		expr.GetType().multipleValueType = NewMultipleValueType(typeList)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *FunctionCallExpression) Generate(ob *OpCodeBuf) {
	generatePushArgument(expr.argumentList, ob)
	expr.funcName.Generate(ob)
	ob.generateCode(expr.Position(), vm.VM_INVOKE)
}

func NewFunctionCallExpression(pos Position, function Expression, argumentList []Expression) *FunctionCallExpression {
	expr := &FunctionCallExpression{
		funcName:     function,
		argumentList: argumentList,
	}
	expr.SetType(NewType(vm.BasicTypeVoid))
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

func (expr *SelectorExpression) Fix() Expression {
	var newExpr Expression

	expr.expression = expr.expression.Fix()
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

func (expr *SelectorExpression) Generate(ob *OpCodeBuf) {}

// 仅限函数
func fixPackageSelectorExpression(expr *SelectorExpression, field string) Expression {
	innerExpr := expr.expression
	innerExpr.GetType().Fix()

	p := innerExpr.(*IdentifierExpression).inner.(*Package)

	fd := p.compiler.SearchFunction(field)
	if fd != nil {
		newExpr := CreateIdentifierExpression(expr.Position(), field)
		newExpr.inner = &FunctionIdentifier{
			functionDefinition: fd,
			Index:              GetCurrentCompiler().GetFuncIndex(fd),
		}

		newExpr.SetType(fd.CopyType())
		newExpr.GetType().Fix()

		return newExpr
	}

	decl := p.compiler.SearchDeclaration(field)
	if decl != nil {
		// TODO: 初始值直接给会有问题
		newDecl := NewDeclaration(decl.Position(), decl.Type.Copy(), decl.Name, nil)
		newDecl.PackageName = p.compiler.GetPackageName()
		newDecl.Index = GetCurrentCompiler().AddDeclarationList(newDecl)

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

func (expr *CastExpression) Fix() Expression { return expr }

func (expr *CastExpression) Generate(ob *OpCodeBuf) {
	expr.operand.Generate(ob)

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

func (expr *ArrayExpression) Fix() Expression {
	elemType := expr.GetType().arrayType.ElementType

	for i := 1; i < len(expr.List); i++ {
		expr.List[i] = expr.List[i].Fix()
		expr.List[i] = CreateAssignCast(expr.List[i], elemType)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *ArrayExpression) Generate(ob *OpCodeBuf) {
	for _, subExpr := range expr.List {
		subExpr.Generate(ob)
	}

	count := len(expr.List)
	ob.generateCode(expr.Position(), vm.VM_NEW_ARRAY, count)
}

func CreateArrayExpression(pos Position, typ *Type, exprList []Expression) *ArrayExpression {
	expr := &ArrayExpression{
		List: exprList,
	}
	expr.SetType(typ)
	expr.SetPosition(pos)

	return expr
}

// MapExpression
type MapExpression struct {
	ExpressionBase
	KeyList   []Expression
	ValueList []Expression
}

func (expr *MapExpression) Fix() Expression {
	if len(expr.KeyList) != len(expr.ValueList) {
		panic("TODO")
	}

	keyType := expr.GetType().mapType.Key
	valueType := expr.GetType().mapType.Value

	for i := 1; i < len(expr.KeyList); i++ {
		expr.KeyList[i] = expr.KeyList[i].Fix()
		expr.KeyList[i] = CreateAssignCast(expr.KeyList[i], keyType)
	}

	for i := 1; i < len(expr.ValueList); i++ {
		expr.ValueList[i] = expr.ValueList[i].Fix()
		expr.ValueList[i] = CreateAssignCast(expr.ValueList[i], valueType)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *MapExpression) Generate(ob *OpCodeBuf) {
	for _, subExpr := range expr.ValueList {
		subExpr.Generate(ob)
	}

	for _, subExpr := range expr.KeyList {
		subExpr.Generate(ob)
	}

	size := len(expr.KeyList)

	ob.generateCode(expr.Position(), vm.VM_NEW_MAP, size)
}

func CreateMapExpression(pos Position, typ *Type) *MapExpression {
	expr := &MapExpression{}
	expr.SetType(typ)
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

func (expr *IndexExpression) Fix() Expression {

	expr.array = expr.array.Fix()
	expr.index = expr.index.Fix()

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

func (expr *IndexExpression) Generate(ob *OpCodeBuf) {
	expr.array.Generate(ob)
	expr.index.Generate(ob)

	code := vm.VM_PUSH_ARRAY
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
