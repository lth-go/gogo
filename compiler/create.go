package compiler

const DEFAULT_PACKAGE = "gogogogo.lang"

func setRequireList(require_list []*Require) {

	compiler := getCurrentCompiler()

	//current_package_name = strings.Join(compiler.packageNameList, '.')

	// 添加默认包
	//if current_package_name != DEFAULT_PACKAGE {
	//    require_list = add_default_package(require_list)
	//}

	compiler.requireList = require_list
}

func createFunctionDefinition(typ *TypeSpecifier, identifier string, parameter_list []*Parameter, block *Block) *FunctionDefinition {
	compiler := getCurrentCompiler()

	fd := &FunctionDefinition{}

	fd.typeSpecifier = typ
	fd.packageNameList = compiler.packageNameList
	fd.name = identifier
	fd.parameterList = parameter_list
	fd.block = block

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	compiler.funcList = append(compiler.funcList, fd)

	return fd
}

func create_function_derive_type(fd *FunctionDefinition) *TypeSpecifier {

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
	cd.extends = extends

	cd.SetPosition(pos)

	if compiler.currentClassDefinition != nil {
		panic("TODO")
	}

	compiler.currentClassDefinition = cd
}

func endClassDefine(member_list []MemberDeclaration) {
	compiler := getCurrentCompiler()

	cd := compiler.currentClassDefinition

	if cd == nil {
		panic("TODO")
	}

	if compiler.classDefinitionList == nil {
		compiler.classDefinitionList = []*ClassDefinition{}
	}
	compiler.classDefinitionList = append(compiler.classDefinitionList, cd)

	cd.member = member_list
	compiler.currentClassDefinition = nil
}

// 类方法定义
func methodFunctionDefine(typ *TypeSpecifier, identifier string, parameter_list []*Parameter, block *Block) *FunctionDefinition {

	fd := createFunctionDefinition(typ, identifier, parameter_list, block)

	return fd
}
