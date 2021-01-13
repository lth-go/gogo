package vm

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
	VM_PUSH_NULL
	/**********/
	VM_PUSH_STACK_INT
	VM_PUSH_STACK_DOUBLE
	VM_PUSH_STACK_OBJECT
	VM_POP_STACK_INT
	VM_POP_STACK_DOUBLE
	VM_POP_STACK_OBJECT
	/**********/
	VM_PUSH_HEAP_INT
	VM_PUSH_HEAP_FLOAT
	VM_PUSH_HEAP_OBJECT
	VM_POP_HEAP_INT
	VM_POP_HEAP_FLOAT
	VM_POP_HEAP_OBJECT
	/**********/
	VM_PUSH_ARRAY_INT
	VM_PUSH_ARRAY_DOUBLE
	VM_PUSH_ARRAY_OBJECT
	VM_POP_ARRAY_INT
	VM_POP_ARRAY_DOUBLE
	VM_POP_ARRAY_OBJECT
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
	VM_UP_CAST
	VM_DOWN_CAST
	VM_EQ_INT
	VM_EQ_DOUBLE
	VM_EQ_OBJECT
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
	VM_NE_OBJECT
	VM_NE_STRING
	VM_LOGICAL_AND
	VM_LOGICAL_OR
	VM_LOGICAL_NOT
	VM_POP
	VM_DUPLICATE
	VM_DUPLICATE_OFFSET
	VM_JUMP
	VM_JUMP_IF_TRUE
	VM_JUMP_IF_FALSE
	/**********/
	VM_PUSH_FUNCTION
	VM_INVOKE
	VM_RETURN
	/**********/
	VM_NEW_ARRAY_LITERAL_INT
	VM_NEW_ARRAY_LITERAL_DOUBLE
	VM_NEW_ARRAY_LITERAL_OBJECT
)

type opcodeInfo struct {
	// 注记符
	Mnemonic string

	// 参数类型，
	// `b` 一个字节整数
	// `s` 两个字节整数
	// `p` 常量池索引值
	Parameter      string
	stackIncrement int
}

var OpcodeInfo []opcodeInfo = []opcodeInfo{
	{"push_int_1byte", "b", 1},
	{"push_int_2byte", "s", 1},
	{"push_int", "p", 1},
	{"push_double_0", "", 1},
	{"push_double_1", "", 1},
	{"push_double", "p", 1},
	{"push_string", "p", 1},
	{"push_null", "", 1},
	/**********/
	{"push_stack_int", "s", 1},
	{"push_stack_double", "s", 1},
	{"push_stack_object", "s", 1},
	{"pop_stack_int", "s", -1},
	{"pop_stack_double", "s", -1},
	{"pop_stack_object", "s", -1},
	/**********/
	{"push_static_int", "s", 1},
	{"push_static_double", "s", 1},
	{"push_static_object", "s", 1},
	{"pop_static_int", "s", -1},
	{"pop_static_double", "s", -1},
	{"pop_static_object", "s", -1},
	/**********/
	{"push_array_int", "", 1},
	{"push_array_double", "", 1},
	{"push_array_object", "", 1},
	{"pop_array_int", "", -1},
	{"pop_array_double", "", -1},
	{"pop_array_object", "", -1},
	/**********/
	{"add_int", "", -1},
	{"add_double", "", -1},
	{"add_string", "", -1},
	{"sub_int", "", -1},
	{"sub_double", "", -1},
	{"mul_int", "", -1},
	{"mul_double", "", -1},
	{"div_int", "", -1},
	{"div_double", "", -1},
	{"mod_int", "", -1},
	{"mod_double", "", -1},
	{"minus_int", "", 0},
	{"minus_double", "", 0},
	{"increment", "", 0},
	{"decrement", "", 0},
	{"cast_int_to_double", "", 0},
	{"cast_double_to_int", "", 0},
	{"cast_boolean_to_string", "", 0},
	{"cast_int_to_string", "", 0},
	{"cast_double_to_string", "", 0},
	{"up_cast", "s", 0},
	{"down_cast", "s", 0},
	{"eq_int", "", -1},
	{"eq_double", "", -1},
	{"eq_object", "", -1},
	{"eq_string", "", -1},
	{"gt_int", "", -1},
	{"gt_double", "", -1},
	{"gt_string", "", -1},
	{"ge_int", "", -1},
	{"ge_double", "", -1},
	{"ge_string", "", -1},
	{"lt_int", "", -1},
	{"lt_double", "", -1},
	{"lt_string", "", -1},
	{"le_int", "", -1},
	{"le_double", "", -1},
	{"le_string", "", -1},
	{"ne_int", "", -1},
	{"ne_double", "", -1},
	{"ne_object", "", -1},
	{"ne_string", "", -1},
	{"logical_and", "", -1},
	{"logical_or", "", -1},
	{"logical_not", "", 0},
	{"pop", "", -1},
	{"duplicate", "", 1},
	{"duplicate_offset", "s", 1},
	{"jump", "s", 0},
	{"jump_if_true", "s", -1},
	{"jump_if_false", "s", -1},
	/**********/
	{"push_function", "s", 1},
	{"invoke", "", -1},
	{"return", "", -1},
	/**********/
	{"new_array_literal_int", "s", 1},
	{"new_array_literal_double", "s", 1},
	{"new_array_literal_object", "s", 1},
}
