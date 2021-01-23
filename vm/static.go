package vm

// StaticValue function, variable
type StaticValue interface {
	GetName() string
	GetPackageName() string
}

type StaticBase struct {
	Name        string
	PackageName string
}

func (f *StaticBase) GetName() string        { return f.Name }
func (f *StaticBase) GetPackageName() string { return f.PackageName }

// Static
type Static struct {
	list []StaticValue
}

func NewStatic() *Static {
	return &Static{
		list: make([]StaticValue, 0),
	}
}

func (s *Static) Append(v StaticValue) {
	s.list = append(s.list, v)
}

func (s *Static) Index(packageName string, name string) int {
	for i, v := range s.list {
		if v.GetPackageName() == packageName && v.GetName() == name {
			return i
		}
	}

	return -1
}

func (s *Static) Get(index int) StaticValue {
	return s.list[index]
}

func (s *Static) GetVariableObject(index int) Object {
	return s.list[index].(*StaticVariable).Value.(Object)
}

func (s *Static) SetVariable(index int, value interface{}) {
	s.list[index].(*StaticVariable).Value = value
}

type StaticVariable struct {
	StaticBase
	Value interface{} // function or object
}

func NewStaticVariable(packageName string, name string, value interface{}) *StaticVariable {
	return &StaticVariable{
		StaticBase: StaticBase{
			PackageName: packageName,
			Name:        name,
		},
		Value: value,
	}
}
