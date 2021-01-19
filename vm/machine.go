package vm

import (
	"fmt"
)

var functionNotFound = -1
var callFromNative = -1
var NilObject = &ObjectNil{}

//
// 虚拟机
//
type VirtualMachine struct {
	stack             *Stack        // 栈
	heap              *Heap         // 堆
	static            *Static       // 静态区
	pc                int           // 程序计数器
	currentExecutable *Executable   // 当前exe
	executableList    []*Executable // exe列表
	topLevel          *Executable   // 顶层exe
}

func NewVirtualMachine(exeList *ExecutableList) *VirtualMachine {
	vm := &VirtualMachine{
		stack:             NewStack(),
		heap:              NewHeap(),
		static:            NewStatic(),
		currentExecutable: nil,
	}

	// setVirtualMachine(vm)

	// 添加原生函数
	vm.AddNativeFunctions()

	for _, exe := range exeList.List {
		vm.AddExecutable(exe)
	}

	vm.SetTopExe(exeList.GetTopExe())
	vm.SetMainEntrypoint()

	return vm
}

// 设置全局vm
// var StVirtualMachine *VirtualMachine

// func setVirtualMachine(vm *VirtualMachine) {
//     StVirtualMachine = vm
// }
// func getVirtualMachine() *VirtualMachine {
//     return StVirtualMachine
// }

//
// 虚拟机初始化操作
//

func (vm *VirtualMachine) SetTopExe(exe *Executable) {
	vm.topLevel = exe
}

func (vm *VirtualMachine) SetMainEntrypoint() {
	// TODO: 设置入口为main函数
	// TODO: packageName 是main
	idx := vm.SearchStatic("", "main")
	if idx == -1 {
		panic("TODO")
	}
	b := make([]byte, 2)
	set2ByteInt(b, idx)
	vm.topLevel.CodeList = append(vm.topLevel.CodeList, b...)
	vm.topLevel.CodeList = append(vm.topLevel.CodeList, VM_INVOKE)
}

// 添加单个exe到vm
func (vm *VirtualMachine) AddExecutable(exe *Executable) {
	vm.executableList = append(vm.executableList, exe)

	vm.AddFunctions(exe)
	vm.AddStatic(exe)

	// 修正字节码
	// 方法调用修正
	// 函数下标修正
	vm.ConvertOpCode(exe, exe.CodeList, nil)
	for _, f := range exe.FunctionList {
		vm.ConvertOpCode(exe, f.CodeList, f)
	}
}

// 添加静态区
func (vm *VirtualMachine) AddStatic(exe *Executable) {
	// 变量初始化
	exe.VariableList.Init()

	for _, value := range exe.VariableList.VariableList {
		vm.static.Append(NewStaticVariable(exe.PackageName, value.Name, value.Value))
	}
}

// 添加exe函数到虚拟机
func (vm *VirtualMachine) AddFunctions(exe *Executable) {
	// 检查函数是否重复定义
	for _, exeFunc := range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}
		if vm.SearchStatic(exeFunc.PackageName, exeFunc.Name) != functionNotFound {
			vmError(FUNCTION_MULTIPLE_DEFINE_ERR, exeFunc.PackageName, exeFunc.Name)
		}
	}

	for srcIdx, exeFunc := range exe.FunctionList {
		// 跳过原生,其他包函数
		if !exeFunc.IsImplemented {
			continue
		}
		vmFunc := &GoGoFunction{
			StaticBase: StaticBase{
				PackageName: exeFunc.PackageName,
				Name:        exeFunc.Name,
			},
			Executable: exe,
			Index:      srcIdx,
		}

		vm.static.Append(vmFunc)
	}
}

//
// 虚拟机执行入口
//
func (vm *VirtualMachine) Execute() {
	vm.currentExecutable = vm.topLevel
	vm.pc = 0
	vm.stack.Expand(vm.topLevel.CodeList)
	vm.execute(nil, vm.topLevel.CodeList)
}

