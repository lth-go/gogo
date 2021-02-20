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

	cm := compiler.GetCurrentCompilerManage()

	VM := vm.NewVirtualMachine(exeList, cm.ConstantList, cm.GetVmVariableList())
	VM.Execute()
}
