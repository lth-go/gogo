package compiler

//
// Parameter 形参
//
type Parameter struct {
	Type *Type
	Name string
}

func NewParameter(typ *Type, name string) *Parameter {
	return &Parameter{
		Type: typ,
		Name: name,
	}
}

//
// FunctionDefinition 函数定义
//
type FunctionDefinition struct {
	Type              *Type
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
	fd.Type.fix()

	if fd.block != nil {
		// 修正表达式列表
		FixStatementList(fd.block, fd.block.statementList, fd)
		// 修正返回值
		fd.addReturnFunction()
	}
}

func (fd *FunctionDefinition) GetType() *Type {
	return fd.Type
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {
	for _, param := range fd.parameterList {
		if fd.block.searchDeclaration(param.Name) != nil {
			compileError(param.Type.Position(), PARAMETER_MULTIPLE_DEFINE_ERR, param.Name)
		}

		decl := &Declaration{Name: param.Name, Type: param.Type}
		fd.block.addDeclaration(decl, fd, param.Type.Position())
	}
}

func (fd *FunctionDefinition) addReturnFunction() {
	if fd.block.statementList == nil {
		ret := &ReturnStatement{Value: nil}
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

	ret := &ReturnStatement{Value: nil}
	ret.SetPosition(fd.Type.Position())

	if ret.Value != nil {
		ret.Value.SetPosition(fd.Type.Position())
	}
	ret.fix(fd.block, fd)
	fd.block.statementList = append(fd.block.statementList, ret)
}

func (fd *FunctionDefinition) addLocalVariable(decl *Declaration) {
	decl.Index = len(fd.localVariableList)
	fd.localVariableList = append(fd.localVariableList, decl)
}

func (fd *FunctionDefinition) checkArgument(currentBlock *Block, argumentList []Expression, arrayBase *Type) {
	var tempType *Type

	parameterList := fd.parameterList

	paramLen := len(parameterList)
	argLen := len(argumentList)

	if argLen != paramLen {
		compileError(fd.GetType().Position(), ARGUMENT_COUNT_MISMATCH_ERR, paramLen, argLen)
	}

	for i := 0; i < paramLen; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)

		paramType := parameterList[i].Type
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
	return fd.Type.Copy()
}
