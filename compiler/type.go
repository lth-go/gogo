package compiler

import (
	"fmt"
	"strings"
)

type BasicType int

const (
	BasicTypeNoType BasicType = iota - 1
	BasicTypeBool
	BasicTypeInt
	BasicTypeFloat
	BasicTypeString
	BasicTypeNil
	BasicTypeVoid
	BasicTypePackage
	BasicTypeArray
	BasicTypeMap
	BasicTypeStruct
	BasicTypeFunc
	BasicTypeMultipleValues
	BasicTypeInterface
	BasicTypePointer
)

//
// Type 表达式类型
//
type Type struct {
	PosBase
	basicType         BasicType
	arrayType         *ArrayType
	funcType          *FuncType
	mapType           *MapType
	multipleValueType *MultipleValueType // 用于处理函数多返回值
	structType        *StructType
}

func (t *Type) Fix() {
	// TODO: 修正引用类型别名
	// if t.funcType != nil {
	//     for _, param := range t.funcType.Params {
	//         param.Type.Fix()
	//     }
	// }
}

func (t *Type) GetBasicType() BasicType {
	return t.basicType
}

func (t *Type) SetBasicType(basicType BasicType) {
	t.basicType = basicType
}

func (t *Type) Equal(t2 *Type) bool {
	if t.IsInterface() || t2.IsInterface() {
		return true
	}

	if t.GetBasicType() != t2.GetBasicType() {
		return false
	}

	if !t.arrayType.Equal(t2.arrayType) {
		return false
	}

	if !t.funcType.Equal(t2.funcType) {
		return false
	}

	if !t.multipleValueType.Equal(t2.multipleValueType) {
		return false
	}

	return true
}

func (t *Type) GetResultCount() int {
	// TODO: 无返回值需要返回1
	if t.IsMultipleValues() {
		return len(t.multipleValueType.List)
	} else if t.IsVoid() {
		return 0
	} else {
		return 1
	}
}

func NewType(basicType BasicType) *Type {
	return &Type{
		basicType: basicType,
	}
}

//
// 复合类型
//
type ArrayType struct {
	Len         int64
	ElementType *Type
}

func NewArrayType(elementType *Type) *ArrayType {
	return &ArrayType{
		ElementType: elementType,
	}
}

func (t *ArrayType) Copy() *ArrayType {
	if t == nil {
		return nil
	}

	return NewArrayType(t.ElementType.Copy())
}

func (t *ArrayType) Equal(t2 *ArrayType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}

	return t.ElementType.Equal(t2.ElementType)
}

//
// FuncType
//
type FuncType struct {
	Params  []*Parameter
	Results []*Parameter
}

func NewFuncType(params []*Parameter, results []*Parameter) *FuncType {
	return &FuncType{
		Params:  params,
		Results: results,
	}
}

func (t *FuncType) Copy() *FuncType {
	if t == nil {
		return nil
	}

	copyParams := func(params []*Parameter) []*Parameter {
		newParams := []*Parameter{}

		for _, p := range params {
			newParams = append(newParams, &Parameter{
				Type: p.Type.Copy(),
				Name: p.Name,
			})
		}
		return newParams
	}

	return NewFuncType(copyParams(t.Params), copyParams(t.Results))
}

func (t *FuncType) Equal(t2 *FuncType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}

	if len(t.Params) != len(t2.Params) {
		return false
	}

	if len(t.Results) != len(t2.Results) {
		return false
	}

	for i := 0; i < len(t.Params); i++ {
		if !t.Params[i].Type.Equal(t2.Params[i].Type) {
			return false
		}
	}

	for i := 0; i < len(t.Results); i++ {
		if !t.Results[i].Type.Equal(t2.Results[i].Type) {
			return false
		}
	}

	return true
}

type MapType struct {
	Key   *Type
	Value *Type
}

func NewMapType(keyType, valueType *Type) *MapType {
	return &MapType{
		Key:   keyType,
		Value: valueType,
	}
}

func (t *MapType) Copy() *MapType {
	if t == nil {
		return nil
	}

	return &MapType{
		Key:   t.Key.Copy(),
		Value: t.Value.Copy(),
	}
}

func (t *MapType) Equal(t2 *MapType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}

	if !t.Key.Equal(t2.Key) {
		return false
	}

	if !t.Value.Equal(t2.Value) {
		return false
	}

	return true
}

type MultipleValueType struct {
	List []*Type
}

func NewMultipleValueType(list []*Type) *MultipleValueType {
	return &MultipleValueType{
		List: list,
	}
}

func (t *MultipleValueType) Copy() *MultipleValueType {
	if t == nil {
		return nil
	}

	list := make([]*Type, len(t.List))

	for i, subType := range t.List {
		list[i] = subType.Copy()
	}

	return &MultipleValueType{
		List: list,
	}
}

func (t *MultipleValueType) Equal(t2 *MultipleValueType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}
	if len(t.List) != len(t2.List) {
		return false
	}

	for i := 0; i < len(t.List); i++ {
		if !t.List[i].Equal(t2.List[i]) {
			return false
		}
	}

	return true
}

type StructType struct {
	Fields []*StructField
}

type StructField struct {
	Name  string
	Type  *Type
	Index int
}

func NewStructType(fields []*StructField) *StructType {
	return &StructType{
		Fields: fields,
	}
}

