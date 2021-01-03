package compiler

const defaultPackage = "gogogogo.lang"

func setRequireList(requireList []*ImportSpec) {

	compiler := getCurrentCompiler()

	//current_package_name = strings.Join(compiler.packageNameList, '.')

	// 添加默认包
	//if current_package_name != defaultPackage {
	//    requireLists = add_default_package(requireLists)
	//}

	compiler.requireList = requireList
}

func createFunctionDefinition(typ *TypeSpecifier, identifier string, parameterLists []*Parameter, block *Block) *FunctionDefinition {
	compiler := getCurrentCompiler()

	fd := &FunctionDefinition{}

	fd.typeSpecifier = typ
	fd.packageNameList = compiler.packageNameList
	fd.name = identifier
	fd.parameterList = parameterLists
	fd.block = block

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	compiler.funcList = append(compiler.funcList, fd)

	return fd
}

func createFunctionDeriveType(fd *FunctionDefinition) *TypeSpecifier {

	ret := &TypeSpecifier{basicType: fd.typeSpecifier.basicType}

	*ret = *(fd.typeSpecifier)

	funcDerive := &FunctionDerive{parameterList: fd.parameterList}
	ret.deriveList = []TypeDerive{funcDerive}

	ret.deriveList = append(ret.deriveList, fd.typeSpecifier.deriveList...)

	return ret
}

// yacc类创建
func startClassDefine(identifier string, extends []*Extend, pos Position) {
	compiler := getCurrentCompiler()

	cd := &ClassDefinition{}

	cd.packageNameList = compiler.packageNameList
	cd.name = identifier
	cd.extendList = extends

	cd.SetPosition(pos)

	if compiler.currentClassDefinition != nil {
		panic("TODO")
	}

	compiler.currentClassDefinition = cd
}

func endClassDefine(memberList []MemberDeclaration) {
	compiler := getCurrentCompiler()

	cd := compiler.currentClassDefinition

	if cd == nil {
		panic("TODO")
	}

	if compiler.classDefinitionList == nil {
		compiler.classDefinitionList = []*ClassDefinition{}
	}
	compiler.classDefinitionList = append(compiler.classDefinitionList, cd)

	cd.memberList = memberList
	compiler.currentClassDefinition = nil
}

// 类方法定义
func methodFunctionDefine(typ *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) *FunctionDefinition {

	fd := createFunctionDefinition(typ, identifier, parameterList, block)

	return fd
}
