package compiler

import (
	"fmt"
	"log"
)

func compileError(pos Position, errorNumber int, a ...interface{}) {
	fmt.Println("编译错误:")
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
	BAD_MULTIBYTE_CHARACTER_ERR
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR
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
	PACKAGE_NAME_TOO_LONG_ERR
	REQUIRE_FILE_NOT_FOUND_ERR
	REQUIRE_FILE_ERR
	REQUIRE_DUPLICATE_ERR
	RENAME_HAS_NO_PACKAGED_NAME_ERR
	ABSTRACT_MULTIPLE_SPECIFIED_ERR
	ACCESS_MODIFIER_MULTIPLE_SPECIFIED_ERR
	OVERRIDE_MODIFIER_MULTIPLE_SPECIFIED_ERR
	VIRTUAL_MODIFIER_MULTIPLE_SPECIFIED_ERR
	MEMBER_EXPRESSION_TYPE_ERR
	MEMBER_NOT_FOUND_ERR
	PRIVATE_MEMBER_ACCESS_ERR
	ABSTRACT_METHOD_HAS_BODY_ERR
	CONCRETE_METHOD_HAS_NO_BODY_ERR
	MULTIPLE_INHERITANCE_ERR
	INHERIT_CONCRETE_CLASS_ERR
	NEW_ABSTRACT_CLASS_ERR
	RETURN_IN_VOID_FUNCTION_ERR
	CLASS_NOT_FOUND_ERR
	CONSTRUCTOR_IS_FIELD_ERR
	NOT_CONSTRUCTOR_ERR
	FIELD_CAN_NOT_CALL_ERR
	METHOD_IS_NOT_CALLED_ERR
	ASSIGN_TO_METHOD_ERR
	NON_VIRTUAL_METHOD_OVERRIDED_ERR
	NEED_OVERRIDE_ERR
	ABSTRACT_METHOD_IN_CONCRETE_CLASS_ERR
	HASNT_SUPER_CLASS_ERR
	SUPER_NOT_IN_MEMBER_EXPRESSION_ERR
	FIELD_OF_SUPER_REFERENCED_ERR
	FIELD_OVERRIDED_ERR
	FIELD_NAME_DUPLICATE_ERR
	ARRAY_METHOD_NOT_FOUND_ERR
	STRING_METHOD_NOT_FOUND_ERR
	INSTANCEOF_OPERAND_NOT_REFERENCE_ERR
	INSTANCEOF_TYPE_NOT_REFERENCE_ERR
	INSTANCEOF_FOR_NOT_CLASS_ERR
	INSTANCEOF_MUST_RETURN_TRUE_ERR
	INSTANCEOF_MUST_RETURN_FALSE_ERR
	INSTANCEOF_INTERFACE_ERR
	DOWN_CAST_OPERAND_IS_NOT_CLASS_ERR
	DOWN_CAST_TARGET_IS_NOT_CLASS_ERR
	DOWN_CAST_DO_NOTHING_ERR
	DOWN_CAST_TO_SUPER_CLASS_ERR
	DOWN_CAST_TO_BAD_CLASS_ERR
	DOWN_CAST_INTERFACE_ERR
	REQUIRE_ITSELF_ERR
	IF_CONDITION_NOT_BOOLEAN_ERR
	WHILE_CONDITION_NOT_BOOLEAN_ERR
	FOR_CONDITION_NOT_BOOLEAN_ERR
	DO_WHILE_CONDITION_NOT_BOOLEAN_ERR
	OVERRIDE_METHOD_ACCESSIBILITY_ERR
	BAD_PARAMETER_COUNT_ERR
	BAD_PARAMETER_TYPE_ERR
	BAD_RETURN_TYPE_ERR
	CONSTRUCTOR_CALLED_ERR
	TYPE_NAME_NOT_FOUND_ERR
	INTERFACE_INHERIT_ERR
	PACKAGE_MEMBER_ACCESS_ERR
	PACKAGE_CLASS_ACCESS_ERR
	THIS_OUT_OF_CLASS_ERR
	SUPER_OUT_OF_CLASS_ERR
	EOF_IN_C_COMMENT_ERR
	EOF_IN_STRING_LITERAL_ERR
	TOO_LONG_CHARACTER_LITERAL_ERR
	COMPILE_ERROR_COUNT_PLUS_1
)

