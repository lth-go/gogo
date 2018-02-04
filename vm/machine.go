package vm

import (
	"encoding/binary"
	"fmt"
)

//
// Stack
//
type Stack struct {
	stackPointer int
	stack        []VmValue
}

//
// VmValue
//

type VmValue interface {
	getIntValue() int
	setIntValue(int)

	getDoubleValue() float64
	setDoubleValue(float64)

	getObjectValue() VmObject
	setObjectValue(VmObject)

	isPointer() bool
	setPointer(bool)
}

// VmValueImpl

type VmValueImpl struct {
	// 是否是指针
	pointerFlags bool
}

func (v *VmValueImpl) getIntValue() int {
	panic("error")
}

func (v *VmValueImpl) setIntValue(value int) {
	panic("error")
}

func (v *VmValueImpl) getDoubleValue() float64 {
	panic("error")
}

func (v *VmValueImpl) setDoubleValue(value float64) {
	panic("error")
}

func (v *VmValueImpl) getObjectValue() VmObject {
	panic("error")
}

func (v *VmValueImpl) setObjectValue(value VmObject) {
	panic("error")
}

func (v *VmValueImpl) isPointer() bool {
	return v.pointerFlags
}

func (v *VmValueImpl) setPointer(b bool) {
	v.pointerFlags = b
}

// CallInfo
// 函数返回体
type CallInfo struct {
	VmValueImpl

	caller         *GFunction
	caller_address int
	base           int
}

// VmIntValue
type VmIntValue struct {
	VmValueImpl
	intValue int
}

func (v *VmIntValue) getIntValue() int {
	return v.intValue
}

func (v *VmIntValue) setIntValue(value int) {
	v.intValue = value
}

// VmDoubleValue

type VmDoubleValue struct {
	VmValueImpl
	doubleValue float64
}

func (v *VmDoubleValue) getDoubleValue() float64 {
	return v.doubleValue
}

func (v *VmDoubleValue) setDoubleValue(value float64) {
	v.doubleValue = value
}

// VmObjectValue

type VmObjectValue struct {
	VmValueImpl

	objectValue VmObject
}

func (v *VmObjectValue) getObjectValue() VmObject {
	return v.objectValue
}

func (v *VmObjectValue) setObjectValue(value VmObject) {
	v.objectValue = value
}

//
// VmObject
//
type VmObject interface {
	isMarked() bool
	setMark(bool)

	getString() string
	setString(string)
}

type VmObjectImpl struct {
	// gc用
	marked bool
}

func (obj *VmObjectImpl) isMarked() bool {
	return obj.marked
}

func (obj *VmObjectImpl) setMark(m bool) {
	obj.marked = m
}

func (obj *VmObjectImpl) getString() string {
	panic("TODO")
}

func (obj *VmObjectImpl) setString(v string) {
	panic("TODO")
}

type VmObjectString struct {
	VmObjectImpl

	stringValue string
	isLiteral   bool
}

func (obj *VmObjectString) getString() string {
	return obj.stringValue
}

func (obj *VmObjectString) setString(v string) {
	obj.stringValue = v
}

//
// Heap
//
type Heap struct {
	// TODO:阈值
	currentThreshold int
	objectList       []VmObject
}

//
// Static
//
type Static struct {
	variableList []VmValue
}

//
// Function
//
type Function interface {
	getName() string
}

// 原生函数
type NativeFunction struct {
	Name string

	proc     VmNativeFunctionProc
	argCount int
}

func (f *NativeFunction) getName() string { return f.Name }

type VmNativeFunctionProc func(vm *VmVirtualMachine, argCount int, args VmValue) VmValue

type GFunction struct {
	Name string

	Executable *Executable
	Index      int
}

func (f *GFunction) getName() string { return f.Name }

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

func (vm *VmVirtualMachine) AddNativeFunctions() {
	// TODO
	//vm.addNativeFunction(vm, "print", nv_print_proc, 1)
}

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

	vm.static.variableList = []VmValue{}

	for _, exeValue := range exe.GlobalVariableList {
		newVmValue := vm.initializeValue(exeValue.typeSpecifier.BasicType)
		vm.static.variableList = append(vm.static.variableList, newVmValue)
	}
}

func (vm *VmVirtualMachine) addFunctions(exe *Executable) {
	var exeFunc *VmFunction
	var vmFunc, newVmFunc Function

	for _, exeFunc = range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}
		for _, vmFunc = range vm.function {
			if vmFunc.getName() == exeFunc.Name {
				panic("TODO")
			}
		}
	}

	for srcIdex, exeFunc := range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}

		newVmFunc = &GFunction{Name: exeFunc.Name, Executable: exe, Index: srcIdex}
		vm.function = append(vm.function, newVmFunc)
	}
}

