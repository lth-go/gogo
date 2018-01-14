package parser

// 修正树
func fixTree(c *Compiler) {
	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	for _, funcPos := range c.funcList {
		// TODO 为何跳过
		if funcPos.block == nil {
			continue
		}

		// 添加形参声明
		addParameterAsDeclaration(funcPos)

		// 修正表达式列表
		fixStatementList(funcPos.block, funcPos.block.statementList, funcPos)

		// 修正返回值
		addReturnFunction(funcPos)
	}

	for varCount, dl := range c.declarationList {
		dl.variableIndex = varCount
	}

}

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.Fix(currentBlock, statement, fd)
	}
}
