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
	expr.SetType(NewType(BasicTypeBool))
	expr.GetType().Fix()
	return expr
}

func (expr *BoolExpression) Generate(ob *OpCodeBuf) {
	value := 0
	if expr.Value {
		value = 1
	}

	ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT_1BYTE, value)
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
	expr.SetType(NewType(BasicTypeInt))
	expr.GetType().Fix()

	if expr.Value > 65535 || expr.Value < 0 {
		expr.Index = GetCurrentCompiler().AddConstant(expr.Value)
	}

	return expr
}

func (expr *IntExpression) Generate(ob *OpCodeBuf) {
	if expr.Value >= 0 && expr.Value < 256 {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT_1BYTE, expr.Value)
	} else if expr.Value >= 0 && expr.Value < 65536 {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT_2BYTE, expr.Value)
	} else {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT, expr.Index)
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
	expr.SetType(NewType(BasicTypeFloat))
	expr.GetType().Fix()

	if expr.Value != 0.0 && expr.Value != 1.0 {
		expr.Index = GetCurrentCompiler().AddConstant(expr.Value)
	}
	return expr
}

func (expr *FloatExpression) Generate(ob *OpCodeBuf) {
	if expr.Value == 0.0 {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_FLOAT_0)
	} else if expr.Value == 1.0 {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_FLOAT_1)
	} else {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_FLOAT, expr.Index)
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
	Index int
}

func (expr *StringExpression) Fix() Expression {
	expr.SetType(NewType(BasicTypeString))
	expr.GetType().Fix()

	expr.Index = GetCurrentCompiler().AddConstant(expr.Value)
	return expr
}

func (expr *StringExpression) Generate(ob *OpCodeBuf) {
	ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_STRING, expr.Index)
}

