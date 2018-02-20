package vm

//
// Static
//
// 虚拟机全局静态变量
type Static struct {
	variableList []VmValue
}

func NewStatic() *Static {
	s := &Static{
		variableList: []VmValue{},
	}
	return s
}

func (s *Static) append(value VmValue) {
	s.variableList = append(s.variableList, value)
}

//
// get
//
func (s *Static) getInt(index int) int {
	return s.variableList[index].(*VmIntValue).intValue
}

func (s *Static) getDouble(index int) float64 {
	return s.variableList[index].(*VmDoubleValue).doubleValue
}

func (s *Static) getObject(index int) VmObject {
	return s.variableList[index].(*VmObjectValue).objectValue
}

//
// set
//
func (s *Static) setInt(index int, value int) {
	s.variableList[index].(*VmIntValue).intValue = value
}

func (s *Static) setDouble(index int, value float64) {
	s.variableList[index].(*VmDoubleValue).doubleValue = value
}

func (s *Static) setObject(index int, value VmObject) {
	s.variableList[index].(*VmObjectValue).objectValue = value
}
