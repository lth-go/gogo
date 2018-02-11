package vm

import (
	"fmt"
)

//
// 虚拟机
//
type VmVirtualMachine struct {
	// 栈
	stack Stack
	// 堆
	heap Heap
	// 全局变量
	static Static
	// 全局函数列表
	function []Function
	// 解释器
	executable *Executable
	// 程序计数器
	pc int
}

func NewVirtualMachine() *VmVirtualMachine {
	vm := &VmVirtualMachine{
		stack:      NewStack(),
		heap:       NewHeap(),
		static:     NewStatic(),
		function:   []Function{},
		executable: nil,
	}
	vm.AddNativeFunctions()

	return vm
}

//
// 一些初始化操作
//

// 虚拟机添加解释器
func (vm *VmVirtualMachine) AddExecutable(exe *Executable) {

	vm.executable = exe

	vm.addFunctions(exe)

	vm.convertCode(exe, exe.CodeList, nil)

	for _, f := range exe.FunctionList {
		vm.convertCode(exe, f.CodeList, f)
	}

	vm.addStaticVariables(exe)
}

func (vm *VmVirtualMachine) addStaticVariables(exe *Executable) {

	for _, exeValue := range exe.GlobalVariableList {
		newVmValue := vm.initializeValue(exeValue.typeSpecifier.BasicType)
		vm.static.append(newVmValue)
	}
}

func (vm *VmVirtualMachine) addFunctions(exe *Executable) {

	for _, exeFunc := range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}
		// 不能添加重名函数
		for _, vmFunc := range vm.function {
			if vmFunc.getName() == exeFunc.Name {
				panic("TODO")
			}
		}
	}

	for srcIdex, exeFunc := range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}
		newVmFunc := &GFunction{Name: exeFunc.Name, Executable: exe, Index: srcIdex}
		vm.function = append(vm.function, newVmFunc)
	}
}

func (vm *VmVirtualMachine) AddNativeFunctions() {
	// TODO
	//vm.addNativeFunction(vm, "print", nv_print_proc, 1)
}

//
// 虚拟机执行入口
//
func (vm *VmVirtualMachine) Execute() {
	vm.pc = 0

	vm.execute(nil, vm.executable.CodeList)
}

