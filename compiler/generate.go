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

func newCodeBuf() *OpcodeBuf {
	ob := &OpcodeBuf{
		codeList:       []byte{},
		labelTableList: []*LabelTable{},
		lineNumberList: []*vm.VmLineNumber{},
	}
	return ob
}

func (ob *OpcodeBuf) getLabel() int {
	// 返回栈顶位置
	ob.labelTableList = append(ob.labelTableList, &LabelTable{})
	return len(ob.labelTableList) - 1
}

func (ob *OpcodeBuf) setLabel(label int) {
	// 设置跳转
	ob.labelTableList[label].labelAddress = len(ob.codeList)
}

//
// generateCode
//
func (ob *OpcodeBuf) generateCode(pos Position, code byte, rest ...int) {
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
	ob.addLineNumber(pos.Line, startPc)
}

func (ob *OpcodeBuf) addLineNumber(lineNumber int, startPc int) {

	if len(ob.lineNumberList) == 0 || ob.lineNumberList[len(ob.lineNumberList)-1].LineNumber != lineNumber {
		newLineNumber := &vm.VmLineNumber{
			LineNumber: lineNumber,
			StartPc:    startPc,
			PcCount:    len(ob.codeList) - startPc,
		}
		ob.lineNumberList = append(ob.lineNumberList, newLineNumber)
	} else {
		// 源代码中相同的一行
		topLineNumber := ob.lineNumberList[len(ob.lineNumberList)-1]
		topLineNumber.PcCount += len(ob.codeList) - startPc
	}
}

//
// FIX
//
func (ob *OpcodeBuf) fixOpcodeBuf() []byte {

	ob.fixLabels()
	ob.labelTableList = nil

	return ob.codeList
}

// 修正label, 将正确的跳转地址填入
func (ob *OpcodeBuf) fixLabels() {

	for i := 0; i < len(ob.codeList); i++ {
		if ob.codeList[i] == vm.VM_JUMP ||
			ob.codeList[i] == vm.VM_JUMP_IF_TRUE ||
			ob.codeList[i] == vm.VM_JUMP_IF_FALSE {

			label := get2ByteInt(ob.codeList[i+1:])
			address := ob.labelTableList[label].labelAddress
			set2ByteInt(ob.codeList[i+1:], address)
		}
		info := &vm.OpcodeInfo[ob.codeList[i]]
		for _, p := range []byte(info.Parameter) {
			switch p {
			case 'b':
				i++
			case 's':
				fallthrough
			case 'p':
				i += 2
			default:
				panic("param error")
			}
		}
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
func copyTypeSpecifierNoAlloc(src *TypeSpecifier, dest *vm.VmTypeSpecifier) {

	dest.BasicType = src.basicType
	dest.DeriveList = []vm.VmTypeDerive{}

	for _, derive := range src.deriveList {
		switch realDerive := derive.(type) {
		case *FunctionDerive:
			newDerive := &vm.VmFunctionDerive{ParameterList: copyParameterList(realDerive.parameterList)}
			dest.AppendDerive(newDerive)
		case *ArrayDerive:
			dest.AppendDerive(&vm.VmArrayDerive{})
		default:
			panic("TODO")
		}
	}
}

func copyTypeSpecifier(src *TypeSpecifier) *vm.VmTypeSpecifier {

	dest := &vm.VmTypeSpecifier{}

	copyTypeSpecifierNoAlloc(src, dest)

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
		dest.LocalVariableList = copyLocalVariables(src)
	} else {
		dest.LocalVariableList = nil
	}
}

func copyLocalVariables(fd *FunctionDefinition) []*vm.VmLocalVariable {
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

// TODO 作为exe的方法
func AddTypeSpecifier(src *TypeSpecifier, exe *vm.Executable) int {
	ret := len(exe.TypeSpecifierList)

	newType := &vm.VmTypeSpecifier{}
	copyTypeSpecifierNoAlloc(src, newType)
	exe.TypeSpecifierList = append(exe.TypeSpecifierList, newType)

	return ret
}

func generate_pop_to_lvalue(exe *vm.Executable, block *Block , expr Expression , ob *OpcodeBuf) {
	identifierExpr, ok := expr.(*IdentifierExpression)
	if ok {
		generatePopToIdentifier(identifierExpr.inner.(*Declaration), expr.Position(), ob)
		return 
	}
	indexExpr, ok := expr.(*IndexExpression)
	if !ok {
		panic("TODO")
	}
	indexExpr.array.generate(exe, block, ob)
	indexExpr.index.generate(exe, block, ob)
        
	ob.generateCode(expr.Position(), vm.VM_POP_ARRAY_INT + getOpcodeTypeOffset(expr.typeS()))
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpcodeBuf) {
	var code byte

	offset := getOpcodeTypeOffset(decl.typeSpecifier)
	if decl.isLocal {
		code = vm.VM_POP_STACK_INT
	} else {
		code = vm.VM_POP_STATIC_INT
	}
	ob.generateCode(pos, code+offset, decl.variableIndex)
}
