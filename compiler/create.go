package compiler

const DEFAULT_PACKAGE = "gogogogo.lang"

func setRequireList(require_list []*Require) {

	//compiler := getCurrentCompiler()

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

	ret.deriveList = []ArrayDerive{&FunctionDerive{parameterList: fd.parameterList}}

	ret.deriveList = append(ret.deriveList, fd.typeSpecifier.deriveList...)

	return ret
}
