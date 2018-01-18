package parser

type OpcodeBuf struct {
	size       int
	alloc_size int
	code       string
	LabelTable *label_table
	lineNumber []*LineNumber
}

func generate(compiler *Compiler) *Executable {
	exe := newExecutable()

	addGlobalVariable(compiler, exe)
	addFunctions(compiler, exe)
	addTopLevel(compiler, exe)

	return exe
}

func addGlobalVariable(compiler *Compiler, exe *Executable) {
	var v *Variable

	exe.globalVariableList = []*Variable{}

	for _, dl := range compiler.declarationList {
		v = &Variable{
			name:          dl.name,
			typeSpecifier: copy_type_specifier(dl.typeSpecifier),
		}

		exe.globalVariableList = append(exe.globalVariableList, v)
	}
}

// 为每个函数生成所需的信息
func addFunctions(compiler *Compiler, exe *Executable) {

	var ob OpcodeBuf

	var f *Function

	for _, fd := range compiler.funcList {
		f = &Function{}
		exe.functionList = append(exe.functionList, f)

		copyFunction(fd, f)
		if fd.block == nil {
			// 原生函数
			f.isImplemented = false
			continue
		}

		init_opcode_buf(&ob)
		generateStatementList(exe, fd.block, fd.block.statementList, &ob)

		f.isImplemented = true
		f.code_size = ob.size
		f.code = fix_opcode_buf(&ob)
		f.lineNumber_size = ob.lineNumber_size
		f.lineNumber = ob.lineNumber
		f.need_stack_size = calc_need_stack_size(f.code, f.code_size)
	}
}

// 生成解释器所需的信息
func addTopLevel(compiler *Compiler, exe *Executable) {
	var ob OpcodeBuf

	init_opcode_buf(&ob)
	generateStatementList(exe, nil, compiler.statementList, &ob)

	exe.code_size = ob.size
	exe.code = fix_opcode_buf(&ob)
	exe.lineNumber_size = ob.lineNumber_size
	exe.lineNumber = ob.lineNumber
	exe.need_stack_size = calc_need_stack_size(exe.code, exe.code_size)
}

func generateCode(ob []*Opcode, LineNumber int, code Opcode, rest ...int) {

}

//
// generateStatementList
//

func generateStatementList(exe *Executable, currentBlock *Block, statementList []*Statement, ob *OpcodeBuf) {

	for _, pos := range statementList {
		statement.generate(exe, currentBlock, ob)
	}
}