func CreateStringExpression(pos Position, value string) *StringExpression {
	expr := &StringExpression{
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

//
// InterfaceExpression
//
type InterfaceExpression struct {
	ExpressionBase
	Data Expression
}

func (expr *InterfaceExpression) Fix() Expression {
	expr.SetType(NewType(BasicTypeInterface))
	expr.Data.Fix()
	expr.GetType().Fix()

	return expr
}

func (expr *InterfaceExpression) Generate(ob *OpCodeBuf) {
	expr.Data.Generate(ob)
	ob.GenerateCode(expr.Position(), vm.OP_CODE_NEW_INTERFACE)
}

func CreateInterfaceExpression(pos Position) *InterfaceExpression {
	expr := &InterfaceExpression{
		Data: CreateNilExpression(pos),
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
	expr.SetType(NewType(BasicTypeNil))
	expr.GetType().Fix()

	return expr
}

func (expr *NilExpression) Generate(ob *OpCodeBuf) {
	ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_NIL)
}

func CreateNilExpression(pos Position) *NilExpression {
	expr := &NilExpression{}
	expr.SetPosition(pos)

	return expr
}

func isNilExpression(expr Expression) bool {
	_, ok := expr.(*NilExpression)
	return ok
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

	for i := 0; i < len(expr.List); i++ {
		// TODO: 直接使用value
		keyValueExpr, ok := expr.List[i].(*KeyValueExpression)
		if ok {
			expr.List[i] = keyValueExpr.Value
		}

		expr.List[i] = expr.List[i].Fix()
		expr.List[i] = CreateAssignCast(expr.List[i], elemType)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *ArrayExpression) Generate(ob *OpCodeBuf) {
	for i := len(expr.List) - 1; i >= 0; i-- {
		expr.List[i].Generate(ob)
	}

	ob.GenerateCode(expr.Position(), vm.OP_CODE_NEW_ARRAY, len(expr.List))
}

func CreateArrayExpression(typ *Type, exprList []Expression) *ArrayExpression {
	expr := &ArrayExpression{
		List: exprList,
	}
	expr.SetType(typ)

	return expr
}

//
// MapExpression
//
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

	for i := 0; i < len(expr.KeyList); i++ {
		expr.KeyList[i] = expr.KeyList[i].Fix()
		expr.KeyList[i] = CreateAssignCast(expr.KeyList[i], keyType)
	}

	for i := 0; i < len(expr.ValueList); i++ {
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

	ob.GenerateCode(expr.Position(), vm.OP_CDOE_NEW_MAP, size)
}

func CreateMapExpression(typ *Type, valueList []Expression) *MapExpression {
	expr := &MapExpression{
		KeyList:   make([]Expression, len(valueList)),
		ValueList: make([]Expression, len(valueList)),
	}
	expr.SetType(typ)

	for i, subExpr := range valueList {
		keyValueExpr, ok := subExpr.(*KeyValueExpression)
		if !ok {
			panic("TODO")
		}
		expr.KeyList[i] = keyValueExpr.Key
		expr.ValueList[i] = keyValueExpr.Value
	}

	return expr
}

//
// StructExpression
//
type StructExpression struct {
	ExpressionBase
	FieldList []Expression
}

func (expr *StructExpression) Fix() Expression {
	for i, field := range expr.FieldList {
		fieldType := expr.Type.structType.Fields[i].Type

		if field == nil {
			// TODO: default value
			// expr.FieldList[i] = field
			panic("TODO")
		}
		expr.FieldList[i] = expr.FieldList[i].Fix()
		expr.FieldList[i] = CreateAssignCast(expr.FieldList[i], fieldType)
	}

	expr.GetType().Fix()

	return expr
}

func (expr *StructExpression) Generate(ob *OpCodeBuf) {
	for i := len(expr.FieldList) - 1; i >= 0; i-- {
		expr.FieldList[i].Generate(ob)
	}

	count := len(expr.FieldList)
	ob.GenerateCode(expr.Position(), vm.OP_CODE_NEW_STRUCT, count)
}

func CreateStructExpression(typ *Type, valueList []Expression) *StructExpression {
	fieldCount := len(typ.structType.Fields)

	expr := &StructExpression{
		FieldList: make([]Expression, fieldCount),
	}
	expr.SetType(typ)

	fieldMap := make(map[string]*StructField)
	for i, field := range typ.structType.Fields {
		field.Index = i
		fieldMap[field.Name] = field
	}

	for _, subExpr := range valueList {
		keyValueExpr, ok := subExpr.(*KeyValueExpression)
		if !ok {
			panic("TODO")
		}

		fieldName := keyValueExpr.Key.(*IdentifierExpression).Name
		field, ok := fieldMap[fieldName]
		if !ok {
			panic("TODO")
		}

		expr.FieldList[field.Index] = keyValueExpr.Value
	}

	return expr
}

//
// IdentifierExpression 变量表达式
//
type IdentifierExpression struct {
	ExpressionBase
	PackageName string
	Name        string
	Obj         interface{} // 变量,函数,包,自定义类型(FunctionIdentifier Declaration Package)
	Block       *Block
}

type FunctionIdentifier struct {
	Func  *FunctionDefinition
	Index int
}

// TODO: 变量查找顺序,先查本地变量/函数,再搜全局,再搜内置
func (expr *IdentifierExpression) Fix() Expression {
	c := GetCurrentCompiler()
	pkg := GetCurrentPackage()

	//
	// 判断是否是变量(本地/全局)
	//
	declaration := expr.Block.SearchDeclaration(expr.Name)
	if declaration == nil {
		declaration = c.SearchDeclaration(expr.PackageName, expr.Name)
	}

	if declaration != nil {
		expr.SetType(declaration.Type.Copy())
		expr.Obj = declaration
		expr.GetType().Fix()
		return expr
	}

	//
	// 判断是否是函数
	//
	fd, index := c.SearchFunction(expr.PackageName, expr.Name)
	if fd != nil {
		expr.SetType(fd.CopyType())
		expr.Obj = &FunctionIdentifier{
			Func:  fd,
			Index: index,
		}
		expr.GetType().Fix()

		return expr
	}

	//
	// 判断是否是包引用
	//
	for _, imp := range pkg.importList {
		if imp.packageName == expr.Name {
			expr.SetType(NewType(BasicTypePackage))
			expr.Obj = imp.packageName
			expr.GetType().Fix()
			return expr
		}
	}

	compileError(expr.Position(), IDENTIFIER_NOT_FOUND_ERR, expr.Name)

	return nil
}

func (expr *IdentifierExpression) Generate(ob *OpCodeBuf) {
	switch inner := expr.Obj.(type) {
	case *FunctionIdentifier:
		ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_FUNCTION, inner.Index)

	case *Declaration:
		var code byte

		if inner.IsLocal {
			code = vm.OP_CODE_PUSH_STACK
		} else {
			code = vm.OP_CODE_PUSH_STATIC
		}
		ob.GenerateCode(expr.Position(), code, inner.Index)
	}
}

func CreateIdentifierExpression(pos Position, name string) *IdentifierExpression {
	expr := &IdentifierExpression{Name: name}
	expr.SetPosition(pos)

	c := GetCurrentPackage()
	if c != nil {
		expr.Block = c.currentBlock
		expr.PackageName = c.GetPackageName()
	}

	return expr
}

//
// CallExpression 函数调用表达式
//
type CallExpression struct {
	ExpressionBase
	Func Expression   // 函数名
	Args []Expression // 实参列表
}

func (expr *CallExpression) Fix() Expression {
	expr.Func = expr.Func.Fix()

	funcExpr, ok := expr.Func.(*IdentifierExpression)
	if !ok {
		panic("TODO")
		// compileError(expr.Position(), FUNCTION_NOT_FOUND_ERR, "")
	}

	fd := funcExpr.Obj.(*FunctionIdentifier).Func

	expr.Args = FixArgList(fd, expr.Args)

	FixReturn(fd, expr.Type)

	expr.GetType().Fix()

	return expr
}

func FixArgList(fd *FunctionDefinition, argumentList []Expression) []Expression {
	parameterList := fd.Type.funcType.Params

	paramLen := len(parameterList)

	if paramLen > 0 {
		lastP := parameterList[paramLen-1]
		if lastP.Ellipsis {
			newArgList := make([]Expression, 0)
			for _, expr := range argumentList[:paramLen-1] {
				newArgList = append(newArgList, expr)
			}
			lastArg := CreateArrayExpression(lastP.Type, argumentList[paramLen-1:])
			newArgList = append(newArgList, lastArg)

			argumentList = newArgList
		}
	}

	argLen := len(argumentList)

	if argLen != paramLen {
		compileError(fd.GetType().Position(), ARGUMENT_COUNT_MISMATCH_ERR, paramLen, argLen)
	}

	for i := 0; i < paramLen; i++ {
		argumentList[i] = argumentList[i].Fix()
		if !argumentList[i].GetType().Equal(parameterList[i].Type) {
			compileError(
				argumentList[i].Position(),
				ARGUMENT_COUNT_MISMATCH_ERR,
				parameterList[i].Name,
				parameterList[i].Type.GetBasicType(),
				argumentList[i].GetType().GetBasicType(),
			)
		}
	}

	return argumentList
}

// 设置返回值类型
func FixReturn(fd *FunctionDefinition, typ *Type) {
	resultCount := len(fd.Type.funcType.Results)

	if resultCount == 0 {
		typ.SetBasicType(BasicTypeVoid)
	} else if resultCount == 1 {
		typ.SetBasicType(fd.Type.funcType.Results[0].Type.GetBasicType())
	} else {
		typeList := make([]*Type, resultCount)
		for i, resultType := range fd.Type.funcType.Results {
			typeList[i] = resultType.Type.Copy()
		}
		typ.SetBasicType(BasicTypeMultipleValues)
		typ.multipleValueType = NewMultipleValueType(typeList)
	}
}

// CallExpression
// returns
// args
// return address
// locals -- base
func (expr *CallExpression) Generate(ob *OpCodeBuf) {
	fd := expr.Func.(*IdentifierExpression).Obj.(*FunctionIdentifier).Func

	for _, p := range fd.Type.funcType.Results {
		GetTypeDefaultValue(p.Type, expr.Position()).Generate(ob)
	}

	for _, param := range expr.Args {
		param.Generate(ob)
	}

	expr.Func.Generate(ob)
	ob.GenerateCode(expr.Position(), vm.OP_CODE_INVOKE)
}

func NewFunctionCallExpression(pos Position, function Expression, argumentList []Expression) *CallExpression {
	expr := &CallExpression{
		Func: function,
		Args: argumentList,
	}
	expr.SetType(NewType(BasicTypeVoid))
	expr.SetPosition(pos)

	return expr
}

//
// SelectorExpression
//
type SelectorExpression struct {
	ExpressionBase
	X     Expression // 实例
	Sel   string     // 成员名称
	Index int
}

func (expr *SelectorExpression) Fix() Expression {
	var newExpr Expression

	expr.X = expr.X.Fix()
	typ := expr.X.GetType()

	// TODO: IsType()
	switch {
	case typ.IsPackage():
		newExpr = FixPackageSelectorExpression(expr)
	case typ.IsStruct():
		newExpr = FixStructSelectorExpression(expr)
	default:
		compileError(expr.Position(), MEMBER_EXPRESSION_TYPE_ERR)
	}

	newExpr.GetType().Fix()

	return newExpr
}

func (expr *SelectorExpression) Generate(ob *OpCodeBuf) {
	expr.X.Generate(ob)
	// TODO: 临时用下OP_CODE_PUSH_INT_2BYTE
	ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT_2BYTE, expr.Index)
	ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_STRUCT)
}

