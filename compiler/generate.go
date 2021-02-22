package compiler

import (
	"encoding/binary"

	"github.com/lth-go/gogo/utils"
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

func (ob *OpCodeBuf) GetLabel() int {
	// 返回栈顶位置
	ob.labelTableList = append(ob.labelTableList, &LabelTable{})
	return len(ob.labelTableList) - 1
}

func (ob *OpCodeBuf) SetLabel(label int) {
	// 设置跳转
	ob.labelTableList[label].labelAddress = len(ob.codeList)
}

//
// GenerateCode
//
func (ob *OpCodeBuf) GenerateCode(pos Position, code byte, rest ...int) {
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
	ob.AddLineNumber(pos.Line, startPc)
}

func (ob *OpCodeBuf) AddLineNumber(lineNumber int, startPc int) {

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

// 修正label, 将正确的跳转地址填入
func (ob *OpCodeBuf) FixLabel() []byte {
	for i := 0; i < len(ob.codeList); i++ {
		if ob.codeList[i] == vm.OP_CODE_JUMP ||
			ob.codeList[i] == vm.OP_CODE_JUMP_IF_TRUE ||
			ob.codeList[i] == vm.OP_CODE_JUMP_IF_FALSE {

			label := utils.Get2ByteInt(ob.codeList[i+1:])
			address := ob.labelTableList[label].labelAddress
			utils.Set2ByteInt(ob.codeList[i+1:], address)
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

	ob.labelTableList = nil

	return ob.codeList
}
