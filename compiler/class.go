package compiler

import (
	"strings"
	"../vm"
)

const DEFAULT_CONSTRUCTOR_NAME = "init"

//
// ClassDefinition
//
type ClassDefinition struct {
	PosImpl

	packageNameList []string
	name            string

	extends    []*Extend
	superClass *ClassDefinition

	memberList []MemberDeclaration
}

func (cd *ClassDefinition) getPackageName() string {
	return strings.Join(cd.packageNameList, ".")
}

// 添加类到当前compiler
func (cd *ClassDefinition) addToCurrentCompiler() int {
	var dummy int

	compiler := getCurrentCompiler()

	src_package_name := cd.getPackageName()

	for i, vmClass := range compiler.vmClassList {
		if (src_package_name == vmClass.PackageName) && (cd.name == vmClass.Name) {
			return i
		}
	}

	ret := len(compiler.vmClassList)

	dest := &vm.VmClass{}
	compiler.vmClassList = append(compiler.vmClassList, dest)

	dest.PackageName = src_package_name
	dest.Name = cd.name
	dest.IsImplemented = false

	for _, sup_pos := range cd.extends {
		search_class_and_add(cd.Position(), sup_pos.identifier, &dummy)
	}

	return ret
}

func (cd *ClassDefinition) getSuperFieldMethodCount() (int, int) {
	fieldIndex := 0
	methodIndex := 0

	for superCd := cd.superClass; superCd != nil; superCd = superCd.superClass {
		for _, memberIfs := range superCd.memberList {
			switch member := memberIfs.(type) {
			case *MethodMember:
				if member.methodIndex > methodIndex {
					methodIndex = member.methodIndex
				}
			case *FieldMember:
				if member.fieldIndex > fieldIndex {
					fieldIndex = member.fieldIndex
				}
			default:
				panic("TODO")
			}
		}
	}
	return fieldIndex, methodIndex
}

func (cd *ClassDefinition) searchMemberInSuper(memberName string) MemberDeclaration {
	var member MemberDeclaration

	if cd.superClass == nil {
		return nil
	}

	member = cd.superClass.searchMember(memberName)
	if member != nil {
		return member
	}

	return nil
}

func (cd *ClassDefinition) searchMember(memberName string) MemberDeclaration {

	for _, md := range cd.memberList {
		switch member := md.(type) {
		case *MethodMember:
			if member.functionDefinition.name == memberName {
				return member
			}
		case *FieldMember:
			if member.name == memberName {
				return member
			}
		default:
			panic("TODO")
		}
	}

	// 递归查找
	if cd.superClass != nil {
		member := cd.superClass.searchMember(memberName)
		if member != nil {
			return member
		}
	}

	return nil
}

func (cd *ClassDefinition) fixExtends() {
	var dummy_class_index int

	for _, extend := range cd.extends {
		super := search_class_and_add(cd.Position(), extend.identifier, &dummy_class_index)

		extend.classDefinition = super

		if cd.superClass != nil {
			compileError(cd.Position(), MULTIPLE_INHERITANCE_ERR, super.name)
		}

		cd.superClass = super
	}
}

// ==============================
// Extend
// ==============================
// 继承
type Extend struct {
	identifier      string
	classDefinition *ClassDefinition
}

func createExtendList(identifier string) []*Extend {
	extend := &Extend{
		identifier: identifier,
	}

	return []*Extend{extend}
}

func chainExtendList(list []*Extend, add string) []*Extend {
	newExtendList := createExtendList(add)

	list = append(list, newExtendList...)

	return list
}

// ==============================
// MemberDeclaration
// ==============================
type MemberDeclaration interface{}

func chainMemberDeclaration(list []MemberDeclaration, add []MemberDeclaration) []MemberDeclaration {
	list = append(list, add...)

	return list
}

//
// MethodMember
//
type MethodMember struct {
	PosImpl

	functionDefinition *FunctionDefinition
	methodIndex        int
}

func createMethodMember(function_definition *FunctionDefinition, pos Position) []MemberDeclaration {
	compiler := getCurrentCompiler()

	ret := &MethodMember{}
	ret.SetPosition(pos)

	ret.functionDefinition = function_definition

	if function_definition.block == nil {
		compileError(pos, CONCRETE_METHOD_HAS_NO_BODY_ERR)
	}

	function_definition.classDefinition = compiler.currentClassDefinition

	return []MemberDeclaration{ret}
}

//
// FieldMember
//
type FieldMember struct {
	PosImpl

	name          string
	typeSpecifier *TypeSpecifier
	fieldIndex    int
}

func createFieldMember(typ *TypeSpecifier, name string, pos Position) []MemberDeclaration {
	ret := &FieldMember{
		name:          name,
		typeSpecifier: typ,
	}
	ret.SetPosition(pos)

	return []MemberDeclaration{ret}
}