func FixPackageSelectorExpression(expr *SelectorExpression) Expression {
	packageName := expr.X.(*IdentifierExpression).Obj.(string)

	fd, index := GetCurrentCompiler().SearchFunction(packageName, expr.Sel)
	if fd != nil {
		newExpr := CreateIdentifierExpression(expr.Position(), expr.Sel)
		newExpr.Obj = &FunctionIdentifier{
			Func:  fd,
			Index: index,
		}

		newExpr.SetType(fd.CopyType())
		newExpr.GetType().Fix()

		return newExpr
	}

	decl := GetCurrentCompiler().SearchDeclaration(packageName, expr.Sel)
	if decl != nil {
		newExpr := CreateIdentifierExpression(expr.Position(), expr.Sel)
		newExpr.SetType(decl.Type.Copy())
		newExpr.Obj = decl
		newExpr.GetType().Fix()

		return newExpr
	}

	panic(fmt.Sprintf("package filed not found '%s'", expr.Sel))
}

func FixStructSelectorExpression(expr *SelectorExpression) Expression {
	// innerExpr := expr.X
	// innerExpr.GetType().Fix()
	expr.X = expr.X.Fix()

	for i, field := range expr.X.GetType().structType.Fields {
		if field.Name == expr.Sel {
			expr.Index = i
			expr.SetType(field.Type.Copy())

			return expr
		}
	}

	panic("TODO")
}

