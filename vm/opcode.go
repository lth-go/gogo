package vm

// ==============================
// 字节码
// ==============================

const (
	VM_PUSH_INT_1BYTE byte = iota
	VM_PUSH_INT_2BYTE
	VM_PUSH_INT
	VM_PUSH_FLOAT_0
	VM_PUSH_FLOAT_1
	VM_PUSH_FLOAT
	VM_PUSH_STRING
	VM_PUSH_NIL
	/**********/
	VM_PUSH_STACK_INT
	VM_PUSH_STACK_FLOAT
	VM_PUSH_STACK_OBJECT
	VM_POP_STACK_INT
	VM_POP_STACK_FLOAT
	VM_POP_STACK_OBJECT
	/**********/
	VM_PUSH_STATIC_INT
	VM_PUSH_STATIC_FLOAT
	VM_PUSH_STATIC_OBJECT
	VM_POP_STATIC_INT
	VM_POP_STATIC_FLOAT
	VM_POP_STATIC_OBJECT
	/**********/
	VM_PUSH_ARRAY_OBJECT
	VM_POP_ARRAY_OBJECT
	/**********/
	VM_ADD_INT
	VM_ADD_FLOAT
	VM_ADD_STRING
	VM_SUB_INT
	VM_SUB_FLOAT
	VM_MUL_INT
	VM_MUL_FLOAT
	VM_DIV_INT
	VM_DIV_FLOAT
	VM_MOD_INT
	VM_MOD_FLOAT
	VM_MINUS_INT
	VM_MINUS_FLOAT
	VM_INCREMENT
	VM_DECREMENT
	VM_CAST_INT_TO_FLOAT
	VM_CAST_FLOAT_TO_INT
	VM_CAST_BOOLEAN_TO_STRING
	VM_CAST_INT_TO_STRING
	VM_CAST_FLOAT_TO_STRING
	VM_UP_CAST
	VM_DOWN_CAST
	VM_EQ_INT
	VM_EQ_FLOAT
	VM_EQ_OBJECT
	VM_EQ_STRING
	VM_GT_INT
	VM_GT_FLOAT
	VM_GT_STRING
	VM_GE_INT
	VM_GE_FLOAT
	VM_GE_STRING
	VM_LT_INT
	VM_LT_FLOAT
	VM_LT_STRING
	VM_LE_INT
	VM_LE_FLOAT
	VM_LE_STRING
	VM_NE_INT
	VM_NE_FLOAT
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
	VM_NEW_ARRAY
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
	{"push_float_0", "", 1},
	{"push_float_1", "", 1},
	{"push_float", "p", 1},
	{"push_string", "p", 1},
	{"push_null", "", 1},
	/**********/
	{"push_stack_int", "s", 1},
	{"push_stack_float", "s", 1},
	{"push_stack_object", "s", 1},
	{"pop_stack_int", "s", -1},
	{"pop_stack_float", "s", -1},
	{"pop_stack_object", "s", -1},
	/**********/
	{"push_static_int", "s", 1},
	{"push_static_float", "s", 1},
	{"push_static_object", "s", 1},
	{"pop_static_int", "s", -1},
	{"pop_static_float", "s", -1},
	{"pop_static_object", "s", -1},
	/**********/
	{"push_array_object", "", 1},
	{"pop_array_object", "", -1},
	/**********/
	{"add_int", "", -1},
	{"add_float", "", -1},
	{"add_string", "", -1},
	{"sub_int", "", -1},
	{"sub_float", "", -1},
	{"mul_int", "", -1},
	{"mul_float", "", -1},
	{"div_int", "", -1},
	{"div_float", "", -1},
	{"mod_int", "", -1},
	{"mod_float", "", -1},
	{"minus_int", "", 0},
	{"minus_float", "", 0},
	{"increment", "", 0},
	{"decrement", "", 0},
	{"cast_int_to_float", "", 0},
	{"cast_float_to_int", "", 0},
	{"cast_boolean_to_string", "", 0},
	{"cast_int_to_string", "", 0},
	{"cast_float_to_string", "", 0},
	{"up_cast", "s", 0},
	{"down_cast", "s", 0},
	{"eq_int", "", -1},
	{"eq_float", "", -1},
	{"eq_object", "", -1},
	{"eq_string", "", -1},
	{"gt_int", "", -1},
	{"gt_float", "", -1},
	{"gt_string", "", -1},
	{"ge_int", "", -1},
	{"ge_float", "", -1},
	{"ge_string", "", -1},
	{"lt_int", "", -1},
	{"lt_float", "", -1},
	{"lt_string", "", -1},
	{"le_int", "", -1},
	{"le_float", "", -1},
	{"le_string", "", -1},
	{"ne_int", "", -1},
	{"ne_float", "", -1},
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
	{"new_array", "s", 1},
}
