package parser

type OpcodeBuf struct {
	size       int
	alloc_size int
	code       string
	LabelTable *label_table
	lineNumber []*LineNumber
}

func init_opcode_buf(ob *OpcodeBuf) {
	// TODO
	ob.size = 0
	ob.alloc_size = 0
	ob.code = NULL
	ob.label_table_size = 0
	ob.label_table_alloc_size = 0
	ob.label_table = nil
	ob.line_number_size = 0
	ob.line_number = nil
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

//
// generateCode
//
func generateCode(ob []*Opcode, pos Position, code Opcode, rest ...int) {
	line_number := pos.Line

	// 获取参数类型
	param := dvm_opcode_info[code].parameter

	start_pc = ob.size
	ob.code[ob.size] = code
	ob.size++
	for i := 0; param[i] != '0'; i++ {
		for _, value := range rest {
			switch param[i] {
			case 'b': /* byte */
				ob.code[ob.size] = value
				ob.size++
			case 's': /* short(2byte int) */
				ob.code[ob.size] = (value >> 8)
				ob.code[ob.size+1] = (value & 0xff)
				ob.size += 2
			case 'p': /* constant pool index */
				ob.code[ob.size] = (value >> 8)
				ob.code[ob.size+1] = (value & 0xff)
				ob.size += 2
			default:
				panic("TODO")
			}
		}
	}
	add_line_number(ob, line_number, start_pc)
}

//
// generateStatementList
//

func generateStatementList(exe *Executable, currentBlock *Block, statementList []*Statement, ob *OpcodeBuf) {

	for _, pos := range statementList {
		statement.generate(exe, currentBlock, ob)
	}
}
