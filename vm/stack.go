package vm

const (
	stackAllocSize = 4096
)

// 虚拟机栈
type Stack struct {
	stackPointer int
	stack        []Value
	objectList   []Object
}

func NewStack() *Stack {
	s := &Stack{
		stackPointer: 0,
		stack:        make([]Value, stackAllocSize, (stackAllocSize+1)*2),
		objectList:   make([]Object, stackAllocSize, (stackAllocSize+1)*2),
	}
	return s
}

// 栈伸缩
func (s *Stack) Expand(codeList []byte) {
	needStackSize := getNeedStackSize(codeList)

	rest := s.Len() - s.stackPointer

	if rest <= needStackSize {
		size := s.Len() + needStackSize - rest

		// TODO: remove
		newStack := make([]Value, size, (size+1)*2)
		copy(newStack, s.stack)
		s.stack = newStack

		newObjectList := make([]Object, size, (size+1)*2)
		copy(newObjectList, s.objectList)
		s.objectList = newObjectList
	}
}

func getNeedStackSize(codeList []byte) int {
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

func (s *Stack) Len() int {
	return len(s.stack)
}

// 根据incr以及stackPointer返回栈的位置
func (s *Stack) getIndex(incr int) int {
	index := s.stackPointer + incr
	if index == -1 {
		index = len(s.stack) - 1
	}
	return index
}

func (s *Stack) Get(sp int) Value {
	return s.stack[sp]
}

func (s *Stack) Set(sp int, v Value) {
	switch value := v.(type) {
	case *IntValue:
		s.objectList[sp] = NewObjectInt(value.intValue)
	case *DoubleValue:
		s.objectList[sp] = NewObjectFloat(value.doubleValue)
	case *ObjectRef:
		s.objectList[sp] = value.data
	}

	s.stack[sp] = v
}

// 直据sp返回栈中元素
func (s *Stack) GetInt(sp int) int {
	return s.objectList[sp].(*ObjectInt).Value
}

func (s *Stack) GetFloat(sp int) float64 {
	return s.objectList[sp].(*ObjectFloat).Value
}

func (s *Stack) GetString(sp int) string {
	return s.objectList[sp].(*ObjectString).Value
}

func (s *Stack) GetObject(sp int) *ObjectRef {
	return &ObjectRef{
		ValueImpl: ValueImpl{
			pointerFlags: true,
		},
		data: s.objectList[sp],
	}
}

// 根据incr以及stackPointer返回栈中元素
func (s *Stack) GetIntPlus(incr int) int {
	index := s.getIndex(incr)
	return s.GetInt(index)
}

func (s *Stack) GetFloatPlus(incr int) float64 {
	index := s.getIndex(incr)
	return s.GetFloat(index)
}

func (s *Stack) GetStringPlus(incr int) string {
	index := s.getIndex(incr)
	return s.GetString(index)
}

func (s *Stack) GetObjectPlus(incr int) *ObjectRef {
	index := s.getIndex(incr)
	return s.GetObject(index)
}

// 根据sp向栈中写入元素
func (s *Stack) SetInt(sp int, value int) {
	v := NewIntValue(value)
	v.setPointer(false)
	s.Set(sp, v)
}

func (s *Stack) SetFloat(sp int, value float64) {
	v := NewDoubleValue(value)
	v.setPointer(false)
	s.Set(sp, v)
}

func (s *Stack) SetString(sp int, value string) {
	v := NewObjectString(value)
	s.objectList[sp] = v
}

func (s *Stack) SetObject(sp int, value *ObjectRef) {
	value.setPointer(true)
	s.Set(sp, value)
}

// 根据incr以及stackPointer向栈中写入元素
func (s *Stack) SetIntPlus(incr int, value int) {
	index := s.getIndex(incr)
	s.SetInt(index, value)
}

func (s *Stack) SetDoublePlus(incr int, value float64) {
	index := s.getIndex(incr)
	s.SetFloat(index, value)
}

func (s *Stack) SetStringPlus(incr int, value string) {
	index := s.getIndex(incr)
	s.SetString(index, value)
}

func (s *Stack) SetObjectPlus(incr int, value *ObjectRef) {
	index := s.getIndex(incr)
	s.SetObject(index, value)
}

func (s *Stack) getArrayInt(sp int) *ObjectArrayInt {
	index := s.getIndex(sp)
	return s.GetObject(index).data.(*ObjectArrayInt)
}

func (s *Stack) getArrayDouble(sp int) *ObjectArrayDouble {
	index := s.getIndex(sp)
	return s.GetObject(index).data.(*ObjectArrayDouble)
}

func (s *Stack) getArrayObject(sp int) *ObjectArrayObject {
	index := s.getIndex(sp)
	return s.GetObject(index).data.(*ObjectArrayObject)
}
