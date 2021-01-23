package compiler

type BlockInfo interface{}

type StatementBlockInfo struct {
	statement     Statement
	continueLabel int
	breakLabel    int
}

func NewStatementBlockInfo(statement Statement) *StatementBlockInfo {
	return &StatementBlockInfo{
		statement: statement,
	}
}

type FunctionBlockInfo struct {
	function *FunctionDefinition
	endLabel int
}

//
// Block
//
type Block struct {
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration
	parent          BlockInfo // 块信息，函数块，还是条件语句
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

func (b *Block) SearchDeclaration(name string) *Declaration {
	// 从局部作用域查找
	for block := b; block != nil; block = block.outerBlock {
		for _, declaration := range block.declarationList {
			if declaration.Name == name {
				return declaration
			}
		}
	}

	// 从全局作用域查找
	return GetCurrentCompiler().SearchDeclaration(name)
}

func (b *Block) Fix() {
	for _, statement := range b.statementList {
		statement.Fix()
	}
}
