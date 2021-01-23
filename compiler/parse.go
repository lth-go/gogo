package compiler

func AddDeclList(decl *Declaration) {
	c := GetCurrentCompiler()
	decl.PackageName = c.packageName
	c.AddDeclarationList(decl)
}

func SetPackageName(packageName string) {
	c := GetCurrentCompiler()
	c.packageName = packageName
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
