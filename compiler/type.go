package compiler

import (
	"fmt"
	"strings"

	"github.com/lth-go/gogo/vm"
)

//
// composite type
//
type SliceType struct {
	Len         int64
	ElementType *TypeSpecifier
}

func NewSliceType(elementType *TypeSpecifier) *SliceType {
	return &SliceType{
		ElementType: elementType,
	}
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

// TypeSpecifier 表达式类型
type TypeSpecifier struct {
	PosImpl
	name      string
	basicType vm.BasicType
	sliceType *SliceType
	funcType  *FuncType
}

func (t *TypeSpecifier) fix() {
	if t.funcType != nil {
		for _, parameter := range t.funcType.Params {
			parameter.typeSpecifier.fix()
		}
	}
}

func (t *TypeSpecifier) GetBasicType() vm.BasicType {
	return t.basicType
}

// TODO: 临时使用
func newTypeSpecifier(basicType vm.BasicType) *TypeSpecifier {
	return &TypeSpecifier{
		basicType: basicType,
	}
}

//
// create
//
func createTypeSpecifier(basicType vm.BasicType, pos Position) *TypeSpecifier {
	typ := newTypeSpecifier(basicType)
	typ.SetPosition(pos)
	return typ
}

func createArrayTypeSpecifier(typ *TypeSpecifier) *TypeSpecifier {
	// TODO: 基本类型应该是slice
	newType := newTypeSpecifier(vm.BasicTypeSlice)
	newType.sliceType = NewSliceType(typ)
	return newType
}

func createFuncTypeSpecifier(params []*Parameter, results []*Parameter) *TypeSpecifier {
	newType := newTypeSpecifier(vm.BasicTypeFunc)
	newType.funcType = NewFuncType(params, results)
	return newType
}

// TODO: 改名
func createFuncType(fd *FunctionDefinition) *TypeSpecifier {
	typ := CopyType(fd.typeSpecifier)
	typ.funcType = NewFuncType(fd.parameterList, nil)

	return typ
}

func (t *TypeSpecifier) IsArray() bool {
	// TODO: 根据basic判断
	return t.sliceType != nil
}

func (t *TypeSpecifier) IsFunc() bool {
	// TODO: 根据basic判断
	return t.funcType != nil
}

func (t *TypeSpecifier) IsComposite() bool {
	return t.IsArray() || t.IsFunc()
}

func (t *TypeSpecifier) IsVoid() bool {
	return t.GetBasicType() == vm.BasicTypeVoid
}

func (t *TypeSpecifier) IsBool() bool {
	return t.GetBasicType() == vm.BasicTypeBool
}

func (t *TypeSpecifier) IsInt() bool {
	return t.GetBasicType() == vm.BasicTypeInt
}

func (t *TypeSpecifier) IsFloat() bool {
	return t.GetBasicType() == vm.BasicTypeFloat
}

func (t *TypeSpecifier) IsString() bool {
	return t.GetBasicType() == vm.BasicTypeString
}

func (t *TypeSpecifier) IsModule() bool {
	return t.GetBasicType() == vm.BasicTypeModule
}

func (t *TypeSpecifier) IsObject() bool {
	return t.IsString() || t.IsArray()
}

func (t *TypeSpecifier) IsNil() bool {
	return t.GetBasicType() == vm.BasicTypeNil
}

// TODO:
func (t *TypeSpecifier) IsBase() bool {
	return t.GetBasicType() == vm.BasicTypeBase
}

func (t *TypeSpecifier) GetTypeName() string {
	typeName := getBasicTypeName(t.GetBasicType())

	switch {
	case t.IsArray():
		typeName = "[]" + typeName
	case t.IsFunc():
		paramTypeNameList := []string{}
		resultTypeNameList := []string{}

		for _, p := range t.funcType.Params {
			paramTypeNameList = append(paramTypeNameList, p.typeSpecifier.GetTypeName())
		}

		for _, p := range t.funcType.Results {
			resultTypeNameList = append(resultTypeNameList, p.typeSpecifier.GetTypeName())
		}

		typeName = fmt.Sprintf(
			"func(%s) (%s)",
			strings.Join(paramTypeNameList, ", "),
			strings.Join(resultTypeNameList, ", "),
		)
	}

	return typeName
}

func getBasicTypeName(typ vm.BasicType) string {
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

// utils
func cloneTypeSpecifier(src *TypeSpecifier) *TypeSpecifier {
	typ := &TypeSpecifier{}
	*typ = *src

	return typ
}

func createTypeSpecifierAsName(name string, pos Position) *TypeSpecifier {
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

	return createTypeSpecifier(basicType, pos)
}

func CopyType(srcType *TypeSpecifier) *TypeSpecifier {
	if srcType == nil {
		return nil
	}

	destType := newTypeSpecifier(vm.BasicTypeNoType)

	// TODO: 深拷贝
	*destType = *srcType

	return destType
}
