package compiler

//
// Parameter 形参
//
type Parameter struct {
	typeSpecifier *Type
	name          string
}

//
// FunctionDefinition 函数定义
//
type FunctionDefinition struct {
	typeSpecifier     *Type
	packageName       string
	name              string
	receiver          *Parameter
	parameterList     []*Parameter
	block             *Block
	localVariableList []*Declaration
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

func (fd *FunctionDefinition) typeS() *Type {
	return fd.typeSpecifier
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {
	for _, param := range fd.parameterList {
		if fd.block.searchDeclaration(param.name) != nil {
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

func (fd *FunctionDefinition) checkArgument(currentBlock *Block, argumentList []Expression, arrayBase *Type) {
	var tempType *Type

	parameterList := fd.parameterList

	paramLen := len(parameterList)
	argLen := len(argumentList)

	if argLen != paramLen {
		compileError(fd.typeS().Position(), ARGUMENT_COUNT_MISMATCH_ERR, paramLen, argLen)
	}

	for i := 0; i < paramLen; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)

		paramType := parameterList[i].typeSpecifier
		tempType = paramType
		argumentList[i] = CreateAssignCast(argumentList[i], tempType)
	}
}

func (fd *FunctionDefinition) GetPackageName() string {
	return fd.packageName
}

func (fd *FunctionDefinition) GetName() string {
	return fd.name
}

// 拷贝函数定义的参数类型
func (fd *FunctionDefinition) CopyType() *Type {
	typ := fd.typeSpecifier.CopyType()
	typ.funcType = NewFuncType(fd.parameterList, nil)

	return typ
}
