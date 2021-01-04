package compiler

import (
	"strings"

	"github.com/lth-go/gogogogo/vm"
)

const defaultConstructorName = "init"

//
// ClassDefinition
//
type ClassDefinition struct {
	PosImpl
	packageNameList []string
	name            string
	memberList      []MemberDeclaration
}

func (cd *ClassDefinition) getPackageName() string {
	return strings.Join(cd.packageNameList, ".")
}

// 添加类到当前compiler
func (cd *ClassDefinition) addToCurrentCompiler() int {
	compiler := getCurrentCompiler()

	srcPackageName := cd.getPackageName()

	for i, vmClass := range compiler.vmClassList {
		if (srcPackageName == vmClass.PackageName) && (cd.name == vmClass.Name) {
			return i
		}
	}

	ret := len(compiler.vmClassList)

	dest := &vm.Class{
		PackageName:   srcPackageName,
		Name:          cd.name,
		IsImplemented: false,
	}

	compiler.vmClassList = append(compiler.vmClassList, dest)

	return ret
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

	return nil
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

func createMethodMember(functionDefinition *FunctionDefinition, pos Position) []MemberDeclaration {
	compiler := getCurrentCompiler()

	ret := &MethodMember{}
	ret.SetPosition(pos)

	ret.functionDefinition = functionDefinition

	if functionDefinition.block == nil {
		compileError(pos, CONCRETE_METHOD_HAS_NO_BODY_ERR)
	}

	functionDefinition.classDefinition = compiler.currentClassDefinition

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
