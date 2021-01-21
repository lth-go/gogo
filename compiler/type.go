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

// Type 表达式类型
type Type struct {
	PosBase
	basicType vm.BasicType
	sliceType *ArrayType
	funcType  *FuncType
	mapType   *MapType
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

	if !t.sliceType.Equal(t2.sliceType) {
		return false
	}

	if !t.funcType.Equal(t2.funcType) {
		return false
	}

	return true
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
	newType := CreateType(vm.BasicTypeSlice, pos)
	newType.sliceType = NewArrayType(typ)
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

func (t *Type) IsArray() bool {
	return t.GetBasicType() == vm.BasicTypeSlice
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

	newType.sliceType = t.sliceType.Copy()
	newType.funcType = t.funcType.Copy()

	return newType
}