func (t *StructType) Copy() *StructType {
	if t == nil {
		return nil
	}

	fieldList := make([]*StructField, 0)

	for _, field := range t.Fields {
		fieldList = append(fieldList, &StructField{
			Name: field.Name,
			Type: field.Type.Copy(),
		})
	}

	return &StructType{
		Fields: fieldList,
	}
}

func (t *StructType) Equal(t2 *StructType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}

	if len(t.Fields) != len(t2.Fields) {
		return false
	}

	for i := 0; i < len(t.Fields); i++ {
		f1 := t.Fields[i]
		f2 := t.Fields[i]

		if f1.Name != f2.Name {
			return false
		}

		if !f1.Type.Equal(f2.Type) {
			return false
		}
	}

	return true
}

//
// create
//
func CreateType(basicType BasicType, pos Position) *Type {
	typ := NewType(basicType)
	typ.SetPosition(pos)
	return typ
}

func CreateArrayType(typ *Type, pos Position) *Type {
	newType := CreateType(BasicTypeArray, pos)
	newType.arrayType = NewArrayType(typ)
	return newType
}

func CreateFuncType(params []*Parameter, results []*Parameter) *Type {
	newType := NewType(BasicTypeFunc)
	newType.funcType = NewFuncType(params, results)
	return newType
}

func CreateMapType(keyType *Type, valueType *Type, pos Position) *Type {
	newType := CreateType(BasicTypeMap, pos)
	newType.mapType = NewMapType(keyType, valueType)
	return newType
}

func CreateInterfaceType(pos Position) *Type {
	newType := CreateType(BasicTypeInterface, pos)
	return newType
}

func CreateStructType(pos Position, fieldDeclList []*StructField) *Type {
	newType := CreateType(BasicTypeStruct, pos)
	newType.structType = NewStructType(fieldDeclList)
	return newType
}

func CreateFieldDecl(name string, fieldType *Type) *StructField {
	return &StructField{
		Name: name,
		Type: fieldType,
	}
}

func (t *Type) IsArray() bool {
	return t.GetBasicType() == BasicTypeArray
}

func (t *Type) IsFunc() bool {
	return t.GetBasicType() == BasicTypeFunc
}

func (t *Type) IsComposite() bool {
	return t.IsArray() || t.IsFunc()
}

func (t *Type) IsVoid() bool {
	return t.GetBasicType() == BasicTypeVoid
}

func (t *Type) IsBool() bool {
	return t.GetBasicType() == BasicTypeBool
}

func (t *Type) IsInt() bool {
	return t.GetBasicType() == BasicTypeInt
}

func (t *Type) IsFloat() bool {
	return t.GetBasicType() == BasicTypeFloat
}

func (t *Type) IsString() bool {
	return t.GetBasicType() == BasicTypeString
}

func (t *Type) IsPackage() bool {
	return t.GetBasicType() == BasicTypePackage
}

func (t *Type) IsObject() bool {
	return t.IsString() || t.IsArray()
}

func (t *Type) IsNil() bool {
	return t.GetBasicType() == BasicTypeNil
}

func (t *Type) IsMultipleValues() bool {
	return t.GetBasicType() == BasicTypeMultipleValues
}

func (t *Type) IsMap() bool {
	return t.GetBasicType() == BasicTypeMap
}

func (t *Type) IsInterface() bool {
	return t.GetBasicType() == BasicTypeInterface
}

func (t *Type) IsStruct() bool {
	return t.GetBasicType() == BasicTypeStruct
}

func (t *Type) GetTypeName() string {
	typeName := GetBasicTypeName(t.GetBasicType())

	switch {
	case t.IsArray():
		typeName = "[]" + typeName
	case t.IsFunc():
		paramTypeNameList := []string{}
		resultTypeNameList := []string{}

		for _, p := range t.funcType.Params {
			paramTypeNameList = append(paramTypeNameList, p.Type.GetTypeName())
		}

		for _, p := range t.funcType.Results {
			resultTypeNameList = append(resultTypeNameList, p.Type.GetTypeName())
		}

		typeName = fmt.Sprintf(
			"func(%s) (%s)",
			strings.Join(paramTypeNameList, ", "),
			strings.Join(resultTypeNameList, ", "),
		)
	}

	return typeName
}

func GetBasicTypeName(typ BasicType) string {
	switch typ {
	case BasicTypeBool:
		return "bool"
	case BasicTypeInt:
		return "int"
	case BasicTypeFloat:
		return "float"
	case BasicTypeString:
		return "string"
	case BasicTypeNil:
		return "nil"
	case BasicTypeFunc:
		return "func"
	default:
		panic(fmt.Sprintf("bad case. type..%d\n", typ))
	}
}

// 根据字面量创建基本类型
func CreateTypeByName(name string, pos Position) *Type {
	basicType := BasicTypeNoType

	// TODO:
	basicTypeMap := map[string]BasicType{
		"bool":   BasicTypeBool,
		"int":    BasicTypeInt,
		"float":  BasicTypeFloat,
		"string": BasicTypeString,
	}

	_, ok := basicTypeMap[name]
	if ok {
		basicType = basicTypeMap[name]
	}

	return CreateType(basicType, pos)
}

func (t *Type) Copy() *Type {
	newType := NewType(t.GetBasicType())

	newType.arrayType = t.arrayType.Copy()
	newType.funcType = t.funcType.Copy()
	newType.multipleValueType = t.multipleValueType.Copy()
	newType.mapType = t.mapType.Copy()
	newType.structType = t.structType.Copy()

	return newType
}
