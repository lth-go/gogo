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

func printProc(vm *VirtualMachine, argCount int, args []Value) Value {
	var str = "null"

	ret := &IntValue{
		intValue: 0,
	}

	obj := args[0].(*ObjectRef).data

	if obj != nil {
		str = obj.(*ObjectString).stringValue
	}

	fmt.Println(str)

	return ret
}
