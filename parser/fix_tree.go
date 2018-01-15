package parser

// 修正树
func fixTree(c *Compiler) {
	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	for _, fd := range c.funcList {
		if fd.block == nil {
			continue
		}

		// 添加形参声明
		addParameterAsDeclaration(fd)

		// 修正表达式列表
		fixStatementList(fd.block, fd.block.statementList, fd)

		// 修正返回值
		addReturnFunction(fd)
	}

}

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

// ==============================
// utils
// ==============================

func addParameterAsDeclaration(fd *FunctionDefinition) {

	for _, param := range fd.parameterList {
		if searchDeclaration(param.name, fd.block) != nil {
			compileError(param.Position(), 0, "")
		}
		decl := &Declaration{name: param.name, typeSpecifier: param.typeSpecifier}

		addDeclaration(fd.block, decl, fd, param.Position())
	}
}

func addDeclaration(currentBlock *Block, decl *Declaration, fd *FunctionDefinition, pos Position) {
	if searchDeclaration(decl.name, currentBlock) {
		compileError(pos, 0, "")
	}

	if currentBlock != nil {
		currentBlock.declarationList = append(currentBlock.declarationList, decl)
		addLocalVariable(fd, decl)
		decl.isLocal = BooleanTrue
	} else {
		compiler := getCurrentCompiler()
		compiler.declarationList = append(compiler.declarationList, decl)
		decl.isLocal = BooleanFalse
	}

}

func addReturnFunction(fd *FunctionDefinition) {

	if fd.block.statementList == nil {
		ret := &ReturnStatement{returnValue: nil}
		ret.fix(fd.block, fd)
		fd.block.statementList = []Statement{ret}
		return
	}

	last := fd.block.statementList[-1]
	_, ok := last.(ReturnStatement)
	if ok {
		return
	}
	ret := &ReturnStatement{returnValue: nil}
	ret.fix(fd.block, fd)
	fd.block.statementList = append(fd.block.statementList, ret)
	return
}
