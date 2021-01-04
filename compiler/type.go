package compiler

const (
	TypeBool int = iota
	TypeInt
	TypeFloat
	TypeString
	TypeArray
	TypeSlice
	TypeMap
	TypeStruct
)

type Type interface {
	GetType() int
	GetIdentifier() string
	GetPackageName() string
}

type BaseType struct {
	PosImpl
	PackageName string
	Name        string
	// MethodSet   map[string]string
}

type BoolType struct {
	BaseType
}

func (t *BoolType) GetType() int {
	return TypeBool
}

type IntType struct {
	BaseType
}

func (t *IntType) GetType() int {
	return TypeInt
}

type FloatType struct {
	BaseType
}

func (t *FloatType) GetType() int {
	return TypeFloat
}

type StringType struct {
	BaseType
}

func (t *StringType) GetType() int {
	return TypeString
}

type SliceType struct {
	BaseType
	Len int64
	Elt Type
}

func (t *SliceType) GetType() int {
	return TypeSlice
}

type MapType struct {
	BaseType
	Key   Type
	Value Type
}

func (t *MapType) GetType() int {
	return TypeMap
}

type StructType struct {
	BaseType
	Fields []*Field // list of field declarations
}

func (t *StructType) GetType() int {
	return TypeStruct
}

// type FuncType struct {
//     BaseType
//     Func    token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
//     Params  *FieldList // (incoming) parameters; non-nil
//     Results *FieldList // (outgoing) results; or nil
// }

type Field struct {
	Name string
	Type Type
}
