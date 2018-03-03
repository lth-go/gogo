package main

import (
	"testing"

	"./compiler"
	"./vm"
)

var testFile = "test/shape.4g"

func TestVmMachine(t *testing.T) {
	exeList := compiler.CompileFile(testFile)

	// 打印字节码
	for _, exe := range exeList.List {
		vm.PrintCode(exe.CodeList)
	}

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.SetExecutableList(exeList)

	VM.Execute()
}
