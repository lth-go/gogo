package vm

type SliceType struct {
	Len         int64
	ElementType *Type
}

type FuncType struct {
	ParamTypeList  []*Type
	ResultTypeList []*Type
}

type Type struct {
	BasicType BasicType
	SliceType *SliceType
	FuncType  *FuncType
}

func (t *Type) SetSliceType(typ *Type, length int64) {
	t.SliceType = &SliceType{
		Len:         length,
		ElementType: typ,
	}
}

func (t *Type) SetFuncType(paramsTypeList []*Type, resultTypeList []*Type) {
	t.FuncType = &FuncType{
		ParamTypeList:  paramsTypeList,
		ResultTypeList: resultTypeList,
	}
}

func (t *Type) IsSliceType() bool {
	return t.BasicType == BasicTypeSlice
}
