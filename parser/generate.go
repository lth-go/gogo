package parser

func generate(compiler *Compiler) *Executable {
	exe := newExecutable()

	addGlobalVariable(compiler, exe)
	add_functions(compiler, exe)
	add_top_level(compiler, exe)

	return exe
}

func addGlobalVariable(compiler *Compiler, exe *Executable) {

	exe.globalVariableList = []*Variable{}

	for i, dl := range compiler.declarationList {
		exe.globalVariableList[i].name = dl.name
		exe.globalVariableList[i].typeSpecifier = copy_type_specifier(dl.typeSpecifier)
	}
}

func add_functions(compiler *Compiler, exe *Executable) {

	var ob OpcodeBuf

	for i, fd := range compiler.funcList {
		copy_function(fd, &exe.functionList[i])
		if fd.block == nil {
			// 原生函数
			exe.functionList[i].is_implemented = DVM_FALSE
			continue
		}

		init_opcode_buf(&ob)
		generate_statement_list(exe, fd.block, fd.block.statementList, &ob)

		exe.functionList[i].is_implemented = true
		exe.functionList[i].code_size = ob.size
		exe.functionList[i].code = fix_opcode_buf(&ob)
		exe.functionList[i].line_number_size = ob.line_number_size
		exe.functionList[i].line_number = ob.line_number
		exe.functionList[i].need_stack_size = calc_need_stack_size(exe.functionList[i].code, exe.functionList[i].code_size)
	}

}

func add_top_level(compiler *Compiler, exe *Executable) {
	var ob OpcodeBuf

	init_opcode_buf(&ob)
	generate_statement_list(exe, nil, compiler.statementList, &ob)

	exe.code_size = ob.size
	exe.code = fix_opcode_buf(&ob)
	exe.line_number_size = ob.line_number_size
	exe.line_number = ob.line_number
	exe.need_stack_size = calc_need_stack_size(exe.code, exe.code_size)
}
