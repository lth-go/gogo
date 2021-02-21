package vm

//
// Static
//
type Static struct {
	list []Object
}

func NewStatic() *Static {
	return &Static{
		list: make([]Object, 0),
	}
}

func (s *Static) Append(obj Object) {
	s.list = append(s.list, obj)
}

func (s *Static) Get(index int) Object {
	return s.list[index]
}

func (s *Static) Set(index int, obj Object) {
	s.list[index] = obj
}
