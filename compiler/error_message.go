package compiler

const (
	PARSE_ERR int = iota
	CHARACTER_INVALID_ERR
	FUNCTION_MULTIPLE_DEFINE_ERR
	BAD_MULTIBYTE_CHARACTER_ERR
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR
	PARAMETER_MULTIPLE_DEFINE_ERR
	VARIABLE_MULTIPLE_DEFINE_ERR
	IDENTIFIER_NOT_FOUND_ERR
	DERIVE_TYPE_CAST_ERR
	CAST_MISMATCH_ERR
	MATH_TYPE_MISMATCH_ERR
	COMPARE_TYPE_MISMATCH_ERR
	LOGICAL_TYPE_MISMATCH_ERR
	MINUS_TYPE_MISMATCH_ERR
	LOGICAL_NOT_TYPE_MISMATCH_ERR
	INC_DEC_TYPE_MISMATCH_ERR
	FUNCTION_NOT_IDENTIFIER_ERR
	FUNCTION_NOT_FOUND_ERR
	ARGUMENT_COUNT_MISMATCH_ERR
	NOT_LVALUE_ERR
	LABEL_NOT_FOUND_ERR
)

var errMessageMap = map[int]string{
	PARSE_ERR:                             "在($(token))附近发生语法错误",
	CHARACTER_INVALID_ERR:                 "不正确的字符($(bad_char))",
	FUNCTION_MULTIPLE_DEFINE_ERR:          "函数名重复($(name))",
	BAD_MULTIBYTE_CHARACTER_ERR:           "不正确的多字节字符。",
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR: "预期外的宽字符串。",
	PARAMETER_MULTIPLE_DEFINE_ERR:         "函数的参数名重复($(name))。",
	VARIABLE_MULTIPLE_DEFINE_ERR:          "变量名$(name)重复。",
	IDENTIFIER_NOT_FOUND_ERR:              "找不到变量或函数%s。",
	DERIVE_TYPE_CAST_ERR:                  "不能强制转型为派生类型。",
	CAST_MISMATCH_ERR:                     "不能将$(src)转型为$(dest)。",
	MATH_TYPE_MISMATCH_ERR:                "算数运算符的操作数类型不正确。",
	COMPARE_TYPE_MISMATCH_ERR:             "比较运算符的操作数类型不正确。",
	LOGICAL_TYPE_MISMATCH_ERR:             "逻辑and/or运算符的操作数类型不正确。",
	MINUS_TYPE_MISMATCH_ERR:               "减法运算符的操作数类型不正确。",
	LOGICAL_NOT_TYPE_MISMATCH_ERR:         "逻辑非运算符的操作数类型不正确。",
	INC_DEC_TYPE_MISMATCH_ERR:             "自增/自减运算符的操作数类型不正确。",
	FUNCTION_NOT_IDENTIFIER_ERR:           "函数调用运算符的操作数不是函数名。",
	FUNCTION_NOT_FOUND_ERR:                "找不到函数$(name)。",
	ARGUMENT_COUNT_MISMATCH_ERR:           "函数的参数数量错误。Need: %d, Give: %d",
	NOT_LVALUE_ERR:                        "赋值运算符的左边不是一个左值。",
	LABEL_NOT_FOUND_ERR:                   "标签不存在。",
}
