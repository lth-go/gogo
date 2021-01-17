package vm

import (
	"fmt"
)

// 原生函数
type GoGoNativeFunction struct {
	StaticBase
	proc     NativeFunctionProc
	argCount int
}

type NativeFunctionProc func(vm *VirtualMachine, argCount int, args []Object) Object

// 保存调用函数的索引
type GoGoFunction struct {
	StaticBase
	Executable *Executable
	Index      int
}

func (vm *VirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("_sys", "print", printProc, 1)
}

func (vm *VirtualMachine) addNativeFunction(packageName string, funcName string, proc NativeFunctionProc, argCount int) {
	function := &GoGoNativeFunction{
		StaticBase: StaticBase{
			PackageName: packageName,
			Name:        funcName,
		},
		proc:     proc,
		argCount: argCount,
	}

	vm.static.Append(function)
}

func printProc(vm *VirtualMachine, argCount int, args []Object) Object {
	str := args[0].(*ObjectString).Value

	fmt.Println(str)

	return NilObject
}