//
// 栈操作
//
func (vm *VmVirtualMachine) STI_GET(sp int) int {
	return vm.stack.stack[vm.stack.stackPointer+sp].getIntValue()
}
func (vm *VmVirtualMachine) STI_SET(sp int, v int) {
	vm.stack.stack[vm.stack.stackPointer+sp].setIntValue(v)
}

func (vm *VmVirtualMachine) STD_GET(sp int) float64 {
	return vm.stack.stack[vm.stack.stackPointer+sp].getDoubleValue()
}
func (vm *VmVirtualMachine) STD_SET(sp int, v float64) {
	vm.stack.stack[vm.stack.stackPointer+sp].setDoubleValue(v)
}
func (vm *VmVirtualMachine) STO_GET(sp int) VmObject {
	return vm.stack.stack[vm.stack.stackPointer+sp].getObjectValue()
}
func (vm *VmVirtualMachine) STO_SET(sp int, v VmObject) {
	vm.stack.stack[vm.stack.stackPointer+sp].setObjectValue(v)
}

//

func (vm *VmVirtualMachine) STI_I(sp int) int {
	return vm.stack.stack[sp].getIntValue()
}

func (vm *VmVirtualMachine) STD_I(sp int) float64 {
	return vm.stack.stack[sp].getDoubleValue()
}

func (vm *VmVirtualMachine) STO_I(sp int) VmObject {
	return vm.stack.stack[sp].getObjectValue()
}

//
func (vm *VmVirtualMachine) STI_WRITE(sp int, r int) {
	v := vm.stack.stack[vm.stack.stackPointer+sp]
	v.setIntValue(r)
	v.setPointer(false)
}
func (vm *VmVirtualMachine) STD_WRITE(sp int, r float64) {
	v := vm.stack.stack[vm.stack.stackPointer+sp]
	v.setDoubleValue(r)
	v.setPointer(false)
}

func (vm *VmVirtualMachine) STO_WRITE(sp int, r VmObject) {
	v := vm.stack.stack[vm.stack.stackPointer+sp]
	v.setObjectValue(r)
	v.setPointer(true)
}

//
func (vm *VmVirtualMachine) STI_WRITE_I(sp int, r int) {
	v := vm.stack.stack[sp]
	v.setIntValue(r)
	v.setPointer(false)
}

func (vm *VmVirtualMachine) STD_WRITE_I(sp int, r float64) {
	v := vm.stack.stack[sp]
	v.setDoubleValue(r)
	v.setPointer(false)
}

func (vm *VmVirtualMachine) STO_WRITE_I(sp int, r VmObject) {
	v := vm.stack.stack[sp]
	v.setObjectValue(r)
	v.setPointer(true)
}

//

