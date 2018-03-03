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
	}
	if fd != nil {
		decl.isLocal = true
		fd.addLocalVariable(decl)
	} else {
		compiler := getCurrentCompiler()
		decl.isLocal = false
		compiler.declarationList = append(compiler.declarationList, decl)
	}
}

func (b *Block) getCurrentFunction() *FunctionDefinition {

	for block := b; ; block = block.outerBlock {
		fdBlockInfo, ok := block.parent.(*FunctionBlockInfo)
		if ok {
			return fdBlockInfo.function
		}

	}
	return nil
}
