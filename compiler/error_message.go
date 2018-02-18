package compiler

const (
	PARSE_ERR int = iota
	CHARACTER_INVALID_ERR
	FUNCTION_MULTIPLE_DEFINE_ERR
	BAD_MULTIBYTE_CHARACTER_ERR
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR
	ARRAY_ELEMENT_CAN_NOT_BE_FINAL_ERR
	COMPLEX_ASSIGNMENT_OPERATOR_TO_FINAL_ERR
	PARAMETER_MULTIPLE_DEFINE_ERR
	VARIABLE_MULTIPLE_DEFINE_ERR
	IDENTIFIER_NOT_FOUND_ERR
	FUNCTION_IDENTIFIER_ERR
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
	ARRAY_LITERAL_EMPTY_ERR
	INDEX_LEFT_OPERAND_NOT_ARRAY_ERR
	INDEX_NOT_INT_ERR
	ARRAY_SIZE_NOT_INT_ERR
	DIVISION_BY_ZERO_IN_COMPILE_ERR
)

var errMessageMap = map[int]string{
	PARSE_ERR:                                "在(%s)附近发生语法错误",
	CHARACTER_INVALID_ERR:                    "不正确的字符(%s)",
	FUNCTION_MULTIPLE_DEFINE_ERR:             "函数名重复(%s)",
	BAD_MULTIBYTE_CHARACTER_ERR:              "不正确的多字节字符。",
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR:    "预期外的宽字符串。",
	ARRAY_ELEMENT_CAN_NOT_BE_FINAL_ERR:       "数组元素不能标识为final。",
	COMPLEX_ASSIGNMENT_OPERATOR_TO_FINAL_ERR: "复合赋值运算符不能用于final值",
	PARAMETER_MULTIPLE_DEFINE_ERR:            "函数的参数名重复(%s)。",
	VARIABLE_MULTIPLE_DEFINE_ERR:             "变量名(%s)重复。",
	IDENTIFIER_NOT_FOUND_ERR:                 "找不到变量或函数%s。",
	FUNCTION_IDENTIFIER_ERR:                  "$(name)是函数名，但没有函数调用的()。",
	DERIVE_TYPE_CAST_ERR:                     "不能强制转型为派生类型。",
	CAST_MISMATCH_ERR:                        "不能将(%s)转型为(%s)。",
	MATH_TYPE_MISMATCH_ERR:                   "算数运算符的操作数类型不正确。",
	COMPARE_TYPE_MISMATCH_ERR:                "比较运算符的操作数类型不正确。Left: %s, Right: %s",
	LOGICAL_TYPE_MISMATCH_ERR:                "逻辑and/or运算符的操作数类型不正确。",
	MINUS_TYPE_MISMATCH_ERR:                  "减法运算符的操作数类型不正确。",
	LOGICAL_NOT_TYPE_MISMATCH_ERR:            "逻辑非运算符的操作数类型不正确。",
	INC_DEC_TYPE_MISMATCH_ERR:                "自增/自减运算符的操作数类型不正确。",
	FUNCTION_NOT_IDENTIFIER_ERR:              "函数调用运算符的操作数不是函数名。",
	FUNCTION_NOT_FOUND_ERR:                   "找不到函数(%s)。",
	ARGUMENT_COUNT_MISMATCH_ERR:              "函数的参数数量错误。Need: %d, Give: %d",
	NOT_LVALUE_ERR:                           "赋值运算符的左边不是一个左值。",
	LABEL_NOT_FOUND_ERR:                      "标签不存在。",
	ARRAY_LITERAL_EMPTY_ERR:                  "数组字面量必须至少有一个元素",
	INDEX_LEFT_OPERAND_NOT_ARRAY_ERR:         "下标运算符[]的左边不是数组类型",
	INDEX_NOT_INT_ERR:                        "数组的下标不是int。",
	ARRAY_SIZE_NOT_INT_ERR:                   "数组的大小不是int。",
	DIVISION_BY_ZERO_IN_COMPILE_ERR:               "整数值不能被0除。",
}
