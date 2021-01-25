package compiler

import (
	"fmt"
	"strings"

	"github.com/lth-go/gogo/vm"
)

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

//
// PackageType
//
type PackageType struct {
}

func NewPackageType() *PackageType {
	return &PackageType{}
}

func (t *PackageType) Copy() *PackageType {
	if t == nil {
		return nil
	}

	return &PackageType{}
}

func (t *Package) Equal(t2 *PackageType) bool {
	if t == nil && t2 == nil {
		return true
	}

	if t == nil && t2 != nil {
		return false
	}

	if t != nil && t2 == nil {
		return false
	}

	return true
}

//
// Type 表达式类型
//
type Type struct {
	PosBase
	basicType         vm.BasicType
	arrayType         *ArrayType
	funcType          *FuncType
	mapType           *MapType
	multipleValueType *MultipleValueType // TODO: 用于处理函数多返回值
	packageType       *PackageType
}

func (t *Type) Fix() {
	// TODO: 修正引用类型别名
	// if t.funcType != nil {
	//     for _, param := range t.funcType.Params {
	//         param.Type.Fix()
	//     }
	// }
}

func (t *Type) GetBasicType() vm.BasicType {
	return t.basicType
}

func (t *Type) SetBasicType(basicType vm.BasicType) {
	t.basicType = basicType
}

func (t *Type) Equal(t2 *Type) bool {
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

func NewType(basicType vm.BasicType) *Type {
	return &Type{
		basicType: basicType,
	}
}

//
// create
//
func CreateType(basicType vm.BasicType, pos Position) *Type {
	typ := NewType(basicType)
	typ.SetPosition(pos)
	return typ
}

func CreateArrayType(typ *Type, pos Position) *Type {
	newType := CreateType(vm.BasicTypeArray, pos)
	newType.arrayType = NewArrayType(typ)
	return newType
}

func CreateFuncType(params []*Parameter, results []*Parameter) *Type {
	newType := NewType(vm.BasicTypeFunc)
	newType.funcType = NewFuncType(params, results)
	return newType
}

func CreateMapType(keyType *Type, valueType *Type, pos Position) *Type {
	newType := CreateType(vm.BasicTypeMap, pos)
	newType.mapType = NewMapType(keyType, valueType)
	return newType
}

func CreateInterfaceType(pos Position) *Type {
	newType := CreateType(vm.BasicTypeInterface, pos)
	return newType
}

func (t *Type) IsArray() bool {
	return t.GetBasicType() == vm.BasicTypeArray
}

func (t *Type) IsFunc() bool {
	return t.GetBasicType() == vm.BasicTypeFunc
}

func (t *Type) IsComposite() bool {
	return t.IsArray() || t.IsFunc()
}

func (t *Type) IsVoid() bool {
	return t.GetBasicType() == vm.BasicTypeVoid
}

func (t *Type) IsBool() bool {
	return t.GetBasicType() == vm.BasicTypeBool
}

func (t *Type) IsInt() bool {
	return t.GetBasicType() == vm.BasicTypeInt
}

func (t *Type) IsFloat() bool {
	return t.GetBasicType() == vm.BasicTypeFloat
}

func (t *Type) IsString() bool {
	return t.GetBasicType() == vm.BasicTypeString
}

func (t *Type) IsPackage() bool {
	return t.GetBasicType() == vm.BasicTypePackage
}

func (t *Type) IsObject() bool {
	return t.IsString() || t.IsArray()
}

func (t *Type) IsNil() bool {
	return t.GetBasicType() == vm.BasicTypeNil
}

func (t *Type) IsMultipleValues() bool {
	return t.GetBasicType() == vm.BasicTypeMultipleValues
}

func (t *Type) IsMap() bool {
	return t.GetBasicType() == vm.BasicTypeMap
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

func GetBasicTypeName(typ vm.BasicType) string {
	switch typ {
	case vm.BasicTypeBool:
		return "bool"
	case vm.BasicTypeInt:
		return "int"
	case vm.BasicTypeFloat:
		return "float"
	case vm.BasicTypeString:
		return "string"
	case vm.BasicTypeNil:
		return "null"
	case vm.BasicTypeFunc:
		return "func"
	default:
		panic(fmt.Sprintf("bad case. type..%d\n", typ))
	}
}

// 根据字面量创建基本类型
func CreateTypeByName(name string, pos Position) *Type {
	basicType := vm.BasicTypeNoType

	// TODO:
	basicTypeMap := map[string]vm.BasicType{
		"void":   vm.BasicTypeVoid,
		"bool":   vm.BasicTypeBool,
		"int":    vm.BasicTypeInt,
		"float":  vm.BasicTypeFloat,
		"string": vm.BasicTypeString,
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

	return newType
}
