package compiler

const defaultPackage = "gogogogo.lang"

func setImportList(importList []*ImportSpec) {

	compiler := getCurrentCompiler()

	//current_package_name = strings.Join(compiler.packageNameList, '.')

	// 添加默认包
	//if current_package_name != defaultPackage {
	//    requireLists = add_default_package(requireLists)
	//}

	compiler.importList = importList
}

func createFunctionDeriveType(fd *FunctionDefinition) *TypeSpecifier {

	ret := newTypeSpecifier(fd.typeSpecifier.basicType)

	*ret = *(fd.typeSpecifier)

	ret.deriveType = &FunctionDerive{parameterList: fd.parameterList}

	return ret
}
