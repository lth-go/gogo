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
	Block           *Block
	DeclarationList []*Declaration
	CodeList        []byte
}

// Fix
func (fd *FunctionDefinition) Fix() {
	if fd.Block == nil {
		return
	}

	// 添加形参声明
	fd.FixParam()
	// fd.Type.Fix()

	// 修正表达式列表
	fd.FixBlock()

	// 修正返回值
	fd.FixReturn()
}

func (fd *FunctionDefinition) GetType() *Type {
	return fd.Type
}

// 将形参添加到函数块声明列表,用于函数语句查找变量
// 实际栈位置将修改为函数参数位置
func (fd *FunctionDefinition) FixParam() {
	for i, param := range fd.Type.funcType.Params {
		decl := &Declaration{
			Type:        param.Type,
			PackageName: fd.PackageName,
			Name:        param.Name,
			Value:       nil,
			Index:       i,
			Block:       nil,
			IsLocal:     true,
		}
		fd.Block.declarationList = append(fd.Block.declarationList, decl)
	}
}

func (fd *FunctionDefinition) FixBlock() {
	fd.Block.Fix()
}

// FixReturn
// 确保函数语句里最后一定是return语句
// TODO: 校验参数类型
func (fd *FunctionDefinition) FixReturn() {
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
