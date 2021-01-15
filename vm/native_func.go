package vm

import (
	"fmt"
)

func (vm *VirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("_sys", "print", printProc, 1)
}

func (vm *VirtualMachine) addNativeFunction(packageName string, funcName string, proc NativeFunctionProc, argCount int) {
	function := &NativeFunction{
		PackageName: packageName,
		Name:        funcName,
		proc:        proc,
		argCount:    argCount,
	}

	vm.functionList = append(vm.functionList, function)
}

func printProc(vm *VirtualMachine, argCount int, args []Object) Object {
	str := args[0].(*ObjectString).Value

	fmt.Println(str)

	return NilObject
}
