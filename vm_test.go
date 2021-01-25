package main

import (
	"os"
	"testing"

	"github.com/lth-go/gogo/compiler"
	"github.com/lth-go/gogo/vm"
)

var testFile = "test/test.gogo"

func TestVmMachine(t *testing.T) {
	os.Setenv("IMPORT_SEARCH_PATH", "./test")

	exeList := compiler.CompileFile(testFile)

	// 创建虚拟机
	VM := vm.NewVirtualMachine(exeList)
	VM.Execute()
}
