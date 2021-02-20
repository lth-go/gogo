package main

import (
	"log"
	"os"

	"github.com/lth-go/gogo/compiler"
	"github.com/lth-go/gogo/vm"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("参数错误\n")
	}
	filename := os.Args[1]

	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("文件不存在\n")
	}

	exeList := compiler.CompileFile(filename)

	VM := vm.NewVirtualMachine(exeList)
	VM.Execute()
}
