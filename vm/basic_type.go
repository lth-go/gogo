package vm

// ==============================
// 基本类型
// ==============================

// BasicType 基础类型
type BasicType int

const (
	// BooleanType 布尔类型
	BooleanType BasicType = iota
	// IntType 整形
	IntType
	// DoubleType 浮点
	DoubleType
	// StringType 字符串类型
	StringType
	// Null
	NullType

	// void
	VoidType

	// class
	ClassType

	// array
	BaseType

	// module
	ModuleType
)

