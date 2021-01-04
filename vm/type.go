package vm

type TypeDerive interface{}

type FunctionDerive struct {
	ParameterList []*LocalVariable
}

type ArrayDerive struct {
}

type TypeSpecifier struct {
	BasicType  BasicType
	DeriveList []TypeDerive
}

func (t *TypeSpecifier) AppendDerive(derive TypeDerive) {
	if t.DeriveList == nil {
		t.DeriveList = []TypeDerive{}
	}
	t.DeriveList = append(t.DeriveList, derive)
}

func (t *TypeSpecifier) isArrayDerive() bool {
	return isArray(t)
}

func isArray(t *TypeSpecifier) bool {
	if t.DeriveList == nil || len(t.DeriveList) == 0 {
		return false
	}
	firstElem := t.DeriveList[0]
	_, ok := firstElem.(*ArrayDerive)
	return ok
}
