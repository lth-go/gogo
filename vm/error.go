package vm

import (
	//"fmt"
	"fmt"
	"os"
)

const (
    BAD_MULTIBYTE_CHARACTER_ERR int = iota
    FUNCTION_NOT_FOUND_ERR
    FUNCTION_MULTIPLE_DEFINE_ERR
    INDEX_OUT_OF_BOUNDS_ERR
    DIVISION_BY_ZERO_ERR
    NULL_POINTER_ERR
    LOAD_FILE_NOT_FOUND_ERR
    LOAD_FILE_ERR
    CLASS_MULTIPLE_DEFINE_ERR
    CLASS_NOT_FOUND_ERR
    CLASS_CAST_ERR
    DYNAMIC_LOAD_WITHOUT_PACKAGE_ERR
)

var errMessageList []string = []string{
	"不正确的多字节字符。",
	"找不到函数$(name)。",
	"重复定义了函数$(package)#$(name)。",
	"数组下标越界。数组大小为$(size)，访问的下标为[$(index)]。",
	"整数值不能被0除。",
	"引用了null。",
	"没有找到要加载的文件$(file)",
	"加载文件时发生错误($(status))。",
	"重复定义了类$(package)#$(name)。",
	"没有找到类$(name)。",
	"对象的类型为$(org)。,不能向下转型为$(target)。",
	"由于函数$(name)没有指定包，不能动态加载。",
}

var errMessageMap = map[int]string{
//BAD_MULTIBYTE_CHARACTER_ERR:  "不正确的多字节字符。",
//FUNCTION_NOT_FOUND_ERR:       "找不到函数%s。",
//FUNCTION_MULTIPLE_DEFINE_ERR: "重复定义函数%s。",
//INDEX_OUT_OF_BOUNDS_ERR:      "数组下标越界。数组大小为%d，访问的下标为[%d]。",
//DIVISION_BY_ZERO_ERR:         "整数值不能被0除。",
//NULL_POINTER_ERR:             "引用了null。",
}

func vmError(errorNumber int, a ...interface{}) {
	//vm := getVirtualMachine()

	//exe := vm.currentExecutable.executable
	//functionList := vm.currentFunction
	//pc := vm.pc

	fmt.Println("运行错误")
	//fmt.Printf("Line: %d\n", getLineNumberByPc(exe, functionList, pc))
	println(errorNumber)
	println(errMessageList[errorNumber])
	panic("TODO")
	//fmt.Printf(errMessageMap[errorNumber], a...)
	os.Exit(1)
}

func getLineNumberByPc(exe *Executable, function *GFunction, pc int) int {
	var lineNumber []*LineNumber
	var ret int

	if function != nil {
		lineNumber = exe.FunctionList[function.Index].LineNumberList
	} else {
		lineNumber = exe.LineNumberList
	}

	for _, line := range lineNumber {
		if pc >= line.StartPc && pc < (line.StartPc+line.PcCount) {
			ret = line.LineNumber
		}
	}

	return ret
}
