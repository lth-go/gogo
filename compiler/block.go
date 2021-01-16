package compiler

//
// BlockInfo
//
type BlockInfo interface{}

type StatementBlockInfo struct {
	statement     Statement
	continueLabel int
	breakLabel    int
}

type FunctionBlockInfo struct {
	function *FunctionDefinition
	endLabel int
}

//
// Block
//
type Block struct {
	outerBlock *Block

	statementList   []Statement
	declarationList []*Declaration

	// 块信息，函数块，还是条件语句
	parent BlockInfo
}

func (b *Block) show(indent int) {
	printWithIndent("Block", indent)
	subIndent := indent + 2

	for _, declaration := range b.declarationList {
		declaration.show(subIndent)
	}

	for _, stmt := range b.statementList {
		stmt.show(subIndent)
	}
}

func (b *Block) addDeclaration(declaration *Declaration, fd *FunctionDefinition, pos Position) {
	if b.searchDeclaration(declaration.Name) != nil {
		compileError(pos, VARIABLE_MULTIPLE_DEFINE_ERR, declaration.Name)
	}

	if b != nil {
		b.declarationList = append(b.declarationList, declaration)
	}

	if fd != nil {
		declaration.IsLocal = true
		fd.addLocalVariable(declaration)
	} else {
		compiler := getCurrentCompiler()
		declaration.IsLocal = false
		compiler.AddDeclarationList(declaration)
	}
}

func (b *Block) getCurrentFunction() *FunctionDefinition {
	for block := b; block != nil; block = block.outerBlock {
		fdBlockInfo, ok := block.parent.(*FunctionBlockInfo)
		if ok {
			return fdBlockInfo.function
		}

	}
	return nil
}

func (b *Block) searchDeclaration(name string) *Declaration {
	// 从局部作用域查找
	for block := b; block != nil; block = block.outerBlock {
		for _, declaration := range block.declarationList {
			if declaration.Name == name {
				return declaration
			}
		}
	}

	// 从全局作用域查找
	compiler := getCurrentCompiler()
	for _, declaration := range compiler.declarationList {
		if declaration.Name == name {
			return declaration
		}
	}

	return nil
}
