package parser

import (
	"encoding/binary"
)

type OpcodeBuf struct {
	codeList       []byte
	labelTableList []*LabelTable
	lineNumberList []*LineNumber
}
type LabelTable struct {
	labelAddress int
}

func init_opcode_buf(ob *OpcodeBuf) {
	ob.codeList = []byte{}
	ob.labelTableList = []*LabelTable{}
	ob.lineNumberList = []*LineNumber{}
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
		f.lineNumberLine = ob.lineNumberLine
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
	exe.lineNumberList = ob.lineNumberList
}

//
// generateCode
//
func generateCode(ob *OpcodeBuf, pos Position, code Opcode, rest ...int) {
	// 获取参数类型
	paramList := []byte(opcode_info[code].parameter)

	ob.codeList = append(ob.codeList, code)
	for i, param := range paramList {
		value := rest[i]
		switch param {
			case "b": /* byte */
			ob.codeList = append(ob.codeList, value)
			case "s": /* short(2byte int) */
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(value))
			ob.codeList = append(ob.codeList, b...)
			case "p": /* constant pool index */
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(value))
			ob.codeList = append(ob.codeList, b...)
		case "":
		default:
			panic("TODO")
		}
	}
	add_line_number(ob, pos.Line, len(ob.codeList))
}

//
// generateStatementList
//

func generateStatementList(exe *Executable, currentBlock *Block, statementList []*Statement, ob *OpcodeBuf) {

	for _, pos := range statementList {
		statement.generate(exe, currentBlock, ob)
	}
}

func copy_type_specifier(src *TypeSpecifier) *TypeSpecifier {

	dest := &TypeSpecifier{basicType: src.basicType}

	dest.deriveList = []TypeDerive{}

	for _, derive := range src.deriveList {
		switch f := derive.(type) {
		case *FunctionDerive:
			newDerive := &FunctionDerive{parameterList: copy_parameter_list(f.parameterList)}
			dest.deriveList = append(dest.deriveList, newDerive)
		default:
			panic("derive error")
		}
	}

	return dest
}


//func copy_parameter_list(src *ParameterList, int *param_count_p) *LocalVariable{
//    int param_count = 0;
//    ParameterList *param;
//    DVM_LocalVariable *dest;
//    int i;

//    for (param = src; param; param = param->next) {
//        param_count++;
//    }
//    *param_count_p = param_count;
//    dest = MEM_malloc(sizeof(DVM_LocalVariable) * param_count);

//    for (param = src, i = 0; param; param = param->next, i++) {
//        dest[i].name = MEM_strdup(param->name);
//        dest[i].type = copy_type_specifier(param->type);
//    }

//    return dest;
//}
