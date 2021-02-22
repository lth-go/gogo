package compiler

// StatementBlockInfo 语句块
type StatementBlockInfo struct {
	Statement     Statement
	ContinueLabel int
	BreakLabel    int
}

func NewStatementBlockInfo(statement Statement) *StatementBlockInfo {
	return &StatementBlockInfo{
		Statement: statement,
	}
}

// FunctionBlockInfo 函数块
type FunctionBlockInfo struct {
	Function *FunctionDefinition
	EndLabel int
}

//
// Block
//
type Block struct {
	parent          interface{}    // 块信息，函数块，还是条件语句
	outerBlock      *Block         // 上级块
	declarationList []*Declaration // 用于搜索作用域
	statementList   []Statement    // 语句
}

// GetCurrentFunction
func (b *Block) GetCurrentFunction() *FunctionDefinition {
	for block := b; block != nil; block = block.outerBlock {
		blockInfo, ok := block.parent.(*FunctionBlockInfo)
		if ok {
			return blockInfo.Function
		}

	}
	return nil
}

func (b *Block) SearchDeclaration(name string) *Declaration {
	// 从局部作用域查找
	for block := b; block != nil; block = block.outerBlock {
		for _, decl := range block.declarationList {
			if decl.Name == name {
				return decl
			}
		}
	}

	return nil
}

func (b *Block) Fix() {
	for _, statement := range b.statementList {
		statement.Fix()
	}
}