func CreateSelectorExpression(expression Expression, memberName string) *SelectorExpression {
	expr := &SelectorExpression{
		X:   expression,
		Sel: memberName,
	}
	expr.SetPosition(expression.Position())

	return expr
}

//
// IndexExpression
//
type IndexExpression struct {
	ExpressionBase
	X     Expression
	Index Expression
}

func (expr *IndexExpression) Fix() Expression {

	expr.X = expr.X.Fix()
	expr.Index = expr.Index.Fix()

	if !expr.X.GetType().IsArray() && !expr.X.GetType().IsMap() {
		compileError(expr.Position(), INDEX_LEFT_OPERAND_NOT_ARRAY_ERR)
	}

	if expr.X.GetType().IsArray() {
		expr.SetType(expr.X.GetType().arrayType.ElementType.Copy())

		if !expr.Index.GetType().IsInt() {
			compileError(expr.Position(), INDEX_NOT_INT_ERR)
		}
	} else if expr.X.GetType().IsMap() {
		expr.SetType(expr.X.GetType().mapType.Value.Copy())
	}

	expr.GetType().Fix()

	return expr
}

func (expr *IndexExpression) Generate(ob *OpCodeBuf) {
	expr.X.Generate(ob)
	expr.Index.Generate(ob)

	switch {
	case expr.X.GetType().IsArray():
		code := vm.OP_CODE_PUSH_ARRAY
		ob.GenerateCode(expr.Position(), code)
	case expr.X.GetType().IsMap():
		code := vm.OP_CODE_PUSH_MAP
		ob.GenerateCode(expr.Position(), code)
	default:
		panic("TODO")
	}
}

