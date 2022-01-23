package compiler

//
// 函数定义
//
func CreateFunctionDefine(pos Position, receiver *Parameter, identifier string, typ *Type, block *Block) {
	c := GetCurrentPackage()

	fd := &FunctionDefinition{
		Type:            typ,
		Name:            identifier,
		PackageName:     c.GetPackageName(),
		Block:           block,
		DeclarationList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{Function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

func CreateDeclaration(pos Position, typ *Type, name string, value Expression) *Declaration {
	decl := NewDeclaration(pos, typ, name, value)

	decl.Block = GetCurrentPackage().currentBlock
	if decl.Block != nil {
		decl.IsLocal = true
	}

	return decl
}

func AddDeclList(stmt Statement) {
	decl := stmt.(*Declaration)

	c := GetCurrentPackage()
	decl.PackageName = c.GetPackageName()

	c.declarationList = append(c.declarationList, decl)
}

func SetPackageName(packageName string) {
	c := GetCurrentPackage()
	c.SetPackageName(packageName)
}

func SetImportList(importList []*Import) {
	c := GetCurrentPackage()
	c.importList = importList
}

func PushCurrentBlock() *Block {
	c := GetCurrentPackage()
	c.currentBlock = &Block{outerBlock: c.currentBlock}

	return c.currentBlock
}

func PopCurrentBlock() *Block {
	c := GetCurrentPackage()
	b := c.currentBlock
	c.currentBlock = c.currentBlock.outerBlock

	return b
}
