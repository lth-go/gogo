package compiler

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
		f.codeList = fixOpcodeBuf(&ob)
		f.lineNumberList = ob.lineNumberList
	}
}

// 生成解释器所需的信息
func addTopLevel(compiler *Compiler, exe *Executable) {
	var ob OpcodeBuf

	init_opcode_buf(&ob)
	generateStatementList(exe, nil, compiler.statementList, &ob)

	exe.codeList = fixOpcodeBuf(&ob)
	exe.lineNumberList = ob.lineNumberList
}

//
// generateCode
//
func generateCode(ob *OpcodeBuf, pos Position, code Opcode, rest ...int) {
	// 获取参数类型
	paramList := []byte(opcodeInfo[code].parameter)

	ob.codeList = append(ob.codeList, byte(code))
	for i, param := range paramList {
		value := rest[i]
		switch param {
		case 'b': /* byte */
			ob.codeList = append(ob.codeList, byte(value))
		case 's': /* short(2byte int) */
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(value))
			ob.codeList = append(ob.codeList, b...)
		case 'p': /* constant pool index */
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(value))
			ob.codeList = append(ob.codeList, b...)
		default:
			panic("TODO")
		}
	}
	add_line_number(ob, pos.Line, len(ob.codeList))
}

//
// generateStatementList
//

func generateStatementList(exe *Executable, currentBlock *Block, statementList []Statement, ob *OpcodeBuf) {

	for _, stmt := range statementList {
		stmt.generate(exe, currentBlock, ob)
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

func copy_parameter_list(src []*Parameter) []*LocalVariable {
	var dest []*LocalVariable = []*LocalVariable{}

	for _, param := range src {
		dest = append(dest, &LocalVariable{
			name:          param.name,
			typeSpecifier: copy_type_specifier(param.typeSpecifier),
		})
	}

	return dest
}

func copyFunction(src *FunctionDefinition, dest *Function) {
	dest.typeSpecifier = copy_type_specifier(src.typeSpecifier)
	dest.name = src.name
	dest.parameterList = copy_parameter_list(src.parameterList)
	if src.block != nil {
		dest.localVariableList = copy_local_variables(src)
	} else {
		dest.localVariableList = nil
	}
}

func fixOpcodeBuf(ob *OpcodeBuf) []byte {

	fix_labels(ob)
	ob.labelTableList = nil

	return ob.codeList
}

// TODO 这是啥
func fix_labels(ob *OpcodeBuf) {

	for i := 0; i < len(ob.codeList); i++ {
		if ob.codeList[i] == byte(JUMP) || ob.codeList[i] == byte(JUMP_IF_TRUE) || ob.codeList[i] == byte(JUMP_IF_FALSE) {
			label := (ob.codeList[i+1] << 8) + (ob.codeList[i+2])
			address := ob.labelTableList[label].labelAddress
			ob.codeList[i+1] = (byte)(address >> 8)
			ob.codeList[i+2] = (byte)(address & 0xff)
		}
		info := &opcodeInfo[ob.codeList[i]]
		for _, p := range []byte(info.parameter) {
			switch p {
			case 'b':
				i++
			case 's': /* FALLTHRU */
				fallthrough
			case 'p':
				i += 2
			default:
				panic("param error")
			}
		}
	}
}

func add_line_number(ob *OpcodeBuf, lineNumber int, start_pc int) {
	if ob.lineNumberList == nil || (ob.lineNumberList[len(ob.lineNumberList)-1].lineNumber != lineNumber) {
		ob.lineNumberList = append(ob.lineNumberList, &LineNumber{
			lineNumber: lineNumber,
			startPc:    start_pc,
			pcCount:    len(ob.lineNumberList) - start_pc,
		})
	} else {
		ob.lineNumberList[len(ob.lineNumberList)-1].pcCount += len(ob.lineNumberList) - start_pc
	}
}

func copy_local_variables(fd *FunctionDefinition) []*LocalVariable {
	// TODO 形参占用位置
	var dest []*LocalVariable = []*LocalVariable{}

	for _, v := range fd.localVariableList {
		dest = append(dest, &LocalVariable{name: v.name, typeSpecifier: copy_type_specifier(v.typeSpecifier)})
	}

	return dest
}
