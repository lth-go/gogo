package vm

type SliceType struct {
	Len         int64
	ElementType *Type
}

type Type struct {
	BasicType  BasicType
	SliceType  *SliceType
}

func (t *Type) SetSliceType(typ *Type, length int64) {
	t.SliceType = &SliceType{
		Len:         length,
		ElementType: typ,
	}
}

func (t *Type) IsSliceType() bool {
	// TODO: 根据basicType判断
	return t.SliceType != nil
}

func (t *Type) IsReferenceType() bool {
	if t.BasicType == BasicTypeString || t.IsSliceType() {
		return true
	}
	return false
}
