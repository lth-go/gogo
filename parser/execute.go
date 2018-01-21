package parser

import (
	"encoding/binary"
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
	function []TODOFunction
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
func (vm *VmVirtualMachine) addExecutable(exe *Executable) {

	vm.executable = exe

	vm.addFunctions(exe)

	vm.convertCode(exe, exe.codeList, nil)

	for _, f := range exe.functionList {
		vm.convertCode(exe, f.codeList, f)
	}

	vm.addStaticVariables(exe)
}

//

func (vm *VmVirtualMachine) STD(sp int) {
	// TODO 应该返回指针
	return vm.stack.stack[vm.stack.stackPointer+sp].getNumberValue()
}
func (vm *VmVirtualMachine) STO(sp int) {
	return vm.stack.stack[vm.stack.stackPointer+sp].getObjectValue()
}

//

func (vm *VmVirtualMachine) STD_I(sp int) {
	return vm.stack.stack[sp].getNumberValue()
}

func (vm *VmVirtualMachine) STO_I(sp int) {
	return vm.stack.stack[sp].getObjectValue()
}

//
func (vm *VmVirtualMachine) STD_WRITE(sp int, r float64) {
	v := vm.stack.stack[vm.stack.stackPointer+sp]

	v.setNumberValue(r)
	v.setPointer(false)
}

func (vm *VmVirtualMachine) STO_WRITE(sp int, r VmObject) {
	v := vm.stack.stack[vm.stack.stackPointer+sp]

	v.setObjectValue(r)
	v.setPointer(true)
}

//

func (vm *VmVirtualMachine) STD_WRITE_I(sp int, r float64) {
	v := vm.stack.stack[sp]

	v.setNumberValue(r)
	v.setPointer(false)
}

func (vm *VmVirtualMachine) STO_WRITE_I(sp int, r float64) {
	v := vm.stack.stack[sp]

	v.setObjectValue(r)
	v.setPointer(true)
}

func (vm *VmVirtualMachine) alloc_object_string() {

	//check_gc(vm)
	ret := VmObjectString{}

	ret.marked = false

	vm.heap.objectList = append(vm.heap.objectList, ret)

	return ret
}

func (vm *VmVirtualMachine) literal_to_vm_string_i(value string) VmObject {
	// TODO
	ret := vm.alloc_object_string()

	ret.stringValue = value
	ret.isLiteral = true

	return ret
}

func GET_2BYTE_INT(b []byte, pc int) int {
	return int(binary.BigEndian.Uint16(b))
}
func SET_2BYTE_INT(codeList []byte, pc int, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

func (vm *VmVirtualMachine) Execute() {
	vm.pc = 0

	execute(nil, vm.executable.codeList)
}

func (vm *VmVirtualMachine) execute(funcList []Function, codeList []byte) {
	var base int

	stack := vm.stack.stack
	exe := vm.executable

	for pc := vm.pc; pc < len(codeList); {

		switch codeList[pc] {
		case VM_PUSH_NUMBER_0:
			vm.STD_WRITE(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_NUMBER_1:
			vm.STD_WRITE(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case VM_PUSH_NUMBER:
			vm.STD_WRITE(0, exe.constantPool[GET_2BYTE_INT(codeList[pc+1:pc+3])].getNumber())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STRING:
			vm.STO_WRITE(0, vm.literal_to_vm_string_i(exe.constantPool[GET_2BYTE_INT(codeList[pc+1:pc+3])].getString()))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_NUMBER:
			vm.STD_WRITE(0, vm.STD_I(base+GET_2BYTE_INT(codeList[pc+1:pc+3])))
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STACK_STRING:
			vm.STO_WRITE(0, vm.STO_I(base+GET_2BYTE_INT(codeList[pc+1:pc+3])))
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STACK_NUMBER:
			vm.STD_WRITE_I(base+GET_2BYTE_INT(codeList[pc+1:pc+3]), vm.STD(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STACK_STRING:
			vm.STO_WRITE_I(base+GET_2BYTE_INT(codeList[pc+1:pc+3]), vm.STO(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_STATIC_NUMBER:
			vm.STD_WRITE(0, vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:pc+3])].getNumberValue())
			vm.stack.stackPointer++
			pc += 3
		case VM_PUSH_STATIC_STRING:
			vm.STO_WRITE(0, vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:pc+3])].getObjectValue())
			vm.stack.stackPointer++
			pc += 3
		case VM_POP_STATIC_NUMBER:
			vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:pc+3])].setNumberValue(vm.STD(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_POP_STATIC_STRING:
			vm.static.variableList[GET_2BYTE_INT(codeList[pc+1:pc+3])].setObjectValue(vm.STO(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_ADD_NUMBER:
			vm.STD(-2) = vm.STD(-2) + vm.STD(-1)
			vm.stack.stackPointer--
			pc++
		case VM_ADD_STRING:
			vm.STO(-2) = chain_string(vm, vm.STO(-2), vm.STO(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_NUMBER:
			vm.STD(-2) = vm.STD(-2) - vm.STD(-1)
			vm.stack.stackPointer--
			pc++
		case VM_MUL_NUMBER:
			vm.STD(-2) = vm.STD(-2) * vm.STD(-1)
			vm.stack.stackPointer--
			pc++
		case VM_DIV_NUMBER:
			vm.STD(-2) = vm.STD(-2) / vm.STD(-1)
			vm.stack.stackPointer--
			pc++
		case VM_MINUS_NUMBER:
			vm.STD(-1) = -vm.STD(-1)
			pc++
		case VM_EQ_NUMBER:
			// 用boolValue替换int
			STI(vm, -2) = (vm.STD(-2) == vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			STI_WRITE(vm, -2, !vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string))
			vm.stack.stackPointer--
			pc++
		case VM_GT_NUMBER:
			STI(vm, -2) = (vm.STD(-2) > vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_GT_STRING:
			STI_WRITE(vm, -2, vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string) > 0)
			vm.stack.stackPointer--
			pc++
		case VM_GE_NUMBER:
			STI(vm, -2) = (vm.STD(-2) >= vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_GE_STRING:
			STI_WRITE(vm, -2,
				vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string) >= 0)
			vm.stack.stackPointer--
			pc++
		case VM_LT_NUMBER:
			STI(vm, -2) = (vm.STD(-2) < vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_LT_STRING:
			STI_WRITE(vm, -2,
				vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string) < 0)
			vm.stack.stackPointer--
			pc++
		case VM_LE_NUMBER:
			STI(vm, -2) = (vm.STD(-2) <= vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_LE_STRING:
			STI_WRITE(vm, -2,
				vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string) <= 0)
			vm.stack.stackPointer--
			pc++
		case VM_NE_NUMBER:
			STI(vm, -2) = (vm.STD(-2) != vm.STD(-1))
			vm.stack.stackPointer--
			pc++
		case VM_NE_STRING:
			STI_WRITE(vm, -2,
				vm_wcscmp(vm.STO(-2).u.string.string, vm.STO(-1).u.string.string) != 0)
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_AND:
			STI(vm, -2) = (STI(vm, -2) && STI(vm, -1))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_OR:
			STI(vm, -2) = (STI(vm, -2) || STI(vm, -1))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_NOT:
			STI(vm, -1) = !STI(vm, -1)
			pc++
		case VM_POP:
			vm.stack.stackPointer--
			pc++
		case VM_DUPLICATE:
			stack[vm.stack.stackPointer] = stack[vm.stack.stackPointer-1]
			vm.stack.stackPointer++
			pc++
		case VM_JUMP:
			pc = GET_2BYTE_INT(codeList[pc+1 : pc+3])
		case VM_JUMP_IF_TRUE:
			if STI(vm, -1) {
				pc = GET_2BYTE_INT(codeList[pc+1 : pc+3])
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_JUMP_IF_FALSE:
			if !STI(vm, -1) {
				pc = GET_2BYTE_INT(codeList[pc+1 : pc+3])
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case VM_PUSH_FUNCTION:
			STI_WRITE(vm, 0, GET_2BYTE_INT(codeList[pc+1:pc+3]))
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			func_idx := STI(vm, -1)
			if vm.function[func_idx].kind == NATIVE_FUNCTION {
				invoke_native_function(vm, &vm.function[func_idx], &vm.stack.stackPointer)
				pc++
			} else {
				invoke_g_function(vm, funcList, &vm.function[func_idx], &code, &code_size, &pc, &vm.stack.stackPointer, &base, &exe)
			}
		case VM_RETURN:
			return_function(vm, funcList, &code, &code_size, &pc, &vm.stack.stackPointer, &base, &exe)
		default:
			panic("TODO")
		}
	}
}

func (vm *VmVirtualMachine) addStaticVariables(exe *Executable) {}

func (vm *VmVirtualMachine) addFunctions(exe *Executable) {}

func (vm *VmVirtualMachine) convertCode(exe *Executable, codeList []byte, f *VmFunction) {

	for i, code := range codeList {
		switch code {
		case VM_PUSH_STACK_NUMBER, VM_PUSH_STACK_STRING, VM_POP_STACK_NUMBER, VM_POP_STACK_STRING:

			if f == nil {
				panic("f == nil")
			}

			src_idx = GET_2BYTE_INT(codeList[pc+1 : pc+3])
			if src_idx >= len(f.parameter) {
				dest_idx = src_idx + CALL_INFO_ALIGN_SIZE
			} else {
				dest_idx = src_idx
			}
			SET_2BYTE_INT(codeList[i+1:i+3], dest_idx)

		case VM_PUSH_FUNCTION:

			idx_in_exe := GET_2BYTE_INT(codeList[pc+1 : pc+3])
			func_idx := search_function(vm, exe.function[idx_in_exe].name)
			SET_2BYTE_INT(codeList[i+1:i+3], func_idx)
		}

		info = &OpcodeInfo[code]
		for _, p := range []byte(info.parameter) {
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

func newVirtualMachine() *VmVirtualMachine {
	vm := VmVirtualMachine{}

	vm.stack.stack = []VmValue{}
	vm.stack.stackPointer = 0

	vm.heap.objectList = []VmObject{}

	//vm.heap.currentThreshold = HEAP_THRESHOLD_SIZE;
	vm.function = []TODOFunction{}

	vm.executable = nil

	vm.AddNativeFunctions()

	return vm
}

//
// Stack
//
type Stack struct {
	int   stackPointer
	stack []VmValue
}

//
// VmValue
//

type VmValue interface {
	getNumberValue() float64
	setNumberValue(float64)

	getObjectValue() VmObject
	setObjectValue(VmObject)
	isPointer() bool
}

// VmValueImpl

type VmValueImpl struct {
	// 是否是指针
	pointerFlags bool
}

func (v *VmValueImpl) getNumberValue() float64 {
	panic("error")
}

func (v *VmValueImpl) setNumberValue(value float64) {
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

// VmNumberValue

type VmNumberValue struct {
	VmValueImpl
	numberValue float64
}

func (v *VmNumberValue) getNumberValue() float64 {
	return v.numberValue
}

func (v *VmNumberValue) setNumberValue(value float64) {
	v.numberValue = value
}

// VmObjectValue

type VmObjectValue struct {
	VmValueImpl

	objectValue VmObject
}

func (v *VmObjectValue) getObjectValue() VmObject {
	return v.objectValue
}

func (v *VmObjectValue) setObjectValue(VmObject) {
	v.objectValue = VmObject
}

//
// VmObject
//
type VmObject interface{}

type VmObjectImpl struct {
	// gc用
	marked bool
}

type VmObjectString struct {
	VmObjectImpl

	stringValue string
	isLiteral   bool
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
type TODOFunction interface{}

// 原生函数
type NativeFunction struct {
	name string

	proc     *VmNativeFunctionProc
	argCount int
}

type VmNativeFunctionProc func(vm *VmVirtualMachine, argCount int, args *VmValue) VmValue

type GFunction struct {
	name string

	executable *Executable
	index      int
}

//
// 字节码解释器
//

type Executable struct {
	// 常量池
	constantPool []Constant

	// 全局变量
	globalVariableList []*Variable

	// 函数列表
	functionList []Function

	// 顶层结构代码
	codeList []byte

	// 行号对应表
	// 保存字节码和与之对应的源代码的行号
	lineNumberList []*LineNumber
}

// ==============================
// 常量池
// ==============================

type Constant interface {
	getNumber() float64
	getString() string
}

type ConstantImpl struct{}

func (c *ConstantImpl) getNumber() float64 {
	panic("error")
}

func (c *ConstantImpl) getString() float64 {
	panic("error")
}

type ConstantNumber struct {
	numberValue float64
}

func (c *ConstantNumber) getNumber() float64 {
	return c.numberValue
}

type ConstantString struct {
	stringValue string
}

func (c *ConstantString) getString() float64 {
	return c.stringValue
}

// ==============================
// 全局变量
// ==============================

type Variable struct {
	name          string
	typeSpecifier *TypeSpecifier
}

// ==============================
// 函数
// ==============================

type Function struct {
	// 类型
	typeSpecifier *TypeSpecifier
	// 函数名
	name string
	// 形参列表
	parameterList []*LocalVariable
	// 是否原生函数
	isImplemented bool
	// 局部变量列表
	localVariableList []*LocalVariable
	// 字节码类表
	codeList []byte

	// 行号对应表
	lineNumberList []*LineNumber
}

type LocalVariable struct {
	name          string
	typeSpecifier *TypeSpecifier
}

// ==============================
// 行号对应表
// ==============================

type LineNumber struct {
	// 源代码行号
	lineNumber int

	// 字节码开始的位置
	startPc int

	// 接下来有多少字节码对应相同的行号
	pcCount int
}

func newExecutable() *Executable {
	exe := &Executable{}
	return exe
}

// ==============================
// 字节码
// ==============================

// 字节码
const (
	PUSH_NUMBER_0 byte = iota
	PUSH_NUMBER_1
	PUSH_NUMBER
	PUSH_STRING
	/**********/
	PUSH_STACK_NUMBER
	PUSH_STACK_STRING
	POP_STACK_NUMBER
	POP_STACK_STRING
	/**********/
	PUSH_STATIC_NUMBER
	PUSH_STATIC_STRING
	POP_STATIC_NUMBER
	POP_STATIC_STRING
	/**********/
	ADD_NUMBER
	ADD_STRING
	SUB_NUMBER
	MUL_NUMBER
	DIV_NUMBER
	MOD_NUMBER
	MINUS_NUMBER
	INCREMENT
	DECREMENT
	CAST_BOOLEAN_TO_STRING
	CAST_NUMBER_TO_STRING
	EQ_NUMBER
	EQ_STRING
	GT_NUMBER
	GT_STRING
	GE_NUMBER
	GE_STRING
	LT_NUMBER
	LT_STRING
	LE_NUMBER
	LE_STRING
	NE_NUMBER
	NE_STRING
	LOGICAL_AND
	LOGICAL_OR
	LOGICAL_NOT
	POP
	DUPLICATE
	JUMP
	JUMP_IF_TRUE
	JUMP_IF_FALSE
	/**********/
	PUSH_FUNCTION
	INVOKE
	RETURN
)

type OpcodeInfo struct {
	// 注记符
	mnemonic string

	// 参数类型，
	// `b` 一个字节整数
	// `s` 两个字节整数
	// `p` 常量池索引值
	parameter       string
	stack_increment int
}
