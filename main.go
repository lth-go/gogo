package main

import (
	"os"

	"github.com/lth-go/gogo/compiler"
	"github.com/lth-go/gogo/vm"
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

	VM := vm.NewVirtualMachine(exeList)
	VM.Execute()
}
