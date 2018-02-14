package main

import (
	"io/ioutil"
	"testing"

	"./compiler"
	"./vm"
)

func TestNativeFunc(t *testing.T) {
	code, err := ioutil.ReadFile("test_single.4g")
	if err != nil {
		panic(err)
	}

	compiler, err := compiler.ParseSrc(string(code))
	if err != nil {
		panic(err)
	}

	compiler.Show()

	exe := vm.NewExecutable()

	compiler.Generate(exe)

	// 打印字节码
	vm.PrintCode(exe.CodeList)

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.AddExecutable(exe)

	VM.Execute()
}
