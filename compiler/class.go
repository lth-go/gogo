package compiler

import (
	"strings"
)

const DEFAULT_CONSTRUCTOR_NAME = "init"

//
// ClassDefinition
//
type ClassDefinition struct {
	PosImpl

	packageNameList []string
	name            string

	extends         []*Extend
	superClass *ClassDefinition

	memberList []MemberDeclaration
}

func (cd *ClassDefinition) add_class() {
	compiler := getCurrentCompiler()

	src_package_name := cd.getPackageName()

	for i, vmClass := range compiler.vmClassList {
        if (src_package_name == vmClass.PackageName) && (cd.name == vmClass.Name) {
            return i
        }
    }

	ret := len(compiler.vmClassList)

	dest := &VmClass{}
	compiler.vmClassList = append(compiler.vmClassList, dest)

    dest.PackageName = src_package_name
    dest.Name = cd.name
    dest.IsImplemented = false

	for _, sup_pos := range cd.extends {
        var dummy int
        search_class_and_add(cd.Position, sup_pos.identifier, &dummy)
    }

    return ret
}

func (cd *ClassDefinition) getPackageName() string {
	return strings.Join(cd.packageNameList, ".")
}

func startClassdefinition(identifier string, extends []*Extend, pos Position) {
	compiler := getCurrentCompiler()

	cd := &ClassDefinition{}

	cd.packageNameList = compiler.packageNameList
	cd.name = identifier
	cd.extends = extends

	cd.SetPosition(pos)

	if compiler.currentClassDefinition != nil {
		panic("TODO")
	}

	compiler.currentClassDefinition = cd
}

func classDefine(member_list []MemberDeclaration) {
	compiler := getCurrentCompiler()

	cd := compiler.currentClassDefinition

	if cd == nil {
		panic("TODO")
	}

	if compiler.classDefinitionList == nil {
		compiler.classDefinitionList = []*ClassDefinition{}
	}
	compiler.classDefinitionList = append(compiler.classDefinitionList, cd)

	cd.member = member_list
	compiler.currentClassDefinition = nil
}

func methodFunctionDefine(typ *TypeSpecifier, identifier string, parameter_list []*Parameter, block *Block) *FunctionDefinition {

	fd := createFunctionDefinition(typ, identifier, parameter_list, block)

	return fd
}


// ==============================
// Extend
// ==============================
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

func chainExtendList(list []*Extend, add string) {
	newExtend := createExtendList(add)

	list = append(list, newExtend)

	return list
}

// ==============================
// MemberDeclaration
// ==============================
type MemberDeclaration interface{}

func chainMemberDeclaration(list []MemberDeclaration, add MemberDeclaration) []MemberDeclaration {
	list = append(list, add)

	return list
}

type MethodMember struct {
	PosImpl

	functionDefinition *FunctionDefinition
	methodIndex        int
}

func createMethodMember(function_definition *FunctionDefinition, pos Position) []MemberDeclaration {

	ret := MethodMember{}
	ret.SetPosition(pos)

	ret.functionDefinition = function_definition

	if function_definition.block == nil {
		compileError(pos, CONCRETE_METHOD_HAS_NO_BODY_ERR)
	}

	function_definition.class_definition = compiler.currentClassDefinition

	return ret
}

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

