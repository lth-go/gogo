package compiler

//
// compare
//
func compareType(typ1 *Type, typ2 *Type) bool {
	if typ1.GetBasicType() != typ2.GetBasicType() {
		return false
	}

	t1 := typ1.sliceType
	t2 := typ2.sliceType

	if t1 == nil && t2 == nil {
		return true
	}

	if t1 == nil && t2 != nil {
		return false
	}

	if t1 != nil && t2 == nil {
		return false
	}

	return compareType(t1.ElementType, t2.ElementType)
}