func CreateIndexExpression(pos Position, array, index Expression) *IndexExpression {
	expr := &IndexExpression{
		X:     array,
		Index: index,
	}
	expr.SetPosition(pos)

	return expr
}

type KeyValueExpression struct {
	ExpressionBase
	Key   Expression
	Value Expression
}

func (expr *KeyValueExpression) Fix() Expression {
	if expr.Key != nil {
		expr.Key = expr.Key.Fix()
	}
	if expr.Value != nil {
		expr.Value = expr.Value.Fix()
	}

	expr.GetType().Fix()

	return expr
}

func (expr *KeyValueExpression) Generate(ob *OpCodeBuf) {
	if expr.Value != nil {
		expr.Value.Generate(ob)
	}

	if expr.Key != nil {
		expr.Key.Generate(ob)
	}
}

func CreateKeyValueExpression(pos Position, key, value Expression) *KeyValueExpression {
	expr := &KeyValueExpression{
		Key:   key,
		Value: value,
	}
	expr.SetPosition(pos)

	return expr
}

func CreateCompositeLit(typ *Type, valueList []Expression) Expression {
	if typ.IsArray() {
		return CreateArrayExpression(typ, valueList)
	}

	if typ.IsMap() {
		return CreateMapExpression(typ, valueList)
	}

	if typ.IsStruct() {
		return CreateStructExpression(typ, valueList)
	}

	return nil
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

		ob.GenerateCode(expr.Position(), code+offset)

	case LogicalAndOperator, LogicalOrOperator:
		var jumpCode, logicalCode byte

		if operator == LogicalAndOperator {
			jumpCode = vm.OP_CODE_JUMP_IF_FALSE
			logicalCode = vm.OP_CODE_LOGICAL_AND
		} else {
			jumpCode = vm.OP_CODE_JUMP_IF_TRUE
			logicalCode = vm.OP_CODE_LOGICAL_OR
		}

		label := ob.GetLabel()

		expr.left.Generate(ob)
		ob.GenerateCode(expr.Position(), vm.OP_CODE_DUPLICATE)
		ob.GenerateCode(expr.Position(), jumpCode, label)

		expr.right.Generate(ob)

		ob.GenerateCode(expr.Position(), logicalCode)

		ob.SetLabel(label)
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
	code := vm.OP_CODE_MINUS_INT
	if expr.GetType().IsFloat() {
		code = vm.OP_CODE_MINUS_FLOAT
	}
	ob.GenerateCode(expr.Position(), code)
}

func (expr *UnaryExpression) GenerateNot(ob *OpCodeBuf) {
	expr.Value.Generate(ob)
	ob.GenerateCode(expr.Position(), vm.OP_CODE_LOGICAL_NOT)
}

func NewUnaryExpression(pos Position, operator UnaryOperatorKind, value Expression) *UnaryExpression {
	expr := &UnaryExpression{
		Operator: operator,
		Value:    value,
	}
	expr.SetPosition(pos)

	return expr
}
