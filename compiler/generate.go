package compiler

import (
	"encoding/binary"
	"log"

	"github.com/lth-go/gogogogo/vm"
)

type OpCodeBuf struct {
	codeList       []byte
	labelTableList []*LabelTable
	lineNumberList []*vm.LineNumber
}

type LabelTable struct {
	labelAddress int
}

func newCodeBuf() *OpCodeBuf {
	ob := &OpCodeBuf{
		codeList:       []byte{},
		labelTableList: []*LabelTable{},
		lineNumberList: []*vm.LineNumber{},
	}
	return ob
}

func (ob *OpCodeBuf) getLabel() int {
	// 返回栈顶位置
	ob.labelTableList = append(ob.labelTableList, &LabelTable{})
	return len(ob.labelTableList) - 1
}

func (ob *OpCodeBuf) setLabel(label int) {
	// 设置跳转
	ob.labelTableList[label].labelAddress = len(ob.codeList)
}

//
// generateCode
//
func (ob *OpCodeBuf) generateCode(pos Position, code byte, rest ...int) {
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

func (ob *OpCodeBuf) addLineNumber(lineNumber int, startPc int) {

	if len(ob.lineNumberList) == 0 || ob.lineNumberList[len(ob.lineNumberList)-1].LineNumber != lineNumber {
		newLineNumber := &vm.LineNumber{
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
func (ob *OpCodeBuf) fixOpcodeBuf() []byte {

	ob.fixLabels()
	ob.labelTableList = nil

	return ob.codeList
}

// 修正label, 将正确的跳转地址填入
func (ob *OpCodeBuf) fixLabels() {

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
			case 's', 'p':
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
func generateStatementList(exe *vm.Executable, currentBlock *Block, statementList []Statement, ob *OpCodeBuf) {
	for _, stmt := range statementList {
		stmt.generate(exe, currentBlock, ob)
	}
}

//
// COPY
//
func copyTypeSpecifierNoAlloc(src *TypeSpecifier, dest *vm.TypeSpecifier) {
	dest.BasicType = src.basicType
	if src.IsArray() {
		dest.SetSliceType(copyTypeSpecifier(src.sliceType.ElementType), src.sliceType.Len)
	}
}

func copyTypeSpecifier(src *TypeSpecifier) *vm.TypeSpecifier {
	dest := &vm.TypeSpecifier{}
	copyTypeSpecifierNoAlloc(src, dest)

	return dest
}

func copyParameterList(src []*Parameter) []*vm.LocalVariable {
	dest := []*vm.LocalVariable{}

	for _, param := range src {
		v := &vm.LocalVariable{
			Name:          param.name,
			TypeSpecifier: copyTypeSpecifier(param.typeSpecifier),
		}
		dest = append(dest, v)
	}
	return dest
}

func copyLocalVariables(fd *FunctionDefinition) []*vm.LocalVariable {
	// TODO 形参占用位置
	var dest = []*vm.LocalVariable{}

	localVariableCount := len(fd.localVariableList) - len(fd.parameterList)

	for _, v := range fd.localVariableList[0:localVariableCount] {
		vmV := &vm.LocalVariable{
			Name:          v.name,
			TypeSpecifier: copyTypeSpecifier(v.typeSpecifier),
		}
		dest = append(dest, vmV)
	}

	return dest
}

func generatePopToLvalue(exe *vm.Executable, block *Block, expr Expression, ob *OpCodeBuf) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		generatePopToIdentifier(e.inner.(*Declaration), expr.Position(), ob)
	case *IndexExpression:
		e.array.generate(exe, block, ob)
		e.index.generate(exe, block, ob)
		ob.generateCode(expr.Position(), vm.VM_POP_ARRAY_INT+getOpcodeTypeOffset(expr.typeS()))
	}
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpCodeBuf) {
	var code byte

	offset := getOpcodeTypeOffset(decl.typeSpecifier)
	if decl.isLocal {
		code = vm.VM_POP_STACK_INT
	} else {
		code = vm.VM_POP_STATIC_INT
	}
	ob.generateCode(pos, code+offset, decl.variableIndex)
}

func generatePushArgument(argList []Expression, exe *vm.Executable, currentBlock *Block, ob *OpCodeBuf) {
	for _, arg := range argList {
		arg.generate(exe, currentBlock, ob)
	}
}

func getOpcodeTypeOffset(typ *TypeSpecifier) byte {
	if typ.IsComposite() {
		return byte(2)
	}

	switch {
	case typ.IsVoid():
		panic("basic type is void")
	case typ.IsBool(), typ.IsInt():
		return byte(0)
	case typ.IsFloat():
		return byte(1)
	case typ.IsString():
		return byte(2)
	case typ.IsNil(), typ.IsBase():
		fallthrough
	default:
		log.Fatalf("TODO")
	}

	return byte(0)
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}

func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}
