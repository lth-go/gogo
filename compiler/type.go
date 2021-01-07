package compiler

import (
	"fmt"

	"github.com/lth-go/gogogogo/vm"
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

func NewFuncType(params []*Parameter) *FuncType {
	return &FuncType{
		Params: params,
	}
}

// TypeSpecifier 表达式类型
type TypeSpecifier struct {
	PosImpl
	name      string
	basicType vm.BasicType // 基本类型
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
	typ := &TypeSpecifier{basicType: basicType}
	typ.SetPosition(pos)
	return typ
}

func createArrayTypeSpecifier(typ *TypeSpecifier) *TypeSpecifier {
	// TODO: 基本类型应该是slice
	newType := newTypeSpecifier(typ.basicType)
	newType.sliceType = NewSliceType(typ)
	return newType
}

func (t *TypeSpecifier) IsArray() bool {
	// TODO: 根据basic判断
	return t.sliceType != nil
}

func (t *TypeSpecifier) IsFunc() bool {
	// TODO: 根据basic判断
	return t.funcType != nil
}

func (t *TypeSpecifier) isModule() bool {
	return isModule(t)
}

func (t *TypeSpecifier) IsComposite() bool {
	return t.IsArray() || t.IsFunc()
}

func (t *TypeSpecifier) GetTypeName() string {
	typeName := getBasicTypeName(t.basicType)

	switch {
	case t.IsArray():
		typeName = "[]" + typeName
	default:
		print("=====\n", t.Position().Line)
		panic("TODO:derive_tag")
	}

	return typeName
}

func getBasicTypeName(typ vm.BasicType) string {
	switch typ {
	case vm.BooleanType:
		return "bool"
	case vm.IntType:
		return "int"
	case vm.DoubleType:
		return "float"
	case vm.StringType:
		return "string"
	case vm.NullType:
		return "null"
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
	basicType := vm.NoType

	// TODO:
	basicTypeMap := map[string]vm.BasicType{
		"void":   vm.VoidType,
		"bool":   vm.BooleanType,
		"int":    vm.IntType,
		"float":  vm.DoubleType,
		"string": vm.StringType,
	}

	_, ok := basicTypeMap[name]
	if ok {
		basicType = basicTypeMap[name]
	}

	return createTypeSpecifier(basicType, pos)
}