func (vm *VmVirtualMachine) execute(gFunc *GFunction, codeList []byte) {
	var base int

	stack := vm.stack
	static := vm.stack
	exe := vm.executable

	for pc := vm.pc; pc < len(codeList); {

		switch codeList[pc] {
		case VM_PUSH_INT_1BYTE:
			stack.writeInt(0, int(codeList[pc+1]))
			vm.stack.stackPointer++
			pc += 2
		case VM_PUSH_INT_2BYTE:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, index)
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, exe.ConstantPool.getInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_DOUBLE_0:
			stack.writeDouble(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE_1:
			stack.writeDouble(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeDouble(0, exe.ConstantPool.getDouble(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeObject(0, vm.createStringObject(exe.ConstantPool.getString(index)))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, stack.getIntI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeDouble(0, stack.getDoubleI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeObject(0, stack.getObjectI(base+index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STACK_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeIntI(base+index, stack.getInt(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeDoubleI(base+index, stack.getDouble(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeObjectI(base+index, stack.getObject(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, static.getInt(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeDouble(0, static.getDouble(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_STRING:
			index := get2ByteInt(codeList[pc+1:])
			stack.writeObject(0, static.getObject(index))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STATIC_INT:
			index := get2ByteInt(codeList[pc+1:])
			static.setInt(index, stack.getInt(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_DOUBLE:
			index := get2ByteInt(codeList[pc+1:])
			static.setDouble(index, stack.getDouble(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_STRING:
			index := get2ByteInt(codeList[pc+1:])
			static.setObject(index, stack.getObject(-1))
			vm.stack.stackPointer--
			pc += 3
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
				stack.writeObject(-1, vm.createStringObject("true"))
			} else {
				stack.writeObject(-1, vm.createStringObject("false"))
			}
			pc++
		case VM_CAST_INT_TO_STRING:
			buf := fmt.Sprintf("%d", stack.getInt(-1))
			stack.writeObject(-1, vm.createStringObject(buf))
			pc++
		case VM_CAST_DOUBLE_TO_STRING:
			buf := fmt.Sprintf("%f", stack.getDouble(-1))
			stack.writeObject(-1, vm.createStringObject(buf))
			pc++
		case VM_EQ_INT:
			stack.setInt(-2, boolToInt(stack.getInt(-2) == stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_DOUBLE:
			stack.setInt(-2, boolToInt(stack.getDouble(-2) == stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			stack.writeInt(-2, boolToInt(!(stack.getObject(-2).getString() == stack.getObject(-1).getString())))
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
			stack.writeInt(-2, boolToInt(stack.getObject(-2).getString() > stack.getObject(-1).getString()))
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
			stack.writeInt(-2, boolToInt(stack.getObject(-2).getString() >= stack.getObject(-1).getString()))
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
			stack.writeInt(-2, boolToInt(stack.getObject(-2).getString() < stack.getObject(-1).getString()))
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
			stack.writeInt(-2, boolToInt(stack.getObject(-2).getString() <= stack.getObject(-1).getString()))
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
		case VM_NE_STRING:
			stack.writeInt(-2, boolToInt(stack.getObject(-2).getString() != stack.getObject(-1).getString()))
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
			vm.stack.stackPointer++
			pc++
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
			index := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, index)
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			func_idx := stack.getInt(-1)
			switch f := vm.function[func_idx].(type) {
			case *NativeFunction:
				vm.invokeNativeFunction(f, &vm.stack.stackPointer)
				pc++
			case *GFunction:
				vm.invokeGFunction(gFunc, f, codeList, &pc, &vm.stack.stackPointer, &base, exe)
			default:
				panic("TODO")
			}
		case VM_RETURN:
			vm.returnFunction(gFunc, codeList, &pc, &vm.stack.stackPointer, &base, exe)
		default:
			panic("TODO")
		}
	}
}

func (vm *VmVirtualMachine) initializeValue(basicType BasicType) VmValue {
	var value VmValue
	switch basicType {
	case BooleanType:
		fallthrough
	case IntType:
		value = &VmIntValue{intValue: 0}
	case DoubleType:
		value = &VmDoubleValue{doubleValue: 0.0}
	case StringType:
		value = &VmObjectValue{objectValue: vm.createStringObject("")}
	default:
		panic("TODO")
	}
	return value
}

func (vm *VmVirtualMachine) initializeLocalVariables(f *VmFunction, from_sp int) {
	var i, sp_idx int

	for i = 0; i < len(f.LocalVariableList); i++ {
		vm.stack.stack[i].setPointer(false)
	}

	sp_idx = from_sp
	for i = 0; i < len(f.LocalVariableList); i++ {
		vm.stack.stack[sp_idx] = vm.initializeValue(f.LocalVariableList[i].TypeSpecifier.BasicType)

		if f.LocalVariableList[i].TypeSpecifier.BasicType == StringType {
			vm.stack.stack[i].setPointer(true)
		}
		sp_idx++
	}
}

// 修正转换code
func (vm *VmVirtualMachine) convertCode(exe *Executable, codeList []byte, f *VmFunction) {
	var dest_idx int

	for i, code := range codeList {
		switch code {
		case VM_PUSH_STACK_INT, VM_POP_STACK_INT,
			VM_PUSH_STACK_DOUBLE, VM_POP_STACK_DOUBLE,
			VM_PUSH_STACK_STRING, VM_POP_STACK_STRING:

			if f == nil {
				panic("f == nil")
			}

			// 增加返回值的位置
			src_idx := get2ByteInt(codeList[i+1:])
			if src_idx >= len(f.ParameterList) {
				dest_idx = src_idx + 1
			} else {
				dest_idx = src_idx
			}
			set2ByteInt(codeList[i+1:], dest_idx)

		case VM_PUSH_FUNCTION:

			idx_in_exe := get2ByteInt(codeList[i+1:])
			func_idx := vm.SearchFunction(exe.FunctionList[idx_in_exe].Name)
			set2ByteInt(codeList[i+1:], func_idx)
		}

		info := &OpcodeInfo[code]
		for _, p := range []byte(info.Parameter) {
			switch p {
			case 'b':
				i++
			case 's': /* FALLTHRU */
				fallthrough
			case 'p':
				i += 2
				break
			default:
				panic("TODO")
			}
		}
	}
}

// 查找函数
func (vm *VmVirtualMachine) SearchFunction(name string) int {

	for i, f := range vm.function {
		if f.getName() == name {
			return i
		}
	}
	panic("TODO")
}

//
// 函数相关
//
// 函数执行
func (vm *VmVirtualMachine) invokeNativeFunction(f *NativeFunction, sp_p *int) {

	stack := vm.stack.stack
	sp := *sp_p

	ret := f.proc(vm, f.argCount, stack[sp-f.argCount-1])

	stack[sp-f.argCount-1] = ret

	*sp_p = sp - f.argCount
}

// 执行原生函数
func (vm *VmVirtualMachine) invokeGFunction(caller_p *GFunction, callee *GFunction,
	code_p []byte, pc_p *int, sp_p *int, base_p *int,
	exe_p *Executable) {

	// callee 函数指针

	exe_p = callee.Executable
	callee_p := exe_p.FunctionList[callee.Index]

	callInfo := &CallInfo{}
	vm.stack.stack[*sp_p-1] = callInfo
	callInfo.caller = caller_p
	callInfo.caller_address = *pc_p
	callInfo.base = *base_p

	*base_p = *sp_p - len(callee_p.ParameterList) - 1
	caller_p = callee

	vm.initializeLocalVariables(callee_p, *sp_p)

	*sp_p += len(callee_p.LocalVariableList)
	*pc_p = 0

	code_p = exe_p.FunctionList[callee.Index].CodeList
}

func (vm *VmVirtualMachine) returnFunction(func_p *GFunction, code_p []byte, pc_p *int, sp_p *int, base_p *int, exe_p *Executable) {

	return_value := vm.stack.stack[(*sp_p)-1]

	callee_p := exe_p.FunctionList[(*func_p).Index]
	callInfo := &CallInfo{}
	vm.stack.stack[*sp_p-1-len(callee_p.LocalVariableList)-1] = callInfo

	if callInfo.caller != nil {
		exe_p = callInfo.caller.Executable
		caller_p := exe_p.FunctionList[callInfo.caller.Index]
		code_p = caller_p.CodeList
	} else {
		exe_p = vm.executable
		code_p = vm.executable.CodeList
	}
	func_p = callInfo.caller

	*pc_p = callInfo.caller_address + 1
	*base_p = callInfo.base

	*sp_p -= len(callee_p.LocalVariableList) + 1 + len(callee_p.ParameterList)

	vm.stack.stack[*sp_p-1] = return_value
}
