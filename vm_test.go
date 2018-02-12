package main

import (
	"testing"

	"./compiler"
	"./vm"
)

func TestNativeFunc(t *testing.T) {
	code := `int print(string str); print("Hello, World!");`

	compiler, err := compiler.ParseSrc(code)
	if err != nil {
		panic(nil)
	}

	exe := vm.NewExecutable()

	compiler.Generate(exe)

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.AddExecutable(exe)

	VM.Execute()
}
