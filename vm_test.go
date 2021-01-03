package main

import (
	"os"
	"testing"

	"github.com/lth-go/gogogogo/compiler"
	"github.com/lth-go/gogogogo/vm"
)

var testFile = "test/test.4g"

func TestVmMachine(t *testing.T) {
	os.Setenv("IMPORT_SEARCH_PATH", "./test")

	exeList := compiler.CompileFile(testFile)

	//// 打印字节码
	//for _, exe := range exeList.List {
	//	vm.PrintCode(exe.CodeList)
	//}

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.SetExecutableList(exeList)

	VM.Execute()
}
