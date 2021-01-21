package vm

//
// 基本类型
//
type BasicType int

const (
	BasicTypeNoType BasicType = iota - 1
	BasicTypeBool
	BasicTypeInt
	BasicTypeFloat
	BasicTypeString
	BasicTypeNil
	BasicTypeVoid
	BasicTypePackage
	BasicTypeSlice
	BasicTypeMap
	BasicTypeStruct
	BasicTypeFunc
)

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

func (t *Type) GetBasicType() BasicType {
	return t.BasicType
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
