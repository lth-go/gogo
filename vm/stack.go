package vm

const (
	stackAllocSize = 4096
)

//
// Stack
//
// 虚拟机栈
type Stack struct {
	stackPointer int
	stack        []Value
}

func NewStack() *Stack {
	s := &Stack{
		stack:        make([]Value, stackAllocSize, (stackAllocSize+1)*2),
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

		newStack := make([]Value, size, (size+1)*2)
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

func (s *Stack) getObject(sp int) *ObjectRef {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index)
}

// 直据sp返回栈中元素
func (s *Stack) getIntI(sp int) int {
	value := s.stack[sp].(*IntValue)
	return value.intValue
}
func (s *Stack) getDoubleI(sp int) float64 {
	value := s.stack[sp].(*DoubleValue)
	return value.doubleValue
}
func (s *Stack) getObjectI(sp int) *ObjectRef {
	value := s.stack[sp].(*ObjectRef)
	return value
}

// 根据sp以及stackPointer向栈中写入元素
func (s *Stack) setInt(sp int, value int) {
	index := s.getIndexOverSp(sp)
	s.setIntI(index, value)
}
func (s *Stack) setDouble(sp int, value float64) {
	index := s.getIndexOverSp(sp)
	s.setDoubleI(index, value)
}
func (s *Stack) setObject(sp int, value *ObjectRef) {
	index := s.getIndexOverSp(sp)
	s.setObjectI(index, value)
}

// 根据sp向栈中写入元素
func (s *Stack) setIntI(sp int, value int) {
	v := NewIntValue(value)
	v.setPointer(false)

	s.stack[sp] = v
}
func (s *Stack) setDoubleI(sp int, value float64) {
	v := NewDoubleValue(value)
	v.setPointer(false)

	s.stack[sp] = v
}
func (s *Stack) setObjectI(sp int, value *ObjectRef) {
	v := value
	v.setPointer(true)

	s.stack[sp] = v
}

// other get
func (s *Stack) getString(sp int) string {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index).data.(*ObjectString).stringValue
}

func (s *Stack) getArrayInt(sp int) *ObjectArrayInt {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index).data.(*ObjectArrayInt)
}

func (s *Stack) getArrayDouble(sp int) *ObjectArrayDouble {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index).data.(*ObjectArrayDouble)
}

func (s *Stack) getArrayObject(sp int) *ObjectArrayObject {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index).data.(*ObjectArrayObject)
}

func (s *Stack) getClassObject(sp int) *ObjectClassObject {
	index := s.getIndexOverSp(sp)
	return s.getObjectI(index).data.(*ObjectClassObject)
}
