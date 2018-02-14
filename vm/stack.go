package vm

const (
	stackAllocSize int = 4096
)

//
// Stack
//
// 虚拟机栈
type Stack struct {
	stackPointer int
	stack        []VmValue
}

func NewStack() *Stack {
	s := &Stack{
		stack:        make([]VmValue, stackAllocSize, (stackAllocSize+1)*2),
		stackPointer: 0,
	}
	return s
}

// expand
func (s *Stack) expand(codeList []byte) {
	needStackSize := calcNeedStackSize(codeList)

	rest := len(s.stack) - s.stackPointer

	if rest <= needStackSize {
		size := len(s.stack) + needStackSize - rest

		newStack := make([]VmValue, size, (size+1)*2)
		copy(newStack, s.stack)

		s.stack = newStack
	}
}

func calcNeedStackSize(codeList []byte) int {
	stackSize := 0

	for i := 0; i < len(codeList); i++ {
		info := OpcodeInfo[int(codeList[i])]
		if info.stackIncrement > 0 {
			stackSize += info.stackIncrement
		}
		for _, p := range []byte(info.Parameter) {
			switch p {
			case 'b':
				i++
			case 's', 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
	}

	return stackSize
}

// 根据sp以及stackPointer返回栈的位置
func (s *Stack) getIndexOverSp(sp int) int {
	index := s.stackPointer + sp
	if index == -1 {
		index = len(s.stack) -1
	}
	return index
}

// 根据sp以及stackPointer返回栈中元素
func (s *Stack) getInt(sp int) int {
	index := s.getIndexOverSp(sp)
	return s.getIntI(index)
}

func (s *Stack) getDouble(sp int) float64 {
	index := s.getIndexOverSp(sp)
	return s.getDoubleI(index)
}

func (s *Stack) getObject(sp int) VmObject {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index)
}

// 直据sp返回栈中元素
func (s *Stack) getIntI(sp int) int {
	return s.stack[sp].getIntValue()
}
func (s *Stack) getDoubleI(sp int) float64 {
	return s.stack[sp].getDoubleValue()
}
func (s *Stack) getObjectI(sp int) VmObject {
	return s.stack[sp].getObjectValue()
}

// 根据sp以及stackPointer向栈中写入元素
func (s *Stack) writeInt(sp int, value int) {
	index := s.getIndexOverSp(sp)
	s.writeIntI(index, value)
}
func (s *Stack) writeDouble(sp int, value float64) {
	index := s.getIndexOverSp(sp)
	s.writeDoubleI(index, value)
}
func (s *Stack) writeObject(sp int, value VmObject) {
	index := s.getIndexOverSp(sp)
	s.writeObjectI(index, value)
}

// 根据sp向栈中写入元素
func (s *Stack) writeIntI(sp int, value int) {
	v := NewIntValue(value)
	v.setPointer(false)

	s.stack[sp] = v
}
func (s *Stack) writeDoubleI(sp int, value float64) {
	v := NewDoubleValue(value)
	v.setPointer(false)

	s.stack[sp] = v
}
func (s *Stack) writeObjectI(sp int, value VmObject) {
	v := NewObjectValue(value)
	v.setPointer(true)

	s.stack[sp] = v
}
