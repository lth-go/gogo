package compiler

//
// 函数定义
//
func CreateFunctionDefine(pos Position, receiver *Parameter, identifier string, typ *Type, block *Block) {
	c := GetCurrentCompiler()

	fd := &FunctionDefinition{
		Type:            typ,
		Name:            identifier,
		PackageName:     c.GetPackageName(),
		ParamList:   typ.funcType.Params,
		Block:           block,
		DeclarationList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

func AddDeclList(decl *Declaration) {
	c := GetCurrentCompiler()
	decl.PackageName = c.GetPackageName()
	c.AddDeclarationList(decl)
}

func SetPackageName(packageName string) {
	c := GetCurrentCompiler()
	c.SetPackageName(packageName)
}

func SetImportList(importList []*Import) {
	c := GetCurrentCompiler()
	c.importList = importList
}

func PushCurrentBlock() *Block {
	c := GetCurrentCompiler()
	c.currentBlock = &Block{outerBlock: c.currentBlock}

	return c.currentBlock
}

func PopCurrentBlock() *Block {
	c := GetCurrentCompiler()
	b := c.currentBlock
	c.currentBlock = c.currentBlock.outerBlock

	return b
}
