package compiler

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

func check_member_accessibility(pos Position, targetClass *ClassDefinition, member MemberDeclaration, memberName string) {
	compiler := getCurrentCompiler()

	if compiler.getPackageName() != targetClass.getPackageName() {
		compileError(pos, PACKAGE_MEMBER_ACCESS_ERR, memberName)
	}
}

func fixClassMemberExpression(expr *MemberExpression, obj Expression, memberName string) Expression {

	obj.typeS().fix()

	cd := obj.typeS().classRef.classDefinition

	member := cd.searchMember(memberName)
	if member == nil {
		compileError(expr.Position(), MEMBER_NOT_FOUND_ERR, cd.name, memberName)
	}

	check_member_accessibility(obj.Position(), cd, member, memberName)

	expr.declaration = member

	switch m := member.(type) {
	case *MethodMember:
		expr.setType(createFunctionDeriveType(m.functionDefinition))
	case *FieldMember:
		expr.setType(m.typeSpecifier)
	}

	return expr

}

func createFunctionDeriveType(fd *FunctionDefinition ) *TypeSpecifier {
	typ := &TypeSpecifier{}

	*typ = *fd.typeS()

	newFuncDerive := &FunctionDerive{
		parameterList: fd.parameterList,
	}
	typ.appendDerive(newFuncDerive)

	typ.deriveList = append(typ.deriveList, fd.typeS().deriveList...)

    return typ
}

