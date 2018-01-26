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
)

// ==============================
// 字节码
// ==============================

const (
    VM_PUSH_INT_1BYTE byte = iota
    VM_PUSH_INT_2BYTE
    VM_PUSH_INT
    VM_PUSH_DOUBLE_0
    VM_PUSH_DOUBLE_1
    VM_PUSH_DOUBLE
    VM_PUSH_STRING
    /**********/
    VM_PUSH_STACK_INT
    VM_PUSH_STACK_DOUBLE
    VM_PUSH_STACK_STRING
    VM_POP_STACK_INT
    VM_POP_STACK_DOUBLE
    VM_POP_STACK_STRING
    /**********/
    VM_PUSH_STATIC_INT
    VM_PUSH_STATIC_DOUBLE
    VM_PUSH_STATIC_STRING
    VM_POP_STATIC_INT
    VM_POP_STATIC_DOUBLE
    VM_POP_STATIC_STRING
    /**********/
    VM_ADD_INT
    VM_ADD_DOUBLE
    VM_ADD_STRING
    VM_SUB_INT
    VM_SUB_DOUBLE
    VM_MUL_INT
    VM_MUL_DOUBLE
    VM_DIV_INT
    VM_DIV_DOUBLE
    VM_MOD_INT
    VM_MOD_DOUBLE
    VM_MINUS_INT
    VM_MINUS_DOUBLE
    VM_INCREMENT
    VM_DECREMENT
    VM_CAST_INT_TO_DOUBLE
    VM_CAST_DOUBLE_TO_INT
    VM_CAST_BOOLEAN_TO_STRING
    VM_CAST_INT_TO_STRING
    VM_CAST_DOUBLE_TO_STRING
    VM_EQ_INT
    VM_EQ_DOUBLE
    VM_EQ_STRING
    VM_GT_INT
    VM_GT_DOUBLE
    VM_GT_STRING
    VM_GE_INT
    VM_GE_DOUBLE
    VM_GE_STRING
    VM_LT_INT
    VM_LT_DOUBLE
    VM_LT_STRING
    VM_LE_INT
    VM_LE_DOUBLE
    VM_LE_STRING
    VM_NE_INT
    VM_NE_DOUBLE
    VM_NE_STRING
    VM_LOGICAL_AND
    VM_LOGICAL_OR
    VM_LOGICAL_NOT
    VM_POP
    VM_DUPLICATE
    VM_JUMP
    VM_JUMP_IF_TRUE
    VM_JUMP_IF_FALSE
    /**********/
    VM_PUSH_FUNCTION
    VM_INVOKE
    VM_RETURN
)
