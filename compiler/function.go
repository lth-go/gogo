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
	Type            *Type
	PackageName     string
	Name            string
	Receiver        *Parameter
	ParameterList   []*Parameter
	Block           *Block
	DeclarationList []*Declaration
}

func (fd *FunctionDefinition) fix() {
	// 添加形参声明
	fd.addParameterAsDeclaration()
	fd.Type.Fix()

	if fd.Block != nil {
		// 修正表达式列表
		fd.Block.FixStatementList(fd)
		// 修正返回值
		fd.addReturnFunction()
	}
}

func (fd *FunctionDefinition) GetType() *Type {
	return fd.Type
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {
	for _, param := range fd.ParameterList {
		decl := NewDeclaration(param.Type.Position(), param.Type, param.Name, nil)
		fd.Block.AddDeclaration(decl, fd)
	}
}

func (fd *FunctionDefinition) addReturnFunction() {
	if fd.Block.statementList == nil {
		ret := &ReturnStatement{Value: nil}
		ret.fix(fd.Block, fd)
		fd.Block.statementList = []Statement{ret}
		return
	}

	// TODO return 是否有必要一定最后
	last := fd.Block.statementList[len(fd.Block.statementList)-1]
	_, ok := last.(*ReturnStatement)
	if ok {
		return
	}

	ret := &ReturnStatement{Value: nil}
	ret.SetPosition(fd.Type.Position())

	if ret.Value != nil {
		ret.Value.SetPosition(fd.Type.Position())
	}
	ret.fix(fd.Block, fd)
	fd.Block.statementList = append(fd.Block.statementList, ret)
}

func (fd *FunctionDefinition) AddDeclarationList(decl *Declaration) {
	decl.Index = len(fd.DeclarationList)
	fd.DeclarationList = append(fd.DeclarationList, decl)
}

func (fd *FunctionDefinition) checkArgument(currentBlock *Block, argumentList []Expression, arrayBase *Type) {
	var tempType *Type

	parameterList := fd.ParameterList

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
	return fd.PackageName
}

func (fd *FunctionDefinition) GetName() string {
	return fd.Name
}

// 拷贝函数定义的参数类型
func (fd *FunctionDefinition) CopyType() *Type {
	return fd.Type.Copy()
}
