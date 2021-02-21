package vm

// 字节码
const (
	OP_CODE_PUSH_INT_1BYTE byte = iota
	OP_CODE_PUSH_INT_2BYTE
	OP_CODE_PUSH_INT
	OP_CODE_PUSH_FLOAT_0
	OP_CODE_PUSH_FLOAT_1
	OP_CODE_PUSH_FLOAT
	OP_CODE_PUSH_STRING
	OP_CODE_PUSH_NIL

	OP_CODE_PUSH_STACK
	OP_CODE_POP_STACK

	OP_CODE_PUSH_STATIC
	OP_CODE_POP_STATIC

	OP_CODE_PUSH_ARRAY
	OP_CODE_POP_ARRAY
	OP_CODE_PUSH_MAP
	OP_CODE_POP_MAP
	OP_CODE_PUSH_STRUCT
	OP_CODE_POP_STRUCT
	OP_CODE_PUSH_INTERFACE
	OP_CODE_POP_INTERFACE

	OP_CODE_ADD_INT
	OP_CODE_ADD_FLOAT
	OP_CODE_ADD_STRING
	OP_CODE_SUB_INT
	OP_CODE_SUB_FLOAT
	OP_CODE_MUL_INT
	OP_CODE_MUL_FLOAT
	OP_CODE_DIV_INT
	OP_CODE_DIV_FLOAT
	OP_CODE_MOD_INT
	OP_CODE_MOD_FLOAT
	OP_CODE_MINUS_INT
	OP_CODE_MINUS_FLOAT
	OP_CODE_EQ_INT
	OP_CODE_EQ_FLOAT
	OP_CODE_EQ_STRING
	OP_CODE_EQ_OBJECT
	OP_CODE_GT_INT
	OP_CODE_GT_FLOAT
	OP_CODE_GT_STRING
	OP_CODE_GE_INT
	OP_CODE_GE_FLOAT
	OP_CODE_GE_STRING
	OP_CODE_LT_INT
	OP_CODE_LT_FLOAT
	OP_CODE_LT_STRING
	OP_CODE_LE_INT
	OP_CODE_LE_FLOAT
	OP_CODE_LE_STRING
	OP_CODE_NE_INT
	OP_CODE_NE_FLOAT
	OP_CODE_NE_STRING
	OP_CODE_NE_OBJECT
	OP_CODE_LOGICAL_AND
	OP_CODE_LOGICAL_OR
	OP_CODE_LOGICAL_NOT
	OP_CODE_POP
	OP_CODE_DUPLICATE
	OP_CODE_DUPLICATE_OFFSET
	OP_CODE_JUMP
	OP_CODE_JUMP_IF_TRUE
	OP_CODE_JUMP_IF_FALSE

	OP_CODE_PUSH_FUNCTION
	OP_CODE_INVOKE
	OP_CODE_RETURN

	OP_CODE_NEW_ARRAY
	OP_CDOE_NEW_MAP
	OP_CODE_NEW_INTERFACE
	OP_CODE_NEW_STRUCT
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

var OpcodeInfo map[byte]opcodeInfo = map[byte]opcodeInfo{
	OP_CODE_PUSH_INT_1BYTE: {"push_int_1byte", "b", 1},
	OP_CODE_PUSH_INT_2BYTE: {"push_int_2byte", "s", 1},
	OP_CODE_PUSH_INT:       {"push_int", "p", 1},
	OP_CODE_PUSH_FLOAT_0:   {"push_float_0", "", 1},
	OP_CODE_PUSH_FLOAT_1:   {"push_float_1", "", 1},
	OP_CODE_PUSH_FLOAT:     {"push_float", "p", 1},
	OP_CODE_PUSH_STRING:    {"push_string", "p", 1},
	OP_CODE_PUSH_NIL:       {"push_nil", "", 1},

	OP_CODE_PUSH_STACK: {"push_stack", "s", 1},
	OP_CODE_POP_STACK:  {"pop_stack", "s", -1},

	OP_CODE_PUSH_STATIC: {"push_static", "s", 1},
	OP_CODE_POP_STATIC:  {"pop_static", "s", -1},

	OP_CODE_PUSH_ARRAY:     {"push_array", "", 1},
	OP_CODE_POP_ARRAY:      {"pop_array", "", -1},
	OP_CODE_PUSH_MAP:       {"push_map", "", 1},
	OP_CODE_POP_MAP:        {"pop_map", "", -1},
	OP_CODE_PUSH_STRUCT:    {"push_struct", "", 1},
	OP_CODE_POP_STRUCT:     {"pop_struct", "", -1},
	OP_CODE_PUSH_INTERFACE: {"push_interface", "", 1},
	OP_CODE_POP_INTERFACE:  {"pop_interface", "", -1},

	OP_CODE_ADD_INT:          {"add_int", "", -1},
	OP_CODE_ADD_FLOAT:        {"add_float", "", -1},
	OP_CODE_ADD_STRING:       {"add_string", "", -1},
	OP_CODE_SUB_INT:          {"sub_int", "", -1},
	OP_CODE_SUB_FLOAT:        {"sub_float", "", -1},
	OP_CODE_MUL_INT:          {"mul_int", "", -1},
	OP_CODE_MUL_FLOAT:        {"mul_float", "", -1},
	OP_CODE_DIV_INT:          {"div_int", "", -1},
	OP_CODE_DIV_FLOAT:        {"div_float", "", -1},
	OP_CODE_MOD_INT:          {"mod_int", "", -1},
	OP_CODE_MOD_FLOAT:        {"mod_float", "", -1},
	OP_CODE_MINUS_INT:        {"minus_int", "", 0},
	OP_CODE_MINUS_FLOAT:      {"minus_float", "", 0},
	OP_CODE_EQ_INT:           {"eq_int", "", -1},
	OP_CODE_EQ_FLOAT:         {"eq_float", "", -1},
	OP_CODE_EQ_STRING:        {"eq_string", "", -1},
	OP_CODE_EQ_OBJECT:        {"eq_object", "", -1},
	OP_CODE_GT_INT:           {"gt_int", "", -1},
	OP_CODE_GT_FLOAT:         {"gt_float", "", -1},
	OP_CODE_GT_STRING:        {"gt_string", "", -1},
	OP_CODE_GE_INT:           {"ge_int", "", -1},
	OP_CODE_GE_FLOAT:         {"ge_float", "", -1},
	OP_CODE_GE_STRING:        {"ge_string", "", -1},
	OP_CODE_LT_INT:           {"lt_int", "", -1},
	OP_CODE_LT_FLOAT:         {"lt_float", "", -1},
	OP_CODE_LT_STRING:        {"lt_string", "", -1},
	OP_CODE_LE_INT:           {"le_int", "", -1},
	OP_CODE_LE_FLOAT:         {"le_float", "", -1},
	OP_CODE_LE_STRING:        {"le_string", "", -1},
	OP_CODE_NE_INT:           {"ne_int", "", -1},
	OP_CODE_NE_FLOAT:         {"ne_float", "", -1},
	OP_CODE_NE_STRING:        {"ne_string", "", -1},
	OP_CODE_NE_OBJECT:        {"ne_object", "", -1},
	OP_CODE_LOGICAL_AND:      {"logical_and", "", -1},
	OP_CODE_LOGICAL_OR:       {"logical_or", "", -1},
	OP_CODE_LOGICAL_NOT:      {"logical_not", "", 0},
	OP_CODE_POP:              {"pop", "", -1},
	OP_CODE_DUPLICATE:        {"duplicate", "", 1},
	OP_CODE_DUPLICATE_OFFSET: {"duplicate_offset", "s", 1},
	OP_CODE_JUMP:             {"jump", "s", 0},
	OP_CODE_JUMP_IF_TRUE:     {"jump_if_true", "s", -1},
	OP_CODE_JUMP_IF_FALSE:    {"jump_if_false", "s", -1},

	OP_CODE_PUSH_FUNCTION: {"push_function", "s", 1},
	OP_CODE_INVOKE:        {"invoke", "", -1},
	OP_CODE_RETURN:        {"return", "", -1},

	OP_CODE_NEW_ARRAY:     {"new_array", "s", 1},
	OP_CDOE_NEW_MAP:       {"new_map", "s", 1},
	OP_CODE_NEW_INTERFACE: {"new_interface", "s", 1},
	OP_CODE_NEW_STRUCT:    {"new_struct", "s", 1},
}

//
// 行号对应表
//
type LineNumber struct {
	// 源代码行号
	LineNumber int
	// 字节码开始的位置
	StartPc int
	// 接下来有多少字节码对应相同的行号
	PcCount int
}
