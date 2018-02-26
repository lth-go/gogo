package compiler

func check_member_accessibility(pos Position, targetClass *ClassDefinition, member MemberDeclaration, memberName string) {
	compiler := getCurrentCompiler()

	if compiler.getPackageName() != targetClass.getPackageName() {
		compileError(pos, PACKAGE_MEMBER_ACCESS_ERR, memberName)
	}
}

func fixClassMemberExpression(expr *MemberExpression, obj Expression, memberName string) Expression {
	var target_interface *ClassDefinition
	var interface_index int

	obj.typeS().fix()

	cd := obj.typeS().classRef.classDefinition

	member := search_member(cd, memberName)
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
