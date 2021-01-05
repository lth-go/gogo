package vm

type TypeDerive interface{}

type FunctionDerive struct {
	ParameterList []*LocalVariable
}

type ArrayDerive struct {
}

type TypeSpecifier struct {
	BasicType  BasicType
	DeriveType TypeDerive
}

func (t *TypeSpecifier) SetDeriveType(derive TypeDerive) {
	t.DeriveType = derive
}

func (t *TypeSpecifier) isArrayDerive() bool {
	_, ok := t.DeriveType.(*ArrayDerive)
	return ok
}
