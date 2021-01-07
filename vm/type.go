package vm

type SliceType struct {
	Len         int64
	ElementType *TypeSpecifier
}

type TypeSpecifier struct {
	BasicType  BasicType
	SliceType  *SliceType
}

func (t *TypeSpecifier) SetSliceType(typ *TypeSpecifier, length int64) {
	t.SliceType = &SliceType{
		Len:         length,
		ElementType: typ,
	}
}

func (t *TypeSpecifier) IsSliceType() bool {
	// TODO: 根据basicType判断
	return t.SliceType != nil
}

func (t *TypeSpecifier) IsReferenceType() bool {
	if t.BasicType == BasicTypeString || t.IsSliceType() {
		return true
	}
	return false
}
