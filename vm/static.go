package vm

// 虚拟机全局静态变量
type Static struct {
	variableList     []Value
}

func NewStatic() *Static {
	s := &Static{
		variableList: []Value{},
	}
	return s
}

func (s *Static) append(value Value) {
	s.variableList = append(s.variableList, value)
}

//
// get
//
func (s *Static) getInt(index int) int {
	return s.variableList[index].(*IntValue).intValue
}

func (s *Static) getDouble(index int) float64 {
	return s.variableList[index].(*DoubleValue).doubleValue
}

func (s *Static) getObject(index int) *ObjectRef {
	return s.variableList[index].(*ObjectRef)
}

//
// set
//
func (s *Static) setInt(index int, value int) {
	s.variableList[index].(*IntValue).intValue = value
}

func (s *Static) setDouble(index int, value float64) {
	s.variableList[index].(*DoubleValue).doubleValue = value
}

func (s *Static) setObject(index int, value *ObjectRef) {
	s.variableList[index] = value
}
