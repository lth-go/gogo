package vm

//
// Static
//
// 虚拟机全局静态变量
type Static struct {
	variableList []VmValue
}

func NewStatic() Static {
	s := Static{
		variableList: []VmValue{},
	}
	return s
}

func (s *Static) append(value VmValue) {
	s.variableList = append(s.variableList, value)
}

func (s *Static) getInt(index int) int {
	return s.variableList[index].getIntValue()
}

func (s *Static) getDouble(index int) float64 {
	return s.variableList[index].getDoubleValue()
}

func (s *Static) getObject(index int) VmObject {
	return s.variableList[index].getObjectValue()
}

func (s *Static) setInt(index int, value int) {
	s.variableList[index].setIntValue(value)
}

func (s *Static) setDouble(index int, value float64) {
	s.variableList[index].setDoubleValue(value)
}

func (s *Static) setObject(index int, value VmObject) {
	s.variableList[index].setObjectValue(value)
}
