package vm

import (
	"fmt"
	"strconv"
)

// 原生函数
type GoGoNativeFunction struct {
	StaticBase
	proc        NativeFunctionProc // 函数指针
	argCount    int                // 参数数量
	resultCount int                // 返回值数量, TODO: 暂时用不上
}

type NativeFunctionProc func(vm *VirtualMachine, argCount int, args []Object) []Object

// 保存调用函数的索引
type GoGoFunction struct {
	StaticBase
	Executable *Executable
	Index      int
}

func (vm *VirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("_sys", "print", nativeFuncPrint, 1, 0)
	vm.addNativeFunction("_sys", "itoa", nativeFuncItoa, 1, 1)
}

func (vm *VirtualMachine) addNativeFunction(
	packageName string,
	funcName string,
	proc NativeFunctionProc,
	argCount int,
	resultCount int,
) {
	function := &GoGoNativeFunction{
		StaticBase: StaticBase{
			PackageName: packageName,
			Name:        funcName,
		},
		proc:        proc,
		argCount:    argCount,
		resultCount: resultCount,
	}

	vm.static.Append(function)
}

func nativeFuncPrint(vm *VirtualMachine, argCount int, args []Object) []Object {
	str := args[0].(*ObjectString).Value

	fmt.Println(str)

	return nil
}

func nativeFuncItoa(vm *VirtualMachine, argCount int, args []Object) []Object {
	obj := args[0].(*ObjectInt)

	return []Object{NewObjectString(strconv.Itoa(obj.Value))}
}
