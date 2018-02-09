package vm

//
// Stack
//
// 虚拟机栈
type Stack struct {
	stackPointer int
	stack        []VmValue
}

func NewStack() Stack {
	s := Stack{
		stack: []VmValue{},
		stackPointer: 0,
	}
	return s
}
