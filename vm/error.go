package vm

import (
	"fmt"
	"os"
)

const NO_LINE_NUMBER_PC int = -1

const (
	BAD_MULTIBYTE_CHARACTER_ERR int = iota
	FUNCTION_NOT_FOUND_ERR
	FUNCTION_MULTIPLE_DEFINE_ERR
	INDEX_OUT_OF_BOUNDS_ERR
	DIVISION_BY_ZERO_ERR
	NULL_POINTER_ERR
)

var errMessageMap = map[int]string{
	BAD_MULTIBYTE_CHARACTER_ERR:  "不正确的多字节字符。",
	FUNCTION_NOT_FOUND_ERR:       "找不到函数%s。",
	FUNCTION_MULTIPLE_DEFINE_ERR: "重复定义函数%s。",
	INDEX_OUT_OF_BOUNDS_ERR:      "数组下标越界。数组大小为%d，访问的下标为[%d]。",
	DIVISION_BY_ZERO_ERR:         "整数值不能被0除。",
	NULL_POINTER_ERR:             "引用了null。",
}

func vmError(exe *Executable, function *GFunction, pc int, errorNumber int, a ...interface{}) {
	fmt.Println("编译错误")
	fmt.Printf("Line: %d\n", getLineNumberByPc(exe, function, pc))
	fmt.Printf(errMessageMap[errorNumber], a...)
	fmt.Println("\n")
	os.Exit(1)
}

func getLineNumberByPc(exe *Executable, function *GFunction, pc int) int {
	var lineNumber []*VmLineNumber
	var ret int

	if function != nil {
		lineNumber = exe.FunctionList[function.Index].LineNumberList
	} else {
		lineNumber = exe.LineNumberList
	}

	for _, line := range lineNumber {
		if pc >= line.StartPc && pc < (line.StartPc + line.PcCount) {
			ret = line.LineNumber
		}
	}

	return ret
}
