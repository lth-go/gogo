package vm

import (
	"fmt"
)

type Function interface{}

//
// 用户函数
//
type GoGoFunction struct {
	ParamCount   int
	ResultCount  int
	VariableList []Object
	CodeList     []byte
}

//
// 原生函数
//
type GoGoNativeFunction struct {
	ParamCount  int                // 参数数量
	ResultCount int                // 返回值数量
	Proc        NativeFunctionProc // 函数指针
}

type NativeFunctionProc func(vm *VirtualMachine, paramCount int, args []Object) []Object

func (vm *VirtualMachine) AddNativeFunctions() {
	vm.addNativeFunction("_sys", "printf", nativeFuncPrintf, 2, 0)
	vm.addNativeFunction("_sys", "len", nativeFuncLen, 1, 1)
	vm.addNativeFunction("_sys", "append", nativeFuncAppend, 2, 1)
	vm.addNativeFunction("_sys", "delete", nativeFuncDelete, 2, 0)
}

func (vm *VirtualMachine) addNativeFunction(
	packageName string,
	funcName string,
	proc NativeFunctionProc,
	paramCount int,
	resultCount int,
) {
	function := &GoGoNativeFunction{
		Proc:        proc,
		ParamCount:  paramCount,
		ResultCount: resultCount,
	}

	vm.funcList = append(vm.funcList, function)
}

func nativeFuncPrintf(vm *VirtualMachine, paramCount int, args []Object) []Object {
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

func nativeFuncLen(vm *VirtualMachine, paramCount int, args []Object) []Object {
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

func nativeFuncAppend(vm *VirtualMachine, paramCount int, args []Object) []Object {
	obj := args[0].(*ObjectArray)
	arg := args[1].(*ObjectArray)

	obj.List = append(obj.List, arg.List...)

	return []Object{obj}
}

func nativeFuncDelete(vm *VirtualMachine, paramCount int, args []Object) []Object {
	obj := args[0].(*ObjectMap)
	key := args[1]

	obj.Delete(key)

	return nil
}
