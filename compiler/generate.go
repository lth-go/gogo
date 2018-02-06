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

	exe.GlobalVariableList = []*vm.VmVariable{}

	for _, dl := range compiler.declarationList {
		v = vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))

		exe.GlobalVariableList = append(exe.GlobalVariableList, v)
	}
}

// 为每个函数生成所需的信息
func addFunctions(compiler *Compiler, exe *vm.Executable) {

	var ob OpcodeBuf

	var f *vm.VmFunction

	for _, fd := range compiler.funcList {
		f = &vm.VmFunction{}
		exe.FunctionList = append(exe.FunctionList, f)

		copyFunction(fd, f)
		if fd.block == nil {
			// 原生函数
			f.IsImplemented = false
			continue
		}

		iniCodeBuf(&ob)
		generateStatementList(exe, fd.block, fd.block.statementList, &ob)

		f.IsImplemented = true
		f.CodeList = fixOpcodeBuf(&ob)
		f.LineNumberList = ob.lineNumberList
	}
}

// 生成解释器所需的信息
func addTopLevel(compiler *Compiler, exe *vm.Executable) {
	var ob *OpcodeBuf

	iniCodeBuf(ob)
	generateStatementList(exe, nil, compiler.statementList, ob)

	exe.CodeList = fixOpcodeBuf(ob)
	exe.LineNumberList = ob.lineNumberList
}

//
// generateCode
//
func generateCode(ob *OpcodeBuf, pos Position, code byte, rest ...int) {
	// 获取参数类型
	paramList := []byte(vm.OpcodeInfo[int(code)].Parameter)

	startPc := len(ob.codeList)
	ob.codeList = append(ob.codeList, code)

	for i, param := range paramList {
		value := rest[i]
		switch param {
		// byte
		case 'b':
			ob.codeList = append(ob.codeList, byte(value))
		// short(2byte int)
		case 's': 
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, uint16(value))
			ob.codeList = append(ob.codeList, b...)
		// constant pool index
		case 'p': 
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
	if ob.lineNumberList == nil || (ob.lineNumberList[len(ob.lineNumberList)-1].LineNumber != lineNumber) {
		l := &vm.VmLineNumber{
			LineNumber: lineNumber,
			StartPc:    start_pc,
			PcCount:    len(ob.codeList) - start_pc,
		}
		ob.lineNumberList = append(ob.lineNumberList, l)
	} else {
		ob.lineNumberList[len(ob.lineNumberList)-1].PcCount += len(ob.codeList) - start_pc
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

	dest := &vm.VmTypeSpecifier{BasicType: src.basicType}

	for _, derive := range src.deriveList {
		switch f := derive.(type) {
		case *FunctionDerive:
			newDerive := &vm.VmFunctionDerive{ParameterList: copyParameterList(f.parameterList)}
			dest.DeriveList = append(dest.DeriveList, newDerive)
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
			Name:          param.name,
			TypeSpecifier: copyTypeSpecifier(param.typeSpecifier),
		}
		dest = append(dest, v)
	}
	return dest
}

func copyFunction(src *FunctionDefinition, dest *vm.VmFunction) {
	dest.TypeSpecifier = copyTypeSpecifier(src.typeSpecifier)
	dest.Name = src.name
	dest.ParameterList = copyParameterList(src.parameterList)
	if src.block != nil {
		dest.LocalVariableList = copy_local_variables(src)
	} else {
		dest.LocalVariableList = nil
	}
}

func copy_local_variables(fd *FunctionDefinition) []*vm.VmLocalVariable {
	// TODO 形参占用位置
	var dest []*vm.VmLocalVariable = []*vm.VmLocalVariable{}

	localVariableCount := len(fd.localVariableList) - len(fd.parameterList)

	for _, v := range fd.localVariableList[0:localVariableCount] {
		vmV := &vm.VmLocalVariable{
			Name:          v.name,
			TypeSpecifier: copyTypeSpecifier(v.typeSpecifier),
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
		info := &vm.OpcodeInfo[ob.codeList[i]]
		for _, p := range []byte(info.Parameter) {
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
