package compiler

import (
	"encoding/binary"

	"../vm"
)

type OpcodeBuf struct {
	codeList       []byte
	labelTableList []*LabelTable
	lineNumberList []*vm.VmLineNumber
}
type LabelTable struct {
	labelAddress int
}

func iniCodeBuf(ob *OpcodeBuf) {
	ob.codeList = []byte{}
	ob.labelTableList = []*LabelTable{}
	ob.lineNumberList = []*vm.VmLineNumber{}
}

func addGlobalVariable(compiler *Compiler, exe *vm.Executable) {
	var v *vm.VmVariable

	exe.globalVariableList = []*vm.VmVariable{}

	for _, dl := range compiler.declarationList {
		v = &vm.VmVariable{
			name:          dl.name,
			typeSpecifier: copyTypeSpecifier(dl.typeSpecifier),
		}

		exe.globalVariableList = append(exe.globalVariableList, v)
	}
}

// 为每个函数生成所需的信息
func addFunctions(compiler *Compiler, exe *vm.Executable) {

	var ob OpcodeBuf

	var f *vm.VmFunction

	for _, fd := range compiler.funcList {
		f = &vm.VmFunction{}
		exe.functionList = append(exe.functionList, f)

		copyFunction(fd, f)
		if fd.block == nil {
			// 原生函数
			f.isImplemented = false
			continue
		}

		iniCodeBuf(&ob)
		generateStatementList(exe, fd.block, fd.block.statementList, &ob)

		f.isImplemented = true
		f.codeList = fixOpcodeBuf(&ob)
		f.lineNumberList = ob.lineNumberList
	}
}

// 生成解释器所需的信息
func addTopLevel(compiler *Compiler, exe *vm.Executable) {
	var ob OpcodeBuf

	iniCodeBuf(&ob)
	generateStatementList(exe, nil, compiler.statementList, &ob)

	exe.codeList = fixOpcodeBuf(&ob)
	exe.lineNumberList = ob.lineNumberList
}

//
// generateCode
//
func generateCode(ob *OpcodeBuf, pos Position, code byte, rest ...int) {
	// 获取参数类型
	paramList := []byte(vm.opcodeInfo[int(code)].parameter)

	startPc = len(ob.codeList)
	ob.codeList = append(ob.codeList, code)

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
	addLineNumber(ob, pos.Line, startPc)
}

func addLineNumber(ob *OpcodeBuf, lineNumber int, start_pc int) {
	if ob.lineNumberList == nil || (ob.lineNumberList[len(ob.lineNumberList)-1].lineNumber != lineNumber) {
		l := &vm.VmLineNumber{
			lineNumber: lineNumber,
			startPc:    start_pc,
			pcCount:    len(ob.codeList) - start_pc,
		}
		ob.lineNumberList = append(ob.lineNumberList, l)
	} else {
		ob.lineNumberList[len(ob.lineNumberList)-1].pcCount += len(ob.codeList) - start_pc
	}
}

//
// generateStatementList
//

func generateStatementList(exe *vm.Executable, currentBlock *Block, statementList []Statement, ob *OpcodeBuf) {

	for _, stmt := range statementList {
		stmt.generate(exe, currentBlock, ob)
	}
}

//
// COPY
//

func copyTypeSpecifier(src *TypeSpecifier) *vm.VmTypeSpecifier {

	dest := &vm.VmTypeSpecifier{basicType: src.basicType}

	dest.deriveList = []TypeDerive{}

	for _, derive := range src.deriveList {
		switch f := derive.(type) {
		case *FunctionDerive:
			newDerive := &vm.VmFunctionDerive{parameterList: copyParameterList(f.parameterList)}
			dest.deriveList = append(dest.deriveList, newDerive)
		default:
			panic("derive error")
		}
	}

	return dest
}

func copyParameterList(src []*Parameter) []*vm.VmLocalVariable {
	dest := []*vm.VmLocalVariable{}

	for _, param := range src {
		v := &vm.VmLocalVariable{
			name:          param.name,
			typeSpecifier: copyTypeSpecifier(param.typeSpecifier),
		}
		dest = append(dest, v)
	}
	return dest
}

func copyFunction(src *FunctionDefinition, dest *vm.VmFunction) {
	dest.typeSpecifier = copyTypeSpecifier(src.typeSpecifier)
	dest.name = src.name
	dest.parameterList = copyParameterList(src.parameterList)
	if src.block != nil {
		dest.localVariableList = copy_local_variables(src)
	} else {
		dest.localVariableList = nil
	}
}

func copy_local_variables(fd *FunctionDefinition) []*vm.VmLocalVariable {
	// TODO 形参占用位置
	var dest []*vm.VmLocalVariable = []*vm.VmLocalVariable{}

	localVariableCount = len(fd.localVariableList) - len(fd.parameterList)

	for _, v := range fd.localVariableList[0:localVariableCount] {
		vmV := &vm.VmLocalVariable{
			name:          v.name,
			typeSpecifier: copyTypeSpecifier(v.typeSpecifier),
		}
		dest = append(dest, vmV)
	}

	return dest
}

//
// FIX
//

func fixOpcodeBuf(ob *OpcodeBuf) []byte {

	fixLabels(ob)
	ob.labelTableList = nil

	return ob.codeList
}

// TODO 这是啥
func fixLabels(ob *OpcodeBuf) {

	for i := 0; i < len(ob.codeList); i++ {
		if ob.codeList[i] == vm.VM_JUMP ||
			ob.codeList[i] == vm.VM_JUMP_IF_TRUE ||
			ob.codeList[i] == vm.VM_JUMP_IF_FALSE {

			label := int((ob.codeList[i+1] << 8) + ob.codeList[i+2])
			address := ob.labelTableList[label].labelAddress
			ob.codeList[i+1] = (byte)(address >> 8)
			ob.codeList[i+2] = (byte)(address & 0xff)
		}
		info := &vm.opcodeInfo[ob.codeList[i]]
		for _, p := range []byte(info.parameter) {
			switch p {
			case 'b':
				i++
			case 's': /* FALLTHRU */
				// TODO 这不是报错了么
				fallthrough
			case 'p':
				i += 2
			default:
				panic("param error")
			}
		}
	}
}
