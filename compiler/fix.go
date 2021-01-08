package compiler

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

// 仅限函数
func fixModuleMemberExpression(expr *MemberExpression, memberName string) Expression {
	innerExpr := expr.expression

	innerExpr.typeS().fix()

	module := innerExpr.(*IdentifierExpression).inner.(*Module)

	moduleCompiler := module.compiler

	fd := moduleCompiler.searchFunction(memberName)
	if fd == nil {
		panic("TODO")
	}

	// TODO 得用当前compiler来添加
	currentCompiler := getCurrentCompiler()

	newExpr := &IdentifierExpression{
		name: memberName,
		inner: &FunctionIdentifier{
			functionDefinition: fd,
			functionIndex:      currentCompiler.addToVmFunctionList(fd),
		},
	}

	newExpr.setType(createFuncType(fd))
	newExpr.typeS().fix()

	return newExpr
}
