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

func (b *Block) show(ident int) {
	printWithIdent("Block", ident)
	subIdent := ident + 2

	for _, decl := range b.declarationList {
		decl.show(subIdent)
	}

	for _, stmt := range b.statementList {
		stmt.show(subIdent)
	}
}

func (b *Block) addDeclaration(decl *Declaration, fd *FunctionDefinition, pos Position) {
	if searchDeclaration(decl.name, b) != nil {
		compileError(pos, VARIABLE_MULTIPLE_DEFINE_ERR, decl.name)
	}

	if b != nil {
		b.declarationList = append(b.declarationList, decl)
		fd.addLocalVariable(decl)
		decl.isLocal = true
	} else {
		compiler := getCurrentCompiler()
		compiler.declarationList = append(compiler.declarationList, decl)
		decl.isLocal = false
	}
}
