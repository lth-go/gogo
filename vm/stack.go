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
		stack:        make([]VmValue, stackAllocSize, (stackAllocSize +1)*2),
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

func (s *Stack) getInt(sp int) int {
	index := s.stackPointer + sp
	if index == -1 {
		index = len(s.stack) - 1
	}
	return s.stack[index].getIntValue()
}
func (s *Stack) setInt(sp int, v int) {
	s.stack[s.stackPointer+sp].setIntValue(v)
}
func (s *Stack) getDouble(sp int) float64 {
	index := s.stackPointer + sp
	if index == -1 {
		index = len(s.stack) - 1
	}
	return s.stack[index].getDoubleValue()
}
func (s *Stack) setDouble(sp int, v float64) {
	s.stack[s.stackPointer+sp].setDoubleValue(v)
}
func (s *Stack) getObject(sp int) VmObject {
	index := s.stackPointer + sp
	if index == -1 {
		index = len(s.stack) - 1
	}
	return s.stack[index].getObjectValue()
}
func (s *Stack) setObject(sp int, v VmObject) {
	s.stack[s.stackPointer+sp].setObjectValue(v)
}

func (s *Stack) getIntI(sp int) int {
	return s.stack[sp].getIntValue()
}
func (s *Stack) getDoubleI(sp int) float64 {
	return s.stack[sp].getDoubleValue()
}
func (s *Stack) getObjectI(sp int) VmObject {
	return s.stack[sp].getObjectValue()
}

func (s *Stack) writeInt(sp int, r int) {
	v := &VmIntValue{
		intValue: r,
	}
	s.stack[s.stackPointer+sp] = v

	v.setPointer(false)
}
func (s *Stack) writeDouble(sp int, r float64) {
	v := &VmDoubleValue{
		doubleValue: r,
	}
	s.stack[s.stackPointer+sp] = v

	v.setPointer(false)
}
func (s *Stack) writeObject(sp int, r VmObject) {
	v := &VmObjectValue{
		objectValue: r,
	}
	s.stack[s.stackPointer+sp] = v

	v.setPointer(true)
}

func (s *Stack) writeIntI(sp int, r int) {
	v := s.stack[sp]
	v.setIntValue(r)
	v.setPointer(false)
}
func (s *Stack) writeDoubleI(sp int, r float64) {
	v := s.stack[sp]
	v.setDoubleValue(r)
	v.setPointer(false)
}
func (s *Stack) writeObjectI(sp int, r VmObject) {
	v := s.stack[sp]
	v.setObjectValue(r)
	v.setPointer(true)
}
