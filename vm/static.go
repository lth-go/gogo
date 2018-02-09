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
