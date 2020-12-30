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

//
// TypeSpecifier
//
type classRef struct {
	identifier      string
	classDefinition *ClassDefinition
	classIndex      int
}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl

	// 基本类型
	basicType vm.BasicType

	// 类引用
	classRef classRef

	// 派生类型
	deriveList []TypeDerive
}

func (t *TypeSpecifier) fix() {

	for _, deriveIfs := range t.deriveList {
		derive, ok := deriveIfs.(*FunctionDerive)
		if ok {
			for _, parameter := range derive.parameterList {
				parameter.typeSpecifier.fix()
			}
		}
	}

	if t.basicType == vm.ClassType && t.classRef.classDefinition == nil {

		cd := searchClass(t.classRef.identifier)
		if cd == nil {
			compileError(t.Position(), TYPE_NAME_NOT_FOUND_ERR, t.classRef.identifier)
			return
		}

		t.classRef.classDefinition = cd
		t.classRef.classIndex = cd.addToCurrentCompiler()
		return
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
func createClassTypeSpecifier(identifier string, pos Position) *TypeSpecifier {

	typ := &TypeSpecifier{
		basicType: vm.ClassType,
		classRef: classRef{
			identifier: identifier,
			classDefinition: nil,
		},
	}
	typ.SetPosition(pos)

	return typ
}

func createArrayTypeSpecifier(typ *TypeSpecifier) *TypeSpecifier {
	typ.appendDerive(&ArrayDerive{})
	return typ
}

func (t *TypeSpecifier) appendDerive(derive TypeDerive) {
	if t.deriveList == nil {
		t.deriveList = []TypeDerive{}
	}
	t.deriveList = append(t.deriveList, derive)
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
