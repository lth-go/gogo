package compiler

//
// Parameter 形参
//
type Parameter struct {
	Type     *Type
	Name     string
	Ellipsis bool
}

func NewParameter(typ *Type, name string, ellipsis bool) *Parameter {
	if ellipsis {
		typ = CreateArrayType(typ, typ.Position())
	}

	return &Parameter{
		Type:     typ,
		Name:     name,
		Ellipsis: ellipsis,
	}
}

//
// FunctionDefinition 函数定义
//
type FunctionDefinition struct {
	Type            *Type
	PackageName     string
	Name            string
	ParameterList   []*Parameter
	Block           *Block
	DeclarationList []*Declaration
}

func (fd *FunctionDefinition) Fix() {
	// 添加形参声明
	fd.addParameterAsDeclaration()
	fd.Type.Fix()

	if fd.Block != nil {
		// 修正表达式列表
		fd.Block.Fix()
		// 修正返回值
		fd.FixReturnStatement()
	}
}

func (fd *FunctionDefinition) GetType() *Type {
	return fd.Type
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {
	for _, param := range fd.ParameterList {
		decl := NewDeclaration(param.Type.Position(), param.Type, param.Name, nil)
		decl.IsLocal = true
		// TODO: 啥时候为空
		if fd.Block != nil {
			fd.Block.declarationList = append(fd.Block.declarationList, decl)
		}
		fd.AddDeclarationList(decl)
	}
}

// 确保函数语句里最后一定是return语句
// TODO: 校验参数类型
func (fd *FunctionDefinition) FixReturnStatement() {
	isNeedAddReturn := func() bool {
		if len(fd.Block.statementList) == 0 {
			return true
		}

		last := fd.Block.statementList[len(fd.Block.statementList)-1]
		_, ok := last.(*ReturnStatement)
		return !ok
	}()
	if !isNeedAddReturn {
		return
	}

	returnStmt := NewReturnStatement(fd.Type.Position(), nil)
	returnStmt.Block = fd.Block
	returnStmt.Fix()

	if fd.Block.statementList == nil {
		fd.Block.statementList = []Statement{}
	}
	fd.Block.statementList = append(fd.Block.statementList, returnStmt)
}

func (fd *FunctionDefinition) AddDeclarationList(decl *Declaration) {
	decl.Index = len(fd.DeclarationList)
	fd.DeclarationList = append(fd.DeclarationList, decl)
}

func (fd *FunctionDefinition) FixArgument(argumentList []Expression) []Expression {
	parameterList := fd.ParameterList

	paramLen := len(parameterList)

	if paramLen > 0 {
		lastP := parameterList[paramLen-1]
		if lastP.Ellipsis {
			newArgList := make([]Expression, 0)
			for _, expr := range argumentList[:paramLen-1] {
				newArgList = append(newArgList, expr)
			} 
			lastArg := CreateArrayExpression(lastP.Type, argumentList[paramLen-1:])
			newArgList = append(newArgList, lastArg)

			argumentList = newArgList
		}
	}

	argLen := len(argumentList)

	if argLen != paramLen {
		compileError(fd.GetType().Position(), ARGUMENT_COUNT_MISMATCH_ERR, paramLen, argLen)
	}

	for i := 0; i < paramLen; i++ {
		argumentList[i] = argumentList[i].Fix()
		if !argumentList[i].GetType().Equal(parameterList[i].Type) {
			compileError(
				argumentList[i].Position(),
				ARGUMENT_COUNT_MISMATCH_ERR,
				parameterList[i].Name,
				parameterList[i].Type.GetBasicType(),
				argumentList[i].GetType().GetBasicType(),
			)
		}
	}

	return argumentList
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
