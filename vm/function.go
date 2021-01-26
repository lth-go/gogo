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
	vm.addNativeFunction("_sys", "printf", nativeFuncPrintf, 2, 0)
	vm.addNativeFunction("_sys", "itoa", nativeFuncItoa, 1, 1)
	vm.addNativeFunction("_sys", "len", nativeFuncLen, 1, 1)
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

func nativeFuncItoa(vm *VirtualMachine, argCount int, args []Object) []Object {
	obj := args[0].(*ObjectInt)

	return []Object{NewObjectString(strconv.Itoa(obj.Value))}
}

func nativeFuncPrintf(vm *VirtualMachine, argCount int, args []Object) []Object {
	format := args[0].(*ObjectString).Value

	switch a := args[1].(type) {
	case *ObjectNil:
		fmt.Printf(format)
	case *ObjectArray:
		list := make([]interface{}, 0)
		for _, valueIFS := range a.List {
			switch value := valueIFS.(type) {
			case *ObjectInt:
				list = append(list, value.Value)
			case *ObjectFloat:
				list = append(list, value.Value)
			case *ObjectString:
				list = append(list, value.Value)
			default:
				panic("TODO")
			}
		}
		fmt.Printf(format, list...)

	default:
		panic("TODO")
	}

	fmt.Printf("")

	return nil
}

func nativeFuncLen(vm *VirtualMachine, argCount int, args []Object) []Object {
	var length int

	switch obj := args[0].(type) {
	case *ObjectString:
		length = len(obj.Value)
	case *ObjectArray:
		length = obj.Len()
	case *ObjectMap:
		length = len(obj.Map)
	default:
		panic("TODO")
	}

	return []Object{NewObjectInt(length)}
}
