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

	cm := compiler.NewCompilerManager()

	cm.CompileFile(testFile)

	VM := vm.NewVirtualMachine(
		cm.ConstantList,
		cm.GetVmVariableList(),
		cm.GetVmFunctionList(),
		cm.CodeList,
	)
	VM.Execute()
}
