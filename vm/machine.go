package vm

import (
	"fmt"
)

var functionNotFound = -1
var callFromNative = -1
var vmNullObjectRef = &ObjectRef{}

//
// 虚拟机
//
type VirtualMachine struct {
	// 栈
	stack *Stack
	// 堆
	heap *Heap
	// 程序计数器
	pc int

	// 全局函数列表
	functionList []ExecFunction

	// 当前exe
	currentExecutable *Executable
	// exe列表
	executableList []*Executable

	// 顶层exe
	topLevel *Executable
}

func NewVirtualMachine() *VirtualMachine {
	vm := &VirtualMachine{
		stack:             NewStack(),
		heap:              NewHeap(),
		functionList:      []ExecFunction{},
		currentExecutable: nil,
	}
	// setVirtualMachine(vm)
	vm.AddNativeFunctions()

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

// 添加executableList
func (vm *VirtualMachine) SetExecutableList(exeList *ExecutableList) {
	topExe := exeList.GetTopExe()
	for _, exe := range exeList.List {
		vm.addExecutable(exe, exe == topExe)
	}

}

// 添加单个exe到vm
func (vm *VirtualMachine) addExecutable(exe *Executable, isTopLevel bool) {
	vm.executableList = append(vm.executableList, exe)
	vm.addFunctions(exe)

	// 修正字节码
	// 方法调用修正
	// 函数下标修正
	vm.convertCode(exe, exe.CodeList, nil)
	for _, f := range exe.FunctionList {
		vm.convertCode(exe, f.CodeList, f)
	}

	exe.VariableList.Init()

	if isTopLevel {
		vm.topLevel = exe
	}
}

func (vm *VirtualMachine) addFunctions(exe *Executable) {
	// check func name
	for _, exeFunc := range exe.FunctionList {
		// TODO 实现默认函数后去除
		if !exeFunc.IsImplemented {
			continue
		}
		if vm.searchFunction(exeFunc.PackageName, exeFunc.Name) != functionNotFound {
			vmError(FUNCTION_MULTIPLE_DEFINE_ERR, exeFunc.PackageName, exeFunc.Name)
		}
	}

	// add exe func to vm func
	// TODO: 忽略未实现函数
	for srcIdx, exeFunc := range exe.FunctionList {
		vmFunc := &GFunction{
			PackageName: exeFunc.PackageName,
			Name:        exeFunc.Name,
			Executable:  exe,
			Index:       srcIdx,
		}

		vm.functionList = append(vm.functionList, vmFunc)
	}
}

//
// 虚拟机执行入口
//
func (vm *VirtualMachine) Execute() {
	vm.currentExecutable = vm.topLevel
	// vm.currentFunction = nil
	vm.pc = 0

	vm.stack.expand(vm.topLevel.CodeList)

	vm.execute(nil, vm.topLevel.CodeList)
}

func (vm *VirtualMachine) execute(gFunc *GFunction, codeList []byte) Value {
	var ret Value
	var base int

	stack := vm.stack
	exe := vm.currentExecutable

	for pc := vm.pc; pc < len(codeList); {
		vl := exe.VariableList

		switch codeList[pc] {
		case VM_PUSH_INT_1BYTE:
			stack.setInt(0, int(codeList[pc+1]))
			vm.stack.stackPointer++
			pc += 2
		case VM_PUSH_INT_2BYTE:
			index := get2ByteInt(codeList[pc+1:])
			stack.setInt(0, index)
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setInt(0, exe.ConstantPool.getInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_DOUBLE_0:
			stack.setDouble(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE_1:
			stack.setDouble(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.setDouble(0, exe.ConstantPool.getDouble(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.setObject(0, vm.createStringObject(exe.ConstantPool.getString(index)))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_NULL:
			stack.setObject(0, vmNullObjectRef)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setInt(0, stack.getIntI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.setDouble(0, stack.getDoubleI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setObject(0, stack.getObjectI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setIntI(base+index, stack.getInt(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.setDoubleI(base+index, stack.getDouble(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setObjectI(base+index, stack.getObject(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setInt(0, vl.getInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.setDouble(0, vl.getDouble(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			stack.setObject(0, vl.getObject(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			vl.setInt(index, stack.getInt(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			vl.setDouble(index, stack.getDouble(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			vl.setObject(index, stack.getObject(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_ARRAY_INT:
			array := stack.getArrayInt(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			intValue := array.getInt(index)

			stack.setInt(-2, intValue)
			vm.stack.stackPointer--
			pc++
		case VM_PUSH_ARRAY_DOUBLE:
			array := stack.getArrayDouble(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			doubleValue := array.getDouble(index)

			stack.setDouble(-2, doubleValue)
			vm.stack.stackPointer--
			pc++
		case VM_PUSH_ARRAY_OBJECT:
			array := stack.getArrayObject(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			object := array.getObject(index)

			stack.setObject(-2, object)
			vm.stack.stackPointer--
			pc++
		case VM_POP_ARRAY_INT:
			value := stack.getInt(-3)
			array := stack.getArrayInt(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			array.setInt(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_POP_ARRAY_DOUBLE:
			value := stack.getDouble(-3)
			array := stack.getArrayDouble(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			array.setDouble(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_POP_ARRAY_OBJECT:
			value := stack.getObject(-3)
			array := stack.getArrayObject(-2)
			index := stack.getInt(-1)

			vm.restorePc(exe, gFunc, pc)
			array.setObject(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_ADD_INT:
			stack.setInt(-2, stack.getInt(-2)+stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_DOUBLE:
			stack.setDouble(-2, stack.getDouble(-2)+stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_STRING:
			stack.setObject(-2, vm.chainStringObject(stack.getObject(-2), stack.getObject(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_INT:
			stack.setInt(-2, stack.getInt(-2)-stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_DOUBLE:
			stack.setDouble(-2, stack.getDouble(-2)-stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_INT:
			stack.setInt(-2, stack.getInt(-2)*stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_DOUBLE:
			stack.setDouble(-2, stack.getDouble(-2)*stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_INT:
			if stack.getInt(-1) == 0 {
				vmError(DIVISION_BY_ZERO_ERR)
			}
			stack.setInt(-2, stack.getInt(-2)/stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_DOUBLE:
			stack.setDouble(-2, stack.getDouble(-2)/stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MINUS_INT:
			stack.setInt(-1, -stack.getInt(-1))
			pc++
		case VM_MINUS_DOUBLE:
			stack.setDouble(-1, -stack.getDouble(-1))
			pc++
		case VM_CAST_INT_TO_DOUBLE:
			stack.setDouble(-1, float64(stack.getInt(-1)))
			pc++
		case VM_CAST_DOUBLE_TO_INT:
			stack.setInt(-1, int(stack.getDouble(-1)))
			pc++
		case VM_CAST_BOOLEAN_TO_STRING:
			if stack.getInt(-1) != 0 {
				stack.setObject(-1, vm.createStringObject("true"))
			} else {
				stack.setObject(-1, vm.createStringObject("false"))
			}
			pc++
		case VM_CAST_INT_TO_STRING:
			// TODO 啥意思
			vm.restorePc(exe, gFunc, pc)
			buf := fmt.Sprintf("%d", stack.getInt(-1))
			stack.setObject(-1, vm.createStringObject(buf))
			pc++
		case VM_CAST_DOUBLE_TO_STRING:
			// TODO 啥意思
			vm.restorePc(exe, gFunc, pc)
			buf := fmt.Sprintf("%f", stack.getDouble(-1))
			stack.setObject(-1, vm.createStringObject(buf))
			pc++
		case VM_EQ_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) == stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) == stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_OBJECT:
			stack.setInt(-2, boolToInt(stack.getObject(-2).data == stack.getObject(-1).data))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			stack.setInt(-2, boolToInt(!(stack.getString(-2) == stack.getString(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_GT_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) > stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) > stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_STRING:
			stack.setInt(-2, boolToInt(stack.getString(-2) > stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) >= stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) >= stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_STRING:
			stack.setInt(-2, boolToInt(stack.getString(-2) >= stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) < stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) < stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_STRING:
			stack.setInt(-2, boolToInt(stack.getString(-2) < stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) <= stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) <= stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_STRING:
			stack.setInt(-2, boolToInt(stack.getString(-2) <= stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) != stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) != stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_OBJECT:
			stack.setInt(-2, boolToInt(stack.getObject(-2).data != stack.getObject(-1).data))
			vm.stack.stackPointer--
			pc++
		case VM_NE_STRING:
			stack.setInt(-2, boolToInt(stack.getString(-2) != stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_AND:
			stack.setInt(-2, boolToInt(intToBool(stack.getInt(-2)) && intToBool(stack.getInt(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_OR:
			stack.setInt(-2, boolToInt(intToBool(stack.getInt(-2)) || intToBool(stack.getInt(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_NOT:
			stack.setInt(-1, boolToInt(!intToBool(stack.getInt(-1))))
			pc++
		case VM_POP:
			vm.stack.stackPointer--
			pc++
		case VM_DUPLICATE:
			// TODO
			stack.stack[vm.stack.stackPointer] = stack.stack[vm.stack.stackPointer-1]
			stack.stack[vm.stack.stackPointer].setPointer(stack.stack[vm.stack.stackPointer-1].isPointer())
			vm.stack.stackPointer++
			pc++
		case VM_DUPLICATE_OFFSET:
			offset := get2ByteInt(codeList[pc+1:])
			stack.stack[vm.stack.stackPointer] = stack.stack[vm.stack.stackPointer-1-offset]
			vm.stack.stackPointer++
			pc += 3
		case VM_JUMP:
			index := get2ByteInt(codeList[pc+1:])
			pc = index
		case VM_JUMP_IF_TRUE:
			if intToBool(stack.getInt(-1)) {
				index := get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_JUMP_IF_FALSE:
			if !intToBool(stack.getInt(-1)) {
				index := get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_PUSH_FUNCTION:
			value := get2ByteInt(codeList[pc+1:])
			stack.setInt(0, value)
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			funcIdx := stack.getInt(-1)
			switch f := vm.functionList[funcIdx].(type) {
			case *NativeFunction:
				vm.restorePc(exe, gFunc, pc)
				vm.InvokeNativeFunction(f, &vm.stack.stackPointer)
				pc++
			case *GFunction:
				vm.InvokeFunction(&gFunc, f, &codeList, &pc, &vm.stack.stackPointer, &base, &exe)
			default:
				panic("TODO")
			}
		case VM_RETURN:
			if vm.returnFunction(&gFunc, &codeList, &pc, &base, &exe) {
				ret = stack.stack[stack.stackPointer-1]
				// TODO goto
				return ret
			}
		case VM_NEW_ARRAY_LITERAL_INT:
			size := get2ByteInt(codeList[pc+1:])

			vm.restorePc(exe, gFunc, pc)
			array := vm.createArrayLiteralInt(size)
			vm.stack.stackPointer -= size
			stack.setObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		case VM_NEW_ARRAY_LITERAL_DOUBLE:
			size := get2ByteInt(codeList[pc+1:])

			vm.restorePc(exe, gFunc, pc)
			array := vm.createArrayLiteralFloat(size)
			vm.stack.stackPointer -= size
			stack.setObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		case VM_NEW_ARRAY_LITERAL_OBJECT:
			size := get2ByteInt(codeList[pc+1:])

			vm.restorePc(exe, gFunc, pc)
			array := vm.createArrayLiteralObject(size)
			vm.stack.stackPointer -= size
			stack.setObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		default:
			panic("TODO")
		}
	}
	return ret
}

func (vm *VirtualMachine) initializeLocalVariables(f *Function, fromSp int) {
	var i, spIdx int

	for i = 0; i < len(f.LocalVariableList); i++ {
		vm.stack.stack[i].setPointer(false)
	}

	spIdx = fromSp
	for i = 0; i < len(f.LocalVariableList); i++ {
		vm.stack.stack[spIdx] = initializeValue(f.LocalVariableList[i].TypeSpecifier)

		if f.LocalVariableList[i].TypeSpecifier.IsReferenceType() {
			vm.stack.stack[i].setPointer(true)
		}
		spIdx++
	}
}

// 修正转换code
func (vm *VirtualMachine) convertCode(exe *Executable, codeList []byte, f *Function) {
	var destIdx int

	for i := 0; i < len(codeList); i++ {
		code := codeList[i]
		switch code {
		// 函数内的本地声明
		case VM_PUSH_STACK_INT, VM_POP_STACK_INT,
			VM_PUSH_STACK_DOUBLE, VM_POP_STACK_DOUBLE,
			VM_PUSH_STACK_OBJECT, VM_POP_STACK_OBJECT:

			var parameterCount int

			if f == nil {
				panic("can't find line, need debug!!!")
			}

			if f.IsMethod {
				parameterCount = len(f.ParameterList) + 1 /* for this */
			} else {
				parameterCount = len(f.ParameterList)
			}

			// 增加返回值的位置
			srcIdx := get2ByteInt(codeList[i+1:])
			if srcIdx >= parameterCount {
				destIdx = srcIdx + 1
			} else {
				destIdx = srcIdx
			}
			set2ByteInt(codeList[i+1:], destIdx)

		case VM_PUSH_FUNCTION:
			idxInExe := get2ByteInt(codeList[i+1:])
			funcIdx := vm.searchFunction(exe.FunctionList[idxInExe].PackageName, exe.FunctionList[idxInExe].Name)
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
func (vm *VirtualMachine) searchFunction(packageName, name string) int {

	for i, f := range vm.functionList {
		if f.getPackageName() == packageName && f.getName() == name {
			return i
		}
	}
	return functionNotFound
}

//
// 函数相关
//
// 执行原生函数
func (vm *VirtualMachine) InvokeNativeFunction(f *NativeFunction, spP *int) {

	stack := vm.stack.stack
	sp := *spP

	ret := f.proc(vm, f.argCount, stack[sp-f.argCount-1:])

	stack[sp-f.argCount-1] = ret
	*spP = sp - f.argCount
}

// 函数执行
func (vm *VirtualMachine) InvokeFunction(caller **GFunction, callee *GFunction, codeP *[]byte, pcP *int, spP *int, baseP *int, exe **Executable) {
	// caller 调用者, 当前所属的函数调用域
	// callee 要调用的函数的基本信息

	*exe = callee.Executable

	// 包含调用函数的全部信息
	calleeP := (*exe).FunctionList[callee.Index]

	// 拓展栈大小
	vm.stack.expand(calleeP.CodeList)

	// 设置返回值信息
	callInfo := &CallInfo{
		caller:        *caller,
		callerAddress: *pcP,
		base:          *baseP,
	}

	// 栈上保存返回信息
	vm.stack.stack[*spP-1] = callInfo

	// 设置base
	*baseP = *spP - len(calleeP.ParameterList) - 1
	if calleeP.IsMethod {
		*baseP--
	}

	// 设置调用者
	*caller = callee

	// 初始化参数
	vm.initializeLocalVariables(calleeP, *spP)

	// 设置栈位置
	*spP += len(calleeP.LocalVariableList)
	*pcP = 0

	// 设置字节码为函数的字节码
	*codeP = calleeP.CodeList
}

func (vm *VirtualMachine) returnFunction(funcP **GFunction, codeP *[]byte, pcP *int, baseP *int, exe **Executable) bool {

	returnValue := vm.stack.stack[vm.stack.stackPointer-1]
	vm.stack.stackPointer--

	calleeFunc := (*exe).FunctionList[(*funcP).Index]

	ret := doReturn(vm, funcP, codeP, pcP, baseP, exe)

	vm.stack.stack[vm.stack.stackPointer] = returnValue
	vm.stack.stack[vm.stack.stackPointer].setPointer(calleeFunc.TypeSpecifier.IsReferenceType())
	vm.stack.stackPointer++

	return ret
}

func doReturn(vm *VirtualMachine, funcP **GFunction, codeP *[]byte, pcP *int, baseP *int, exeP **Executable) bool {

	calleeP := (*exeP).FunctionList[(*funcP).Index]

	argCount := len(calleeP.ParameterList)
	// method
	if calleeP.IsMethod {
		argCount++
	}
	callInfo := vm.stack.stack[*baseP+argCount].(*CallInfo)

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

	return callInfo.callerAddress == callFromNative
}

func (vm *VirtualMachine) createArrayLiteralInt(size int) *ObjectRef {

	array := vm.createArrayInt(size)
	for i := 0; i < size; i++ {
		array.data.(*ObjectArrayInt).intArray[i] = vm.stack.getInt(-size + i)
	}

	return array
}

func (vm *VirtualMachine) createArrayLiteralFloat(size int) *ObjectRef {

	array := vm.createArrayDouble(size)
	for i := 0; i < size; i++ {
		array.data.(*ObjectArrayDouble).doubleArray[i] = vm.stack.getDouble(-size + i)
	}

	return array
}

func (vm *VirtualMachine) createArrayLiteralObject(size int) *ObjectRef {
	array := vm.createArrayObject(size)
	for i := 0; i < size; i++ {
		array.data.(*ObjectArrayObject).objectArray[i] = vm.stack.getObject(-size + i)
	}

	return array
}

func (vm *VirtualMachine) restorePc(ee *Executable, function *GFunction, pc int) {
	vm.currentExecutable = ee
	// vm.currentFunction = function
	vm.pc = pc
}