var errMessageMap map[int]string = map[int]string{
	PARSE_ERR:                                "在($(token))附近发生语法错误",
	CHARACTER_INVALID_ERR:                    "不正确的字符($(bad_char))",
	FUNCTION_MULTIPLE_DEFINE_ERR:             "函数名重复($(name))",
	BAD_MULTIBYTE_CHARACTER_ERR:              "不正确的多字节字符。",
	UNEXPECTED_WIDE_STRING_IN_COMPILE_ERR:    "预期外的宽字符串。",
	PARAMETER_MULTIPLE_DEFINE_ERR:            "函数的参数名重复($(name))。",
	VARIABLE_MULTIPLE_DEFINE_ERR:             "变量名$(name)重复。",
	IDENTIFIER_NOT_FOUND_ERR:                 "找不到变量或函数$(name)。",
	FUNCTION_IDENTIFIER_ERR:                  "$(name)是函数名，但没有函数调用的()。",
	DERIVE_TYPE_CAST_ERR:                     "不能强制转型为派生类型。",
	CAST_MISMATCH_ERR:                        "不能将%+v转型为%v。",
	MATH_TYPE_MISMATCH_ERR:                   "算数运算符的操作数类型不正确。",
	COMPARE_TYPE_MISMATCH_ERR:                "比较运算符的操作数类型不正确。",
	LOGICAL_TYPE_MISMATCH_ERR:                "逻辑and/or运算符的操作数类型不正确。",
	MINUS_TYPE_MISMATCH_ERR:                  "减法运算符的操作数类型不正确。",
	LOGICAL_NOT_TYPE_MISMATCH_ERR:            "逻辑非运算符的操作数类型不正确。",
	INC_DEC_TYPE_MISMATCH_ERR:                "自增/自减运算符的操作数类型不正确。",
	FUNCTION_NOT_IDENTIFIER_ERR:              "函数调用运算符的操作数不是函数名。",
	FUNCTION_NOT_FOUND_ERR:                   "找不到函数$(name)。",
	ARGUMENT_COUNT_MISMATCH_ERR:              "函数的参数数量错误。",
	NOT_LVALUE_ERR:                           "赋值运算符的左边不是一个左边值。",
	LABEL_NOT_FOUND_ERR:                      "标签$(label)不存在。",
	ARRAY_LITERAL_EMPTY_ERR:                  "数组字面量必须至少有一个元素",
	INDEX_LEFT_OPERAND_NOT_ARRAY_ERR:         "下标运算符[]的左边不是数组类型",
	INDEX_NOT_INT_ERR:                        "数组的下标不是int。",
	ARRAY_SIZE_NOT_INT_ERR:                   "数组的大小不是int。",
	DIVISION_BY_ZERO_IN_COMPILE_ERR:          "整数值不能被0除。",
	PACKAGE_NAME_TOO_LONG_ERR:                "package名称过长",
	REQUIRE_FILE_NOT_FOUND_ERR:               "被import的文件不存在($(file))",
	REQUIRE_FILE_ERR:                         "import时发生错误($(status))。",
	REQUIRE_DUPLICATE_ERR:                    "源文件中重复import了包($(package_name))。",
	RENAME_HAS_NO_PACKAGED_NAME_ERR:          "rename后的名称必须指定package。",
	ABSTRACT_MULTIPLE_SPECIFIED_ERR:          "重复声明了abstract。",
	ACCESS_MODIFIER_MULTIPLE_SPECIFIED_ERR:   "重复声明了访问修饰符。",
	OVERRIDE_MODIFIER_MULTIPLE_SPECIFIED_ERR: "重复声明了override。",
	VIRTUAL_MODIFIER_MULTIPLE_SPECIFIED_ERR:  "重复声明了virtual。",
	MEMBER_EXPRESSION_TYPE_ERR:               "该类型不能使用成员运算符。",
	MEMBER_NOT_FOUND_ERR:                     "在类型$(class_name)中不存在成员$(member_name)。",
	PRIVATE_MEMBER_ACCESS_ERR:                "成员$(member_name)是private的，不能访问。",
	ABSTRACT_METHOD_HAS_BODY_ERR:             "没有实现abstract方法。",
	CONCRETE_METHOD_HAS_NO_BODY_ERR:          "必须实现非abstract方法。",
	MULTIPLE_INHERITANCE_ERR:                 "继承了多个类。",
	INHERIT_CONCRETE_CLASS_ERR:               "Diksam中只能继承abstract类(类$(name)不是abstract类)。",
	NEW_ABSTRACT_CLASS_ERR:                   "不能对abstract类($(name))使用new。",
	RETURN_IN_VOID_FUNCTION_ERR:              "void类型的函数不能有返回值。",
	CLASS_NOT_FOUND_ERR:                      "没有找到类$(name)。",
	CONSTRUCTOR_IS_FIELD_ERR:                 "被指定为构造方法的成员$(member_name)不是一个方法。",
	NOT_CONSTRUCTOR_ERR:                      "用来new的方法$(member_name)并不是构造方法。",
	FIELD_CAN_NOT_CALL_ERR:                   "不能调用字段$(member_name)",
	METHOD_IS_NOT_CALLED_ERR:                 "方法$(member_name)不能出现在函数调用之外的位置。",
	ASSIGN_TO_METHOD_ERR:                     "尝试为方法$(member_name)赋值。",
	NON_VIRTUAL_METHOD_OVERRIDED_ERR:         "不能覆盖非virtual方法$(name)。",
	NEED_OVERRIDE_ERR:                        "覆盖方法时必须使用override关键字($(name))。",
	ABSTRACT_METHOD_IN_CONCRETE_CLASS_ERR:    "在abstract类中，存在非abstract方法$(method_name)。",
	HASNT_SUPER_CLASS_ERR:                    "在没有超类的类中使用了super。",
	SUPER_NOT_IN_MEMBER_EXPRESSION_ERR:       "方法调用以外不能使用super。",
	FIELD_OF_SUPER_REFERENCED_ERR:            "不能引用super的字段。",
	FIELD_OVERRIDED_ERR:                      "$(name)是字段，不能覆盖。",
	FIELD_NAME_DUPLICATE_ERR:                 "重复的字段名$(name)。",
	ARRAY_METHOD_NOT_FOUND_ERR:               "数组中没有$(name)方法。",
	STRING_METHOD_NOT_FOUND_ERR:              "数组中没有$(name)方法。",
	INSTANCEOF_OPERAND_NOT_REFERENCE_ERR:     "instanceof的操作数必须是引用类型。",
	INSTANCEOF_TYPE_NOT_REFERENCE_ERR:        "instanceof的右边的类型必须是引用类型。",
	INSTANCEOF_FOR_NOT_CLASS_ERR:             "instanceof的目标必须是类。",
	INSTANCEOF_MUST_RETURN_TRUE_ERR:          "instanceof语句一直为真。",
	INSTANCEOF_MUST_RETURN_FALSE_ERR:         "instanceof语句一直为假。",
	INSTANCEOF_INTERFACE_ERR:                 "因为Diksam的接口间没有父子关系, instanceof语句一直为假。",
	DOWN_CAST_OPERAND_IS_NOT_CLASS_ERR:       "向下转型的源类型必须是类。",
	DOWN_CAST_TARGET_IS_NOT_CLASS_ERR:        "向下转型的目标类型必须是类。",
	DOWN_CAST_DO_NOTHING_ERR:                 "不需要进行向下转型。",
	DOWN_CAST_TO_SUPER_CLASS_ERR:             "尝试将父类转换为子类。",
	DOWN_CAST_TO_BAD_CLASS_ERR:               "尝试转换没有继承关系的类。",
	DOWN_CAST_INTERFACE_ERR:                  "因为Diksam的接口间没有父子关系, 不能向下转型。",
	REQUIRE_ITSELF_ERR:                       "不能import文件本身。",
	IF_CONDITION_NOT_BOOLEAN_ERR:             "if语句的条件表达式不是boolean型。",
	WHILE_CONDITION_NOT_BOOLEAN_ERR:          "while语句的条件表达式不是boolean型。",
	FOR_CONDITION_NOT_BOOLEAN_ERR:            "for语句的条件表达式不是boolean型。",
	DO_WHILE_CONDITION_NOT_BOOLEAN_ERR:       "do while语句的条件表达式不是boolean型。",
	OVERRIDE_METHOD_ACCESSIBILITY_ERR:        "被覆盖的方法$(name)的访问修饰符必须比超类的更严格。",
	BAD_PARAMETER_COUNT_ERR:                  "方法或函数$(name)的参数数量错误。",
	BAD_PARAMETER_TYPE_ERR:                   "方法或函数$(func_name)的第$(index)个参数, $(param_name)的类型错误。",
	BAD_RETURN_TYPE_ERR:                      "方法或函数$(name)的返回值类型错误。",
	CONSTRUCTOR_CALLED_ERR:                   "不能直接调用构造方法。",
	TYPE_NAME_NOT_FOUND_ERR:                  "找不到类型名$(name)。",
	INTERFACE_INHERIT_ERR:                    "Diksam的接口之间不能继承（至今为止）。",
	PACKAGE_MEMBER_ACCESS_ERR:                "不能从包外访问成员$(member_name)。",
	PACKAGE_CLASS_ACCESS_ERR:                 "不能从包外访问类$(class_name)。",
	THIS_OUT_OF_CLASS_ERR:                    "不能在类外使用this。",
	SUPER_OUT_OF_CLASS_ERR:                   "不能在类外使用super。",
	EOF_IN_C_COMMENT_ERR:                     "在C样式的注释中终止了文件。",
	EOF_IN_STRING_LITERAL_ERR:                "在字符串字面量中终止了文件。",
	TOO_LONG_CHARACTER_LITERAL_ERR:           "字符字面量中包含了2个以上的字符。",
}
