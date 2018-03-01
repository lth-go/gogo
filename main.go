package main

import (
	"os"

	"./compiler"
	"./vm"
)

func main() {
	if len(os.Args) != 2 {
		panic("参数错误")
	}
	filename := os.Args[1]

	_, err := os.Stat(filename)
	if err != nil {
		panic("文件不存在")
	}

	exeList := compiler.CompileFile(filename)

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.SetExecutableList(exeList)

	VM.Execute()
}