func (vm *VirtualMachine) execute(gogoFunc *GoGoFunction, codeList []byte) Object {
	var ret Object
	var base int

	stack := vm.stack
	exe := vm.currentExecutable
	static := vm.static

	for pc := vm.pc; pc < len(codeList); {
		switch codeList[pc] {
		case VM_PUSH_INT_1BYTE:
			stack.SetIntPlus(0, int(codeList[pc+1]))
			vm.stack.stackPointer++
			pc += 2
		case VM_PUSH_INT_2BYTE:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, index)
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, exe.ConstantPool.GetInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_FLOAT_0:
			stack.SetFloatPlus(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_FLOAT_1:
			stack.SetFloatPlus(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_FLOAT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetFloatPlus(0, exe.ConstantPool.GetFloat(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetStringPlus(0, exe.ConstantPool.GetString(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_NIL:
			stack.SetPlus(0, NilObject)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, stack.GetInt(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_FLOAT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetFloatPlus(0, stack.GetFloat(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetPlus(0, stack.Get(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetInt(base+index, stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_FLOAT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetFloat(base+index, stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.Set(base+index, stack.GetPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, static.GetVariableInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_FLOAT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetFloatPlus(0, static.GetVariableFloat(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.SetPlus(0, static.GetVariableObject(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			static.SetVariable(index, stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_FLOAT:
			index := get2ByteInt(codeList[pc+1:])
			static.SetVariable(index, stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			static.SetVariable(index, stack.GetPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_ARRAY_OBJECT:
			array := stack.GetArrayPlus(-2)
			index := stack.GetIntPlus(-1)

			object := array.Get(index)

			stack.SetPlus(-2, object)
			vm.stack.stackPointer--
			pc++
		case VM_POP_ARRAY_OBJECT:
			value := stack.GetPlus(-3)
			array := stack.GetArrayPlus(-2)
			index := stack.GetIntPlus(-1)

			array.Set(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_ADD_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)+stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)+stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_STRING:
			stack.SetStringPlus(-2, stack.GetStringPlus(-2)+stack.GetStringPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)-stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)-stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)*stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)*stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_INT:
			if stack.GetIntPlus(-1) == 0 {
				vmError(DIVISION_BY_ZERO_ERR)
			}
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)/stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)/stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MINUS_INT:
			stack.SetIntPlus(-1, -stack.GetIntPlus(-1))
			pc++
		case VM_MINUS_FLOAT:
			stack.SetFloatPlus(-1, -stack.GetFloatPlus(-1))
			pc++
		case VM_CAST_INT_TO_FLOAT:
			stack.SetFloatPlus(-1, float64(stack.GetIntPlus(-1)))
			pc++
		case VM_CAST_FLOAT_TO_INT:
			stack.SetIntPlus(-1, int(stack.GetFloatPlus(-1)))
			pc++
		case VM_CAST_BOOLEAN_TO_STRING:
			if stack.GetIntPlus(-1) != 0 {
				stack.SetStringPlus(-1, "true")
			} else {
				stack.SetStringPlus(-1, "false")
			}
			pc++
		case VM_CAST_INT_TO_STRING:
			buf := fmt.Sprintf("%d", stack.GetIntPlus(-1))
			stack.SetStringPlus(-1, buf)
			pc++
		case VM_CAST_FLOAT_TO_STRING:
			buf := fmt.Sprintf("%f", stack.GetFloatPlus(-1))
			stack.SetStringPlus(-1, buf)
			pc++
		case VM_EQ_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) == stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) == stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_OBJECT:
			stack.SetIntPlus(-2, boolToInt(stack.GetPlus(-2) == stack.GetPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			stack.SetIntPlus(-2, boolToInt(!(stack.GetStringPlus(-2) == stack.GetStringPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_GT_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) > stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) > stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_STRING:
			stack.SetIntPlus(-2, boolToInt(stack.GetStringPlus(-2) > stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) >= stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) >= stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_STRING:
			stack.SetIntPlus(-2, boolToInt(stack.GetStringPlus(-2) >= stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) < stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) < stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_STRING:
			stack.SetIntPlus(-2, boolToInt(stack.GetStringPlus(-2) < stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) <= stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) <= stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_STRING:
			stack.SetIntPlus(-2, boolToInt(stack.GetStringPlus(-2) <= stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_INT:
			stack.SetIntPlus(-2, boolToInt(stack.GetIntPlus(-2) != stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_FLOAT:
			stack.SetIntPlus(-2, boolToInt(stack.GetFloatPlus(-2) != stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_OBJECT:
			stack.SetIntPlus(-2, boolToInt(stack.GetPlus(-2) != stack.GetPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_STRING:
			stack.SetIntPlus(-2, boolToInt(stack.GetStringPlus(-2) != stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_AND:
			stack.SetIntPlus(-2, boolToInt(intToBool(stack.GetIntPlus(-2)) && intToBool(stack.GetIntPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_OR:
			stack.SetIntPlus(-2, boolToInt(intToBool(stack.GetIntPlus(-2)) || intToBool(stack.GetIntPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_NOT:
			stack.SetIntPlus(-1, boolToInt(!intToBool(stack.GetIntPlus(-1))))
			pc++
		case VM_POP:
			vm.stack.stackPointer--
			pc++
		case VM_DUPLICATE:
			stack.Set(vm.stack.stackPointer, stack.Get(vm.stack.stackPointer-1))
			vm.stack.stackPointer++
			pc++
		case VM_DUPLICATE_OFFSET:
			offset := get2ByteInt(codeList[pc+1:])
			stack.Set(vm.stack.stackPointer, stack.Get(vm.stack.stackPointer-1-offset))
			vm.stack.stackPointer++
			pc += 3
		case VM_JUMP:
			index := get2ByteInt(codeList[pc+1:])
			pc = index
		case VM_JUMP_IF_TRUE:
			if intToBool(stack.GetIntPlus(-1)) {
				index := get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_JUMP_IF_FALSE:
			if !intToBool(stack.GetIntPlus(-1)) {
				index := get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_PUSH_FUNCTION:
			value := get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, value)
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			funcIdx := stack.GetIntPlus(-1)
			switch f := vm.static.Get(funcIdx).(type) {
			case *GoGoNativeFunction:
				vm.restorePc(exe, gogoFunc, pc)
				vm.InvokeNativeFunction(f, &vm.stack.stackPointer)
				pc++
			case *GoGoFunction:
				vm.InvokeFunction(&gogoFunc, f, &codeList, &pc, &vm.stack.stackPointer, &base, &exe)
			default:
				panic("TODO")
			}
		case VM_RETURN:
			if vm.returnFunction(&gogoFunc, &codeList, &pc, &base, &exe) {
				// TODO: 目前执行不进去
				ret = stack.Get(stack.stackPointer - 1)
				return ret
			}
		case VM_NEW_ARRAY:
			size := get2ByteInt(codeList[pc+1:])
			array := vm.NewObjectArray(size)

			vm.stack.stackPointer -= size
			stack.SetPlus(0, array)
			vm.stack.stackPointer++
			pc += 3
		default:
			panic("TODO")
		}
	}
	return ret
}

func (vm *VirtualMachine) initLocalVariables(f *Function, fromSp int) {
	var i, spIdx int

	spIdx = fromSp
	for i = 0; i < len(f.LocalVariableList); i++ {
		vm.stack.Set(spIdx, GetObjectByType(f.LocalVariableList[i].Type))
		spIdx++
	}
}

// 修正转换code
func (vm *VirtualMachine) ConvertOpCode(exe *Executable, codeList []byte, f *Function) {
	var destIdx int

	for i := 0; i < len(codeList); i++ {
		code := codeList[i]
		switch code {
		// 函数内的本地声明
		case VM_PUSH_STACK_INT, VM_POP_STACK_INT,
			VM_PUSH_STACK_FLOAT, VM_POP_STACK_FLOAT,
			VM_PUSH_STACK_OBJECT, VM_POP_STACK_OBJECT:

			var parameterCount int

			if f == nil {
				panic("can't find line, need debug!!!")
			}

			parameterCount = f.GetParamCount()
			if f.IsMethod {
				parameterCount += 1 /* for this */
			}

			// 增加返回值的位置
			srcIdx := get2ByteInt(codeList[i+1:])
			if srcIdx >= parameterCount {
				destIdx = srcIdx + 1
			} else {
				destIdx = srcIdx
			}
			set2ByteInt(codeList[i+1:], destIdx)
		case VM_PUSH_STATIC_INT, VM_PUSH_STATIC_FLOAT, VM_PUSH_STATIC_OBJECT,
			VM_POP_STATIC_INT, VM_POP_STATIC_FLOAT, VM_POP_STATIC_OBJECT:

			idxInExe := get2ByteInt(codeList[i+1:])
			funcIdx := vm.SearchStatic(exe.PackageName, exe.VariableList.VariableList[idxInExe].Name)
			set2ByteInt(codeList[i+1:], funcIdx)

		case VM_PUSH_FUNCTION:
			idxInExe := get2ByteInt(codeList[i+1:])
			funcIdx := vm.SearchStatic(exe.FunctionList[idxInExe].PackageName, exe.FunctionList[idxInExe].Name)
			set2ByteInt(codeList[i+1:], funcIdx)
		}

		info := &OpcodeInfo[code]
		for _, p := range []byte(info.Parameter) {
			switch p {
			case 'b':
				i++
			case 's', 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
	}
}

// 查找函数
func (vm *VirtualMachine) SearchStatic(packageName, name string) int {
	return vm.static.Index(packageName, name)
}

//
// 函数相关
//
// 执行原生函数
func (vm *VirtualMachine) InvokeNativeFunction(f *GoGoNativeFunction, spP *int) {
	sp := *spP

	ret := f.proc(vm, f.argCount, vm.stack.objectList[sp-f.argCount-1:])

	vm.stack.Set(sp-f.argCount-1, ret)

	*spP = sp - f.argCount
}

// 函数执行
func (vm *VirtualMachine) InvokeFunction(caller **GoGoFunction, callee *GoGoFunction, codeP *[]byte, pcP *int, spP *int, baseP *int, exe **Executable) {
	// caller 调用者, 当前所属的函数调用域
	// callee 要调用的函数的基本信息

	*exe = callee.Executable

	// 包含调用函数的全部信息
	calleeP := (*exe).FunctionList[callee.Index]

	// 拓展栈大小
	vm.stack.Expand(calleeP.CodeList)

	// 设置返回值信息
	callInfo := &ObjectCallInfo{
		caller:        *caller,
		callerAddress: *pcP,
		base:          *baseP,
	}

	// 栈上保存返回信息
	vm.stack.Set(*spP-1, callInfo)

	// 设置base
	*baseP = *spP - calleeP.GetParamCount() - 1
	if calleeP.IsMethod {
		*baseP--
	}

	// 设置调用者
	*caller = callee

	// 初始化参数
	vm.initLocalVariables(calleeP, *spP)

	// 设置栈位置
	*spP += len(calleeP.LocalVariableList)
	*pcP = 0

	// 设置字节码为函数的字节码
	*codeP = calleeP.CodeList
}

// 保存返回值,并恢复栈
func (vm *VirtualMachine) returnFunction(funcP **GoGoFunction, codeP *[]byte, pcP *int, baseP *int, exeP **Executable) bool {
	// calleeP := (*exeP).FunctionList[(*funcP).Index]
	// argCount := len(calleeP.ParameterList)

	// 获取返回值,用于恢复
	returnValue := vm.stack.Get(vm.stack.stackPointer - 1)
	vm.stack.stackPointer--

	// 恢复调用栈
	ret := doReturn(vm, funcP, codeP, pcP, baseP, exeP)

	vm.stack.Set(vm.stack.stackPointer, returnValue)
	vm.stack.stackPointer++

	return ret
}

// 恢复到父调用栈
func doReturn(vm *VirtualMachine, funcP **GoGoFunction, codeP *[]byte, pcP *int, baseP *int, exeP **Executable) bool {

	calleeP := (*exeP).FunctionList[(*funcP).Index]
	argCount := calleeP.GetParamCount()

	if calleeP.IsMethod {
		argCount++
	}

	callInfo := vm.stack.Get(*baseP + argCount).(*ObjectCallInfo)

	if callInfo.caller != nil {
		*exeP = callInfo.caller.Executable
		callerP := (*exeP).FunctionList[callInfo.caller.Index]
		*codeP = callerP.CodeList
	} else {
		*exeP = vm.topLevel
		*codeP = vm.topLevel.CodeList
	}

	*funcP = callInfo.caller
	vm.stack.stackPointer = *baseP
	*pcP = callInfo.callerAddress + 1
	*baseP = callInfo.base

	// TODO: 目前无效
	return callInfo.callerAddress == callFromNative
}

func (vm *VirtualMachine) restorePc(ee *Executable, function *GoGoFunction, pc int) {
	vm.currentExecutable = ee
	vm.pc = pc
}
