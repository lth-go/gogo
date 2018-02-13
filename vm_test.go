package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"./compiler"
	"./vm"
)

func TestNativeFunc(t *testing.T) {
	code, err := ioutil.ReadFile("test.4g")
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
	print("=====\n", "code list start\n=====\n")
	for i := 0; i < len(exe.CodeList); {
		code := exe.CodeList[i]
		info := vm.OpcodeInfo[int(code)]
		paramList := []byte(info.Parameter)

		fmt.Println(info.Mnemonic)
		for _, param := range paramList {
			switch param {
			case 'b':
				i += 1
			case 's':
				fallthrough
			case 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
		i += 1
	}
	print("=====\n", "code list end\n=====\n")

	// 创建虚拟机
	VM := vm.NewVirtualMachine()

	VM.AddExecutable(exe)

	VM.Execute()
}
