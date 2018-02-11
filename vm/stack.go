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
func (s *Stack) getInt(sp int) int {
	return s.stack[s.stackPointer+sp].getIntValue()
}
func (s *Stack) setInt(sp int, v int) {
	s.stack[s.stackPointer+sp].setIntValue(v)
}
func (s *Stack) getDouble(sp int) float64 {
	return s.stack[s.stackPointer+sp].getDoubleValue()
}
func (s *Stack) setDouble(sp int, v float64) {
	s.stack[s.stackPointer+sp].setDoubleValue(v)
}
func (s *Stack) getObject(sp int) VmObject {
	return s.stack[s.stackPointer+sp].getObjectValue()
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
	v := s.stack[s.stackPointer+sp]
	v.setIntValue(r)
	v.setPointer(false)
}
func (s *Stack) writeDouble(sp int, r float64) {
	v := s.stack[s.stackPointer+sp]
	v.setDoubleValue(r)
	v.setPointer(false)
}
func (s *Stack) writeObject(sp int, r VmObject) {
	v := s.stack[s.stackPointer+sp]
	v.setObjectValue(r)
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
