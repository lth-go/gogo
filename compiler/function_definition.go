package compiler

import (
	"strings"

	"github.com/lth-go/gogogogo/vm"
)

//
// Parameter 形参
//
type Parameter struct {
	typeSpecifier *TypeSpecifier

	name string
}

//
// FunctionDefinition 函数定义
//
type FunctionDefinition struct {
	typeSpecifier *TypeSpecifier

	packageNameList []string
	name            string

	parameterList []*Parameter
	block         *Block

	localVariableList []*Declaration

	index int
}

func (fd *FunctionDefinition) fix() {
	// 添加形参声明
	fd.addParameterAsDeclaration()
	fd.typeSpecifier.fix()

	if fd.block != nil {
		// 修正表达式列表
		fixStatementList(fd.block, fd.block.statementList, fd)

		// 修正返回值
		fd.addReturnFunction()
	}
}

func (fd *FunctionDefinition) typeS() *TypeSpecifier {
	return fd.typeSpecifier
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {

	for _, param := range fd.parameterList {
		if searchDeclaration(param.name, fd.block) != nil {
			compileError(param.typeSpecifier.Position(), PARAMETER_MULTIPLE_DEFINE_ERR, param.name)
		}
		decl := &Declaration{name: param.name, typeSpecifier: param.typeSpecifier}

		fd.block.addDeclaration(decl, fd, param.typeSpecifier.Position())
	}
}

func (fd *FunctionDefinition) addReturnFunction() {

	if fd.block.statementList == nil {
		ret := &ReturnStatement{returnValue: nil}
		ret.fix(fd.block, fd)
		fd.block.statementList = []Statement{ret}
		return
	}

	// TODO return 是否有必要一定最后
	last := fd.block.statementList[len(fd.block.statementList)-1]
	_, ok := last.(*ReturnStatement)
	if ok {
		return
	}

	ret := &ReturnStatement{returnValue: nil}
	ret.SetPosition(fd.typeSpecifier.Position())

	if ret.returnValue != nil {
		ret.returnValue.SetPosition(fd.typeSpecifier.Position())
	}
	ret.fix(fd.block, fd)
	fd.block.statementList = append(fd.block.statementList, ret)
}

func (fd *FunctionDefinition) addLocalVariable(decl *Declaration) {
	decl.variableIndex = len(fd.localVariableList)
	fd.localVariableList = append(fd.localVariableList, decl)
}

func (fd *FunctionDefinition) checkArgument(currentBlock *Block, argumentList []Expression, arrayBase *TypeSpecifier) {
	var tempType *TypeSpecifier

	parameterList := fd.parameterList

	paramLen := len(parameterList)
	argLen := len(argumentList)

	if argLen != paramLen {
		compileError(fd.typeS().Position(), ARGUMENT_COUNT_MISMATCH_ERR, paramLen, argLen)
	}

	for i := 0; i < paramLen; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)

		paramType := parameterList[i].typeSpecifier
		if paramType.basicType == vm.BaseType {
			tempType = arrayBase
		} else {
			tempType = paramType
		}
		argumentList[i] = createAssignCast(argumentList[i], tempType)
	}
}

func (fd *FunctionDefinition) getPackageName() string {
	return strings.Join(fd.packageNameList, ".")
}

func (fd *FunctionDefinition) getVmFuncName() string {
	return fd.name
}
