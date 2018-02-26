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
	DIVISION_BY_ZERO_IN_COMPILE_ERR:          "整数值不能被0除。",
}

var  errMessageList []string {
    "在($(token))附近发生语法错误",
    "不正确的字符($(bad_char))",
    "函数名重复($(name))",
    "不正确的多字节字符。",
    "预期外的宽字符串。",
    "函数的参数名重复($(name))。",
    "变量名$(name)重复。",
    "找不到变量或函数$(name)。",
    "$(name)是函数名，但没有函数调用的()。",
    "不能强制转型为派生类型。",
    "不能将$(src)转型为$(dest)。",
    "算数运算符的操作数类型不正确。",
    "比较运算符的操作数类型不正确。",
    "逻辑and/or运算符的操作数类型不正确。",
    "减法运算符的操作数类型不正确。",
    "逻辑非运算符的操作数类型不正确。",
    "自增/自减运算符的操作数类型不正确。",
    "函数调用运算符的操作数不是函数名。",
    "找不到函数$(name)。",
    "函数的参数数量错误。",
    "赋值运算符的左边不是一个左边值。",
    "标签$(label)不存在。",
    "数组字面量必须至少有一个元素",
    "下标运算符[]的左边不是数组类型",
    "数组的下标不是int。",
    "数组的大小不是int。",
    "整数值不能被0除。",
    "package名称过长",
    "被require的文件不存在($(file))",
    "require时发生错误($(status))。",
    "源文件中重复require了包($(package_name))。",
    "rename后的名称必须指定package。",
    "重复声明了abstract。",
    "重复声明了访问修饰符。",
    "重复声明了override。",
    "重复声明了virtual。",
    "该类型不能使用成员运算符。",
    "在类型$(class_name)中不存在成员$(member_name)。",
    "成员$(member_name)是private的，不能访问。",
    "没有实现abstract方法。",
    "必须实现非abstract方法。",
    "继承了多个类。",
    "Diksam中只能继承abstract类(类$(name)不是abstract类)。",
    "不能对abstract类($(name))使用new。",
    "void类型的函数不能有返回值。",
    "没有找到类$(name)。",
    "被指定为构造方法的成员$(member_name)不是一个方法。",
    "用来new的方法$(member_name)并不是构造方法。",
    "不能调用字段$(member_name)",
    "方法$(member_name)不能出现在函数调用之外的位置。",
    "尝试为方法$(member_name)赋值。",
    "不能覆盖非virtual方法$(name)。",
    "覆盖方法时必须使用override关键字($(name))。",
    "在abstract类中，存在非abstract方法$(method_name)。",
    "在没有超类的类中使用了super。",
    "方法调用以外不能使用super。",
    "不能引用super的字段。",
    "$(name)是字段，不能覆盖。",
    "重复的字段名$(name)。",
    "数组中没有$(name)方法。",
    "数组中没有$(name)方法。",
    "instanceof的操作数必须是引用类型。",
    "instanceof的右边的类型必须是引用类型。",
    "instanceof的目标必须是类。",
    "instanceof语句一直为真。",
    "instanceof语句一直为假。",
    "因为Diksam的接口间没有父子关系, instanceof语句一直为假。",
    "向下转型的源类型必须是类。",
    "向下转型的目标类型必须是类。",

    "不需要进行向下转型。",
    "尝试将父类转换为子类。",
    "尝试转换没有继承关系的类。",
	"因为Diksam的接口间没有父子关系, 不能向下转型。",
    "不能require文件本身。",
    "if语句的条件表达式不是boolean型。",
    "while语句的条件表达式不是boolean型。",
    "for语句的条件表达式不是boolean型。",
    "do while语句的条件表达式不是boolean型。",

    "被覆盖的方法$(name)的访问修饰符必须比超类的更严格。",
    "方法或函数$(name)的参数数量错误。",
	"方法或函数$(func_name)的第$(index)个参数, $(param_name)的类型错误。",
    "方法或函数$(name)的返回值类型错误。",
    "不能直接调用构造方法。",
    "找不到类型名$(name)。",
    "Diksam的接口之间不能继承（至今为止）。",
    "不能从包外访问成员$(member_name)。",
    "不能从包外访问类$(class_name)。",
    "不能在类外使用this。",
    "不能在类外使用super。",
    "在C样式的注释中终止了文件。",
    "在字符串字面量中终止了文件。",
    "字符字面量中包含了2个以上的字符。",
}
