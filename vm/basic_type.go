package vm

// ==============================
// 基本类型
// ==============================

// BasicType 基础类型
type BasicType int

const (
	BasicTypeNoType BasicType = iota - 1
	BasicTypeBool
	BasicTypeInt
	BasicTypeFloat
	BasicTypeString
	BasicTypeNil
	BasicTypeVoid
	BasicTypeBase
	BasicTypeModule
	BasicTypeSlice
	BasicTypeStruct
)
