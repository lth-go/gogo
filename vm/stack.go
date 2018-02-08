package vm

//
// Stack
//
// 虚拟机栈
type Stack struct {
	stackPointer int
	stack        []VmValue
}
