package vm

const (
	stackAllocSize = 4096
)

// 虚拟机栈
type Stack struct {
	stackPointer int      // 栈偏移量, 指向当前最大空栈
	objectList   []Object // 对象栈
}

func NewStack() *Stack {
	s := &Stack{
		stackPointer: 0,
		objectList:   make([]Object, stackAllocSize, (stackAllocSize+1)*2),
	}
	return s
}

//
// 栈伸缩
//
func (s *Stack) Expand(codeList []byte) {
	needStackSize := getNeedStackSize(codeList)

	rest := s.Len() - s.stackPointer

	if rest <= needStackSize {
		size := s.Len() + needStackSize - rest

		newObjectList := make([]Object, size, (size+1)*2)
		copy(newObjectList, s.objectList)
		s.objectList = newObjectList
	}
}

func getNeedStackSize(codeList []byte) int {
	stackSize := 0

	for i := 0; i < len(codeList); i++ {
		info := OpcodeInfo[codeList[i]]
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
	return len(s.objectList)
}

// 根据incr以及stackPointer返回栈的位置
func (s *Stack) getIndex(incr int) int {
	index := s.stackPointer + incr
	if index == -1 {
		index = s.Len() - 1
	}
	return index
}

// 直据sp返回栈中元素
func (s *Stack) Get(sp int) Object {
	return s.objectList[sp]
}

// 根据incr以及stackPointer向栈中写入元素
func (s *Stack) GetPlus(incr int) Object {
	index := s.getIndex(incr)
	return s.Get(index)
}

func (s *Stack) Set(sp int, v Object) {
	s.objectList[sp] = v
}

func (s *Stack) SetPlus(incr int, value Object) {
	index := s.getIndex(incr)
	s.Set(index, value)
}

func (s *Stack) GetInt(sp int) int {
	return s.Get(sp).(*ObjectInt).Value
}

func (s *Stack) GetFloat(sp int) float64 {
	return s.Get(sp).(*ObjectFloat).Value
}

func (s *Stack) GetString(sp int) string {
	return s.Get(sp).(*ObjectString).Value
}

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

func (s *Stack) GetArrayPlus(incr int) *ObjectArray {
	index := s.getIndex(incr)
	return s.Get(index).(*ObjectArray)
}

func (s *Stack) GetMapPlus(incr int) *ObjectMap {
	index := s.getIndex(incr)
	return s.Get(index).(*ObjectMap)
}

func (s *Stack) GetStructPlus(incr int) *ObjectStruct {
	index := s.getIndex(incr)
	return s.Get(index).(*ObjectStruct)
}

func (s *Stack) SetInt(sp int, value int) {
	s.Set(sp, NewObjectInt(value))
}

func (s *Stack) SetFloat(sp int, value float64) {
	s.Set(sp, NewObjectFloat(value))
}

func (s *Stack) SetString(sp int, value string) {
	s.Set(sp, NewObjectString(value))
}

func (s *Stack) SetIntPlus(incr int, value int) {
	index := s.getIndex(incr)
	s.SetInt(index, value)
}

func (s *Stack) SetFloatPlus(incr int, value float64) {
	index := s.getIndex(incr)
	s.SetFloat(index, value)
}

func (s *Stack) SetStringPlus(incr int, value string) {
	index := s.getIndex(incr)
	s.SetString(index, value)
}
