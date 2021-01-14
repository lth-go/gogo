package vm

import (
	"fmt"
)

func (vm *VirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("print", printProc, 1)
}

func (vm *VirtualMachine) addNativeFunction(funcName string, proc NativeFunctionProc, argCount int) {
	function := &NativeFunction{
		Name:     funcName,
		proc:     proc,
		argCount: argCount,
	}

	vm.functionList = append(vm.functionList, function)
}

func printProc(vm *VirtualMachine, argCount int, args []Object) Value {
	ret := NewIntValue(0)

	str := args[0].(*ObjectString).Value

	fmt.Println(str)

	return ret
}
