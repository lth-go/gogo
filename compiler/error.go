package compiler

import (
	"fmt"
	"log"
)

func compileError(pos Position, errorNumber int, a ...interface{}) {
	fmt.Println("编译错误:")
	fmt.Printf("Filename: %s\n", GetCurrentPackage().path)
	fmt.Printf("Line: %d:%d\n", pos.Line, pos.Column)
	errMsg := fmt.Sprintf(errMessageMap[errorNumber], a...)
	msg := fmt.Sprintf("%d\n%s", errorNumber, errMsg)
	panic(msg)
	log.Fatalf(msg)
}

const (
	PARSE_ERR int = iota
	CHARACTER_INVALID_ERR
	FUNCTION_MULTIPLE_DEFINE_ERR
	PARAMETER_MULTIPLE_DEFINE_ERR
	VARIABLE_MULTIPLE_DEFINE_ERR
	IDENTIFIER_NOT_FOUND_ERR
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
	ARGUMENT_TYPE_MISMATCH_ERR
	NOT_LVALUE_ERR
	LABEL_NOT_FOUND_ERR
	INDEX_LEFT_OPERAND_NOT_ARRAY_ERR
	INDEX_NOT_INT_ERR
	ARRAY_SIZE_NOT_INT_ERR
	DIVISION_BY_ZERO_IN_COMPILE_ERR
	PACKAGE_NAME_TOO_LONG_ERR
	REQUIRE_FILE_NOT_FOUND_ERR
	REQUIRE_DUPLICATE_ERR
	MEMBER_EXPRESSION_TYPE_ERR
	RETURN_IN_VOID_FUNCTION_ERR
	CLASS_NOT_FOUND_ERR
	FIELD_CAN_NOT_CALL_ERR
	METHOD_IS_NOT_CALLED_ERR
	ASSIGN_TO_METHOD_ERR
	FIELD_OF_SUPER_REFERENCED_ERR
	FIELD_OVERRIDED_ERR
	FIELD_NAME_DUPLICATE_ERR
	ARRAY_METHOD_NOT_FOUND_ERR
	STRING_METHOD_NOT_FOUND_ERR
	IF_CONDITION_NOT_BOOLEAN_ERR
	FOR_CONDITION_NOT_BOOLEAN_ERR
	BAD_PARAMETER_COUNT_ERR
	BAD_PARAMETER_TYPE_ERR
	BAD_RETURN_TYPE_ERR
	TYPE_NAME_NOT_FOUND_ERR
)

var errMessageMap map[int]string = map[int]string{
	PARSE_ERR:                        "在($(token))附近发生语法错误",
	CHARACTER_INVALID_ERR:            "不正确的字符($(bad_char))",
	FUNCTION_MULTIPLE_DEFINE_ERR:     "函数名重复($(name))",
	PARAMETER_MULTIPLE_DEFINE_ERR:    "函数的参数名重复(%s)。",
	VARIABLE_MULTIPLE_DEFINE_ERR:     "变量名$(name)重复。",
	IDENTIFIER_NOT_FOUND_ERR:         "找不到变量或函数(%s)。",
	CAST_MISMATCH_ERR:                "不能将%+v转型为%v。",
	MATH_TYPE_MISMATCH_ERR:           "算数运算符的操作数类型不正确。",
	COMPARE_TYPE_MISMATCH_ERR:        "比较运算符的操作数类型不正确。",
	LOGICAL_TYPE_MISMATCH_ERR:        "逻辑and/or运算符的操作数类型不正确。",
	MINUS_TYPE_MISMATCH_ERR:          "减法运算符的操作数类型不正确。",
	LOGICAL_NOT_TYPE_MISMATCH_ERR:    "逻辑非运算符的操作数类型不正确。",
	INC_DEC_TYPE_MISMATCH_ERR:        "自增/自减运算符的操作数类型不正确。",
	FUNCTION_NOT_IDENTIFIER_ERR:      "函数调用运算符的操作数不是函数名。",
	FUNCTION_NOT_FOUND_ERR:           "找不到函数%s。",
	ARGUMENT_COUNT_MISMATCH_ERR:      "函数的参数数量错误, %s %v %v",
	ARGUMENT_TYPE_MISMATCH_ERR:       "函数的参数类型错误.",
	NOT_LVALUE_ERR:                   "赋值运算符的左边不是一个左边值。",
	LABEL_NOT_FOUND_ERR:              "标签$(label)不存在。",
	INDEX_LEFT_OPERAND_NOT_ARRAY_ERR: "下标运算符[]的左边不是数组类型",
	INDEX_NOT_INT_ERR:                "数组的下标不是int。",
	ARRAY_SIZE_NOT_INT_ERR:           "数组的大小不是int。",
	DIVISION_BY_ZERO_IN_COMPILE_ERR:  "整数值不能被0除。",
	PACKAGE_NAME_TOO_LONG_ERR:        "package名称过长",
	REQUIRE_FILE_NOT_FOUND_ERR:       "被import的文件不存在($(file))",
	REQUIRE_DUPLICATE_ERR:            "源文件中重复import了包($(package_name))。",
	MEMBER_EXPRESSION_TYPE_ERR:       "该类型不能使用成员运算符。",
	RETURN_IN_VOID_FUNCTION_ERR:      "void类型的函数不能有返回值。",
	CLASS_NOT_FOUND_ERR:              "没有找到类$(name)。",
	FIELD_CAN_NOT_CALL_ERR:           "不能调用字段$(member_name)",
	METHOD_IS_NOT_CALLED_ERR:         "方法$(member_name)不能出现在函数调用之外的位置。",
	ASSIGN_TO_METHOD_ERR:             "尝试为方法$(member_name)赋值。",
	FIELD_OF_SUPER_REFERENCED_ERR:    "不能引用super的字段。",
	FIELD_OVERRIDED_ERR:              "$(name)是字段，不能覆盖。",
	FIELD_NAME_DUPLICATE_ERR:         "重复的字段名$(name)。",
	ARRAY_METHOD_NOT_FOUND_ERR:       "数组中没有$(name)方法。",
	STRING_METHOD_NOT_FOUND_ERR:      "数组中没有$(name)方法。",
	IF_CONDITION_NOT_BOOLEAN_ERR:     "if语句的条件表达式不是boolean型。",
	FOR_CONDITION_NOT_BOOLEAN_ERR:    "for语句的条件表达式不是boolean型。",
	BAD_PARAMETER_COUNT_ERR:          "方法或函数$(name)的参数数量错误。",
	BAD_PARAMETER_TYPE_ERR:           "方法或函数$(func_name)的第$(index)个参数, $(param_name)的类型错误。",
	BAD_RETURN_TYPE_ERR:              "方法或函数$(name)的返回值类型错误。",
	TYPE_NAME_NOT_FOUND_ERR:          "找不到类型名$(name)。",
}