func GET_2BYTE_INT(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}
func SET_2BYTE_INT(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

func (vm *VmVirtualMachine) Execute() {
	vm.pc = 0

	vm.execute(nil, vm.executable.CodeList)
}

func (vm *VmVirtualMachine) execute(gFunc *GFunction, codeList []byte) {
	var base int

	stack := vm.stack.stack
	exe := vm.executable

	for pc := vm.pc; pc < len(codeList); {

		switch codeList[pc] {
		case VM_PUSH_INT_1BYTE:
			vm.STI_WRITE(0, int(codeList[pc+1]))
			vm.stack.stackPointer++
			pc += 2
		case VM_PUSH_INT_2BYTE:
			vm.STI_WRITE(0, GET_2BYTE_INT(codeList[pc+1:]))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_INT:
			vm.STI_WRITE(0, exe.ConstantPool[GET_2BYTE_INT(codeList[pc+1:])].getInt())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_DOUBLE_0:
			vm.STD_WRITE(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE_1:
			vm.STD_WRITE(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_DOUBLE:
			vm.STD_WRITE(0, exe.ConstantPool[GET_2BYTE_INT(codeList[pc+1:])].getDouble())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STRING:
			vm.STO_WRITE(0, vm.literal_to_vm_string_i(exe.ConstantPool[GET_2BYTE_INT(codeList[pc+1:])].getString()))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_INT:
			vm.STI_WRITE(0, vm.STI_I(base+GET_2BYTE_INT(codeList[pc+1:])))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_DOUBLE:
			vm.STD_WRITE(0, vm.STD_I(base+GET_2BYTE_INT(codeList[pc+1:])))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_STRING:
			vm.STO_WRITE(0, vm.STO_I(base+GET_2BYTE_INT(codeList[pc+1:])))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STACK_INT:
			vm.STI_WRITE_I(base+GET_2BYTE_INT(codeList[pc+1:]), vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_DOUBLE:
			vm.STD_WRITE_I(base+GET_2BYTE_INT(codeList[pc+1:]), vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_STRING:
			vm.STO_WRITE_I(base+GET_2BYTE_INT(codeList[pc+1:]), vm.STO_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_STATIC_INT:
			vm.STI_WRITE(0, vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].getIntValue())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_DOUBLE:
			vm.STD_WRITE(0, vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].getDoubleValue())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_STRING:
			vm.STO_WRITE(0, vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].getObjectValue())
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STATIC_INT:
			vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].setIntValue(vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_DOUBLE:
			vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].setDoubleValue(vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_STRING:
			vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:])].setObjectValue(vm.STO_GET(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_ADD_INT:
			vm.STI_SET(-2, vm.STI_GET(-2)+vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_DOUBLE:
			vm.STD_SET(-2, vm.STD_GET(-2)+vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_STRING:
			vm.STO_SET(-2, vm.chainString(vm.STO_GET(-2), vm.STO_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_INT:
			vm.STI_SET(-2, vm.STI_GET(-2)-vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_DOUBLE:
			vm.STD_SET(-2, vm.STD_GET(-2)-vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_INT:
			vm.STI_SET(-2, vm.STI_GET(-2)*vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_DOUBLE:
			vm.STD_SET(-2, vm.STD_GET(-2)*vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_INT:
			vm.STI_SET(-2, vm.STI_GET(-2)/vm.STI_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_DOUBLE:
			vm.STD_SET(-2, vm.STD_GET(-2)/vm.STD_GET(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MINUS_INT:
			vm.STI_SET(-1, -vm.STI_GET(-1))
			pc++
		case VM_MINUS_DOUBLE:
			vm.STD_SET(-1, -vm.STD_GET(-1))
			pc++
		case VM_CAST_INT_TO_DOUBLE:
			vm.STD_SET(-1, float64(vm.STI_GET(-1)))
			pc++
		case VM_CAST_DOUBLE_TO_INT:
			vm.STI_SET(-1, int(vm.STD_GET(-1)))
			pc++
		case VM_CAST_BOOLEAN_TO_STRING:
			if vm.STI_GET(-1) != 0 {
				vm.STO_WRITE(-1, vm.literal_to_vm_string_i("true"))
			} else {
				vm.STO_WRITE(-1, vm.literal_to_vm_string_i("false"))
			}
			pc++
		case VM_CAST_INT_TO_STRING:

			buf := fmt.Sprintf("%d", vm.STI_GET(-1))
			vm.STO_WRITE(-1, vm.create_vm_string_i(buf))
			pc++
		case VM_CAST_DOUBLE_TO_STRING:
			buf := fmt.Sprintf("%f", vm.STD_GET(-1))
			vm.STO_WRITE(-1, vm.create_vm_string_i(buf))
			pc++
		case VM_EQ_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) == vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) == vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			vm.STI_WRITE(-2, boolToInt(!(vm.STO_GET(-2).getString() == vm.STO_GET(-1).getString())))
			vm.stack.stackPointer--
			pc++
		case VM_GT_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) > vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) > vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_STRING:
			vm.STI_WRITE(-2, boolToInt(vm.STO_GET(-2).getString() > vm.STO_GET(-1).getString()))
			vm.stack.stackPointer--
			pc++
		case VM_GE_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) >= vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) >= vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_STRING:
			vm.STI_WRITE(-2, boolToInt(vm.STO_GET(-2).getString() >= vm.STO_GET(-1).getString()))
			vm.stack.stackPointer--
			pc++
		case VM_LT_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) < vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) < vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_STRING:
			vm.STI_WRITE(-2, boolToInt(vm.STO_GET(-2).getString() < vm.STO_GET(-1).getString()))
			vm.stack.stackPointer--
			pc++
		case VM_LE_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) <= vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) <= vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_STRING:
			vm.STI_WRITE(-2, boolToInt(vm.STO_GET(-2).getString() <= vm.STO_GET(-1).getString()))
			vm.stack.stackPointer--
			pc++
		case VM_NE_INT:
			vm.STI_SET(-2, boolToInt(vm.STI_GET(-2) != vm.STI_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_DOUBLE:
			vm.STI_SET(-2, boolToInt(vm.STD_GET(-2) != vm.STD_GET(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_STRING:
			vm.STI_WRITE(-2, boolToInt(vm.STO_GET(-2).getString() != vm.STO_GET(-1).getString()))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_AND:
			vm.STI_SET(-2, boolToInt(intToBool(vm.STI_GET(-2)) && intToBool(vm.STI_GET(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_OR:
			vm.STI_SET(-2, boolToInt(intToBool(vm.STI_GET(-2)) || intToBool(vm.STI_GET(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_NOT:
			vm.STI_SET(-1, boolToInt(!intToBool(vm.STI_GET(-1))))
			pc++
		case VM_POP:
			vm.stack.stackPointer--
			pc++
		case VM_DUPLICATE:
			stack[vm.stack.stackPointer] = stack[vm.stack.stackPointer-1]
			vm.stack.stackPointer++
			pc++
		case VM_JUMP:
			pc = GET_2BYTE_INT(codeList[pc+1:])
		case VM_JUMP_IF_TRUE:
			if intToBool(vm.STI_GET(-1)) {
				pc = GET_2BYTE_INT(codeList[pc+1:])
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_JUMP_IF_FALSE:
			if !intToBool(vm.STI_GET(-1)) {
				pc = GET_2BYTE_INT(codeList[pc+1:])
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_PUSH_FUNCTION:
			vm.STI_WRITE(0, GET_2BYTE_INT(codeList[pc+1:]))
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			func_idx := vm.STI_GET(-1)
			switch f := vm.function[func_idx].(type) {
			case *NativeFunction:
				invoke_native_function(vm, f, &vm.stack.stackPointer)
				pc++
			case *GFunction:
				invoke_g_function(vm, gFunc, f, codeList, &pc, &vm.stack.stackPointer, &base, exe)
			default:
				panic("TODO")
			}
		case VM_RETURN:
			return_function(vm, gFunc, codeList, &pc, &vm.stack.stackPointer, &base, exe)
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
		value = &VmObjectValue{objectValue: vm.literal_to_vm_string_i("")}
	default:
		panic("TODO")
	}
	return value
}

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
			src_idx := GET_2BYTE_INT(codeList[i+1:])
			if src_idx >= len(f.ParameterList) {
				dest_idx = src_idx + 1
			} else {
				dest_idx = src_idx
			}
			SET_2BYTE_INT(codeList[i+1:], dest_idx)

		case VM_PUSH_FUNCTION:

			idx_in_exe := GET_2BYTE_INT(codeList[i+1:])
			func_idx := vm.SearchFunction(exe.FunctionList[idx_in_exe].Name)
			SET_2BYTE_INT(codeList[i+1:], func_idx)
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

func (vm *VmVirtualMachine) chainString(str1 VmObject, str2 VmObject) VmObject {
	str := str1.getString() + str2.getString()
	ret := vm.create_vm_string_i(str)
	return ret
}

func (vm *VmVirtualMachine) SearchFunction(name string) int {

	for i, f := range vm.function {
		if f.getName() == name {
			return i
		}
	}
	panic("TODO")
}

func NewVirtualMachine() *VmVirtualMachine {
	vm := &VmVirtualMachine{}

	vm.stack.stack = []VmValue{}
	vm.stack.stackPointer = 0

	vm.heap.objectList = []VmObject{}

	//vm.heap.currentThreshold = HEAP_THRESHOLD_SIZE;
	vm.function = []Function{}

	vm.executable = nil

	vm.AddNativeFunctions()

	return vm
}

func invoke_native_function(vm *VmVirtualMachine, f *NativeFunction, sp_p *int) {

	stack := vm.stack.stack
	sp := *sp_p

	ret := f.proc(vm, f.argCount, stack[sp-f.argCount-1])

	stack[sp-f.argCount-1] = ret

	*sp_p = sp - f.argCount
}

func invoke_g_function(vm *VmVirtualMachine, caller_p *GFunction, callee *GFunction,
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

	vm.initialize_local_variables(callee_p, *sp_p)

	*sp_p += len(callee_p.LocalVariableList)
	*pc_p = 0

	code_p = exe_p.FunctionList[callee.Index].CodeList
}

func (vm *VmVirtualMachine) initialize_local_variables(f *VmFunction, from_sp int) {
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

func return_function(vm *VmVirtualMachine, func_p *GFunction, code_p []byte, pc_p *int, sp_p *int, base_p *int, exe_p *Executable) {

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
