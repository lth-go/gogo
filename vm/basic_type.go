package vm

// ==============================
// 基本类型
// ==============================

// BasicType 基础类型
type BasicType int

const (
	NoType BasicType = iota - 1
	BooleanType
	IntType
	DoubleType
	StringType
	NullType
	VoidType
	BaseType
	ModuleType
	StructType
)
