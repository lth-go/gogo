package compiler

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

// 仅限函数
func fixPackageMemberExpression(expr *MemberExpression, memberName string) Expression {
	innerExpr := expr.expression

	innerExpr.typeS().fix()

	p := innerExpr.(*IdentifierExpression).inner.(*Package)

	fd := p.compiler.searchFunction(memberName)
	if fd == nil {
		panic("TODO")
	}

	// TODO 得用当前compiler来添加
	currentCompiler := getCurrentCompiler()

	newExpr := &IdentifierExpression{
		name: memberName,
		inner: &FunctionIdentifier{
			functionDefinition: fd,
			Index:              currentCompiler.AddFuncList(fd),
		},
	}

	newExpr.setType(createFuncType(fd))
	newExpr.typeS().fix()

	return newExpr
}
