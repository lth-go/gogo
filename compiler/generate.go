package compiler

import (
	"encoding/binary"

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
	paramList := []byte(vm.OpcodeInfo[code].Parameter)

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
		if ob.codeList[i] == vm.OP_CODE_JUMP ||
			ob.codeList[i] == vm.OP_CODE_JUMP_IF_TRUE ||
			ob.codeList[i] == vm.OP_CODE_JUMP_IF_FALSE {

			label := get2ByteInt(ob.codeList[i+1:])
			address := ob.labelTableList[label].labelAddress
			set2ByteInt(ob.codeList[i+1:], address)
		}
		info := vm.OpcodeInfo[ob.codeList[i]]
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
		dest.SetArrayType(CopyToVmType(src.arrayType.ElementType), src.arrayType.Len)
	}

	if src.IsStruct() {
		typeList := make([]*vm.Type, 0)

		for _, field := range src.structType.Fields {
			typeList = append(typeList, CopyToVmType(field.Type))
		}

		dest.SetStructType(typeList)
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

func copyVmVariableList(fd *FunctionDefinition) []vm.Object {
	// TODO 形参占用位置
	var dest = []vm.Object{}

	localVariableCount := len(fd.DeclarationList) - len(fd.ParameterList)

	for _, v := range fd.DeclarationList[0:localVariableCount] {
		dest = append(dest, vm.GetObjectByType(CopyToVmType(v.Type)))
	}

	return dest
}

func generatePopToLvalue(expr Expression, ob *OpCodeBuf) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		generatePopToIdentifier(e.inner.(*Declaration), expr.Position(), ob)
	case *IndexExpression:
		if e.X.GetType().IsArray() {
			e.X.Generate(ob)
			e.Index.Generate(ob)
			ob.generateCode(expr.Position(), vm.OP_CODE_POP_ARRAY)
		} else if e.X.GetType().IsMap() {
			e.X.Generate(ob)
			e.Index.Generate(ob)
			ob.generateCode(expr.Position(), vm.OP_CODE_POP_MAP)
		} else {
			panic("TODO")
		}
	case *SelectorExpression:
		if e.expression.GetType().IsStruct() {
			e.expression.Generate(ob)
			ob.generateCode(expr.Position(), vm.OP_CODE_PUSH_INT_2BYTE, e.Index)
			ob.generateCode(expr.Position(), vm.OP_CODE_POP_STRUCT)
		} else {
			panic("TODO")
		}

	default:
		panic("TODO")
	}
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpCodeBuf) {
	var code byte

	if decl.IsLocal {
		code = vm.OP_CODE_POP_STACK
	} else {
		code = vm.OP_CODE_POP_STATIC
	}
	ob.generateCode(pos, code, decl.Index)
}

func generatePushArgument(argList []Expression, ob *OpCodeBuf) {
	for _, arg := range argList {
		arg.Generate(ob)
	}
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}

func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}
