package vm

type TypeDerive interface{}

type FunctionDerive struct {
	ParameterList []*LocalVariable
}

type ArrayDerive struct {
}

type SliceType struct {
	Len         int64
	ElementType *TypeSpecifier
}

type TypeSpecifier struct {
	BasicType  BasicType
	DeriveType TypeDerive // TODO: remove
	SliceType  *SliceType
}

func (t *TypeSpecifier) SetDeriveType(derive TypeDerive) {
	t.DeriveType = derive
}

func (t *TypeSpecifier) SetSliceType(typ *TypeSpecifier, length int64) {
	t.SliceType = &SliceType{
		Len:         length,
		ElementType: typ,
	}
}

func (t *TypeSpecifier) isArrayDerive() bool {
	_, ok := t.DeriveType.(*ArrayDerive)
	return ok
}

func (t *TypeSpecifier) IsReferenceType() bool {
	if ((t.BasicType == StringType) && t.DeriveType == nil) || (t.isArrayDerive()) {
		return true
	}
	return false
}
