package main

import (
	"io/ioutil"
	"os"

	"./compiler"
	"./vm"
)

func main() {
	if len(os.Args) != 2 {
		panic("参数错误")
	}
	filename := os.Args[1]

	if !os.IsExist(filename) {
		panic("文件不存在")
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	code := string(buf)

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

	// clean
}
