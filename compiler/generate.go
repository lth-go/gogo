package compiler

import (
	"encoding/binary"
	"log"

	"github.com/lth-go/gogo/vm"
)

type OpCodeBuf struct {
	codeList       []byte
	labelTableList []*LabelTable
	lineNumberList []*vm.LineNumber
}

type LabelTable struct {
	labelAddress int
}

func NewOpCodeBuf() *OpCodeBuf {
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
func generateStatementList(statementList []Statement, ob *OpCodeBuf) {
	for _, stmt := range statementList {
		stmt.Generate(ob)
	}
}

//
// COPY
//
func CopyToVmType(src *Type) *vm.Type {
	dest := &vm.Type{
		BasicType: src.GetBasicType(),
	}
	if src.IsArray() {
		dest.SetSliceType(CopyToVmType(src.arrayType.ElementType), src.arrayType.Len)
	}

	if src.IsFunc() {
		paramTypeList := []*vm.Type{}
		resultTypeList := []*vm.Type{}
		for _, t := range src.funcType.Params {
			paramTypeList = append(paramTypeList, CopyToVmType(t.Type))
		}

		for _, t := range src.funcType.Results {
			resultTypeList = append(resultTypeList, CopyToVmType(t.Type))
		}

		dest.SetFuncType(paramTypeList, resultTypeList)
	}

	return dest
}

func copyVmVariableList(fd *FunctionDefinition) []*vm.Variable {
	// TODO 形参占用位置
	var dest = []*vm.Variable{}

	localVariableCount := len(fd.DeclarationList) - len(fd.ParameterList)

	for _, v := range fd.DeclarationList[0:localVariableCount] {
		vmV := &vm.Variable{
			Name: v.Name,
			Type: CopyToVmType(v.Type),
		}
		dest = append(dest, vmV)
	}

	return dest
}

func generatePopToLvalue(expr Expression, ob *OpCodeBuf) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		generatePopToIdentifier(e.inner.(*Declaration), expr.Position(), ob)
	case *IndexExpression:
		e.array.generate(ob)
		e.index.generate(ob)
		ob.generateCode(expr.Position(), vm.VM_POP_ARRAY_OBJECT)
	default:
		panic("TODO")
	}
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpCodeBuf) {
	var code byte

	offset := getOpcodeTypeOffset(decl.Type)
	if decl.IsLocal {
		code = vm.VM_POP_STACK_INT
	} else {
		code = vm.VM_POP_STATIC_INT
	}
	ob.generateCode(pos, code+offset, decl.Index)
}

func generatePushArgument(argList []Expression, ob *OpCodeBuf) {
	for _, arg := range argList {
		arg.generate(ob)
	}
}

func getOpcodeTypeOffset(typ *Type) byte {
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
	case typ.IsNil():
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
