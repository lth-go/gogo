package vm

import (
	"fmt"
)

func (vm *VmVirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("print", printProc, 1)
}

func (vm *VmVirtualMachine) addNativeFunction(funcName string, proc VmNativeFunctionProc, argCount int) {
	function := &NativeFunction{
		Name:     funcName,
		proc:     proc,
		argCount: argCount,
	}

	vm.function = append(vm.function, function)
}

func printProc(vm *VmVirtualMachine, argCount int, args []VmValue) VmValue {
	var str = "null"

	ret := &VmIntValue{
		intValue: 0,
	}

	obj := args[0].(*VmObjectValue).objectValue

	if obj != nil {
		str = obj.(*VmObjectString).stringValue
	}
	fmt.Println(str)

	return ret
}
