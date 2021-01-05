package compiler

import (
	"github.com/lth-go/gogogogo/vm"
)

//
// derive
//

type TypeDerive interface{}

type FunctionDerive struct {
	parameterList []*Parameter
}

type ArrayDerive struct{}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl
	// 基本类型
	basicType vm.BasicType
	// 派生类型
	deriveType TypeDerive
}

func (t *TypeSpecifier) fix() {
	derive, ok := t.deriveType.(*FunctionDerive)
	if ok {
		for _, parameter := range derive.parameterList {
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
	typ.deriveType = &ArrayDerive{}
	return typ
}

func (t *TypeSpecifier) isArrayDerive() bool {
	return isArray(t)
}

func (t *TypeSpecifier) isModule() bool {
	return isModule(t)
}

// utils
func cloneTypeSpecifier(src *TypeSpecifier) *TypeSpecifier {
	typ := &TypeSpecifier{}

	*typ = *src

	return typ
}
