package vm

// ==============================
// 基本类型
// ==============================

// BasicType 基础类型
type BasicType int

const (
	BooleanType BasicType = iota
	IntType
	DoubleType
	StringType
	NullType
	VoidType
	BaseType
	ModuleType
	StructType
)
