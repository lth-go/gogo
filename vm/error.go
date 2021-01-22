package vm

import (
	"fmt"
	"log"
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

var errMessageMap map[int]string = map[int]string{
	BAD_MULTIBYTE_CHARACTER_ERR:      "不正确的多字节字符。",
	FUNCTION_NOT_FOUND_ERR:           "找不到函数$(name)。",
	FUNCTION_MULTIPLE_DEFINE_ERR:     "重复定义了函数%s.%s。",
	INDEX_OUT_OF_BOUNDS_ERR:          "数组下标越界。数组大小为$(size)，访问的下标为[$(index)]。",
	DIVISION_BY_ZERO_ERR:             "整数值不能被0除。",
	NULL_POINTER_ERR:                 "引用了null。",
	LOAD_FILE_NOT_FOUND_ERR:          "没有找到要加载的文件$(file)",
	LOAD_FILE_ERR:                    "加载文件时发生错误($(status))。",
	CLASS_MULTIPLE_DEFINE_ERR:        "重复定义了类$(package)#$(name)。",
	CLASS_NOT_FOUND_ERR:              "没有找到类$(name)。",
	CLASS_CAST_ERR:                   "对象的类型为$(org)。,不能向下转型为$(target)。",
	DYNAMIC_LOAD_WITHOUT_PACKAGE_ERR: "由于函数$(name)没有指定包，不能动态加载。",
}

func vmError(errorNumber int, a ...interface{}) {
	//vm := getVirtualMachine()

	//exe := vm.currentExecutable.executable
	//functionList := vm.currentFunction
	//pc := vm.pc

	fmt.Println("运行错误")
	//fmt.Printf("Line: %d\n", getLineNumberByPc(exe, functionList, pc))
	errMsg := fmt.Sprintf(errMessageMap[errorNumber], a...)
	log.Fatalf("%d\n%s", errorNumber, errMsg)
	//fmt.Printf(errMessageMap[errorNumber], a...)
}

func getLineNumberByPc(exe *Executable, function *GoGoFunction, pc int) int {
	var lineNumber []*LineNumber
	var ret int

	// TODO: function == nil
	if function != nil {
		lineNumber = exe.FunctionList[function.Index].LineNumberList
	}

	for _, line := range lineNumber {
		if pc >= line.StartPc && pc < (line.StartPc+line.PcCount) {
			ret = line.LineNumber
		}
	}

	return ret
}
