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
	BasicTypeArray
	BasicTypeMap
	BasicTypeStruct
	BasicTypeFunc
	BasicTypeMultipleValues
	BasicTypeInterface
	BasicTypePointer
)

type ArrayType struct {
	Len         int64
	ElementType *Type
}

type FuncType struct {
	ParamTypeList  []*Type
	ResultTypeList []*Type
}

type StructType struct {
	FieldTypeList []*Type
}

type Type struct {
	BasicType  BasicType
	ArrayType  *ArrayType
	FuncType   *FuncType
	StructType *StructType
}

func (t *Type) GetBasicType() BasicType {
	return t.BasicType
}

func (t *Type) SetArrayType(typ *Type, length int64) {
	t.ArrayType = &ArrayType{
		Len:         length,
		ElementType: typ,
	}
}

func (t *Type) SetStructType(fieldTypeList []*Type) {
	t.StructType = &StructType{
		FieldTypeList: fieldTypeList,
	}
}

func (t *Type) SetFuncType(paramsTypeList []*Type, resultTypeList []*Type) {
	t.FuncType = &FuncType{
		ParamTypeList:  paramsTypeList,
		ResultTypeList: resultTypeList,
	}
}

func (t *Type) IsArrayType() bool {
	return t.BasicType == BasicTypeArray
}

func (t *Type) IsMapType() bool {
	return t.BasicType == BasicTypeMap
}

func (t *Type) IsInterfaceType() bool {
	return t.BasicType == BasicTypeInterface
}

func (t *Type) IsReferenceType() bool {
	return t.IsArrayType() || t.IsMapType() || t.IsInterfaceType()
}
