package vm

import (
	// "fmt"

	"github.com/lth-go/gogo/utils"
)

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

	vm.SetTopExe(exeList.Top())
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

// 设置入口为main函数
func (vm *VirtualMachine) SetMainEntrypoint() {
	idx := vm.SearchStatic("main", "main")
	if idx == -1 {
		panic("TODO")
	}

	b := make([]byte, 2)
	utils.Set2ByteInt(b, idx)
	vm.topLevel.CodeList = append(vm.topLevel.CodeList, b...)
	vm.topLevel.CodeList = append(vm.topLevel.CodeList, OP_CODE_INVOKE)
}

// 添加单个exe到vm
func (vm *VirtualMachine) AddExecutable(exe *Executable) {
	vm.executableList = append(vm.executableList, exe)

	vm.AddFunctions(exe)
	vm.AddStatic(exe)

	// 修正字节码
	// 方法调用修正
	// 函数下标修正
	for _, f := range exe.FunctionList {
		vm.ConvertOpCode(exe, f.CodeList, f)
	}
}

// 添加静态区
func (vm *VirtualMachine) AddStatic(exe *Executable) {
	// 变量初始化
	for _, value := range exe.VariableList {
		value.Init()
	}

	for _, value := range exe.VariableList {
		if exe.PackageName != value.PackageName {
			continue
		}

		if vm.static.Index(value.PackageName, value.Name) == -1 {
			vm.static.Append(NewStaticVariable(value.PackageName, value.Name, value.Value))
		}
	}
}

// 添加exe函数到虚拟机
func (vm *VirtualMachine) AddFunctions(exe *Executable) {
	// 检查函数是否重复定义
	for _, exeFunc := range exe.FunctionList {
		if !exeFunc.IsImplemented {
			continue
		}

		if vm.SearchStatic(exeFunc.PackageName, exeFunc.Name) != -1 {
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

func (vm *VirtualMachine) execute(gogoFunc *GoGoFunction, codeList []byte) {
	var base int

	stack := vm.stack
	exe := vm.currentExecutable
	static := vm.static

	for pc := vm.pc; pc < len(codeList); {
		switch codeList[pc] {
		case OP_CODE_PUSH_INT_1BYTE:
			stack.SetIntPlus(0, int(codeList[pc+1]))
			vm.stack.stackPointer++
			pc += 2
		case OP_CODE_PUSH_INT_2BYTE:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, index)
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_PUSH_INT:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, exe.ConstantPool.GetInt(index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_PUSH_FLOAT_0:
			stack.SetFloatPlus(0, 0.0)
			vm.stack.stackPointer++
			pc++
		case OP_CODE_PUSH_FLOAT_1:
			stack.SetFloatPlus(0, 1.0)
			vm.stack.stackPointer++
			pc++
		case OP_CODE_PUSH_FLOAT:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetFloatPlus(0, exe.ConstantPool.GetFloat(index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_PUSH_STRING:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetStringPlus(0, exe.ConstantPool.GetString(index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_PUSH_NIL:
			stack.SetPlus(0, NilObject)
			vm.stack.stackPointer++
			pc++
		case OP_CODE_PUSH_STACK:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetPlus(0, stack.Get(base+index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_POP_STACK:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.Set(base+index, stack.GetPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case OP_CODE_PUSH_STATIC:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetPlus(0, static.GetVariableObject(index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_POP_STATIC:
			index := utils.Get2ByteInt(codeList[pc+1:])
			static.SetVariable(index, stack.GetPlus(-1))
			vm.stack.stackPointer--
			pc += 3
		case OP_CODE_PUSH_ARRAY:
			array := stack.GetArrayPlus(-2)
			index := stack.GetIntPlus(-1)

			object := array.Get(index)

			stack.SetPlus(-2, object)
			vm.stack.stackPointer--
			pc++
		case OP_CODE_POP_ARRAY:
			value := stack.GetPlus(-3)
			array := stack.GetArrayPlus(-2)
			index := stack.GetIntPlus(-1)

			array.Set(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case OP_CODE_PUSH_MAP:
			map_ := stack.GetMapPlus(-2)
			index := stack.GetPlus(-1)

			object := map_.Get(index)

			stack.SetPlus(-2, object)
			vm.stack.stackPointer--
			pc++
		case OP_CODE_POP_MAP:
			value := stack.GetPlus(-3)
			map_ := stack.GetMapPlus(-2)
			index := stack.GetPlus(-1)

			map_.Set(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case OP_CODE_PUSH_STRUCT:
			struct_ := stack.GetStructPlus(-2)
			index := stack.GetIntPlus(-1)

			object := struct_.GetField(index)

			stack.SetPlus(-2, object)
			vm.stack.stackPointer--
			pc++
		case OP_CODE_POP_STRUCT:
			value := stack.GetPlus(-3)
			struct_ := stack.GetStructPlus(-2)
			index := stack.GetIntPlus(-1)

			struct_.SetField(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case OP_CODE_ADD_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)+stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_ADD_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)+stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_ADD_STRING:
			stack.SetStringPlus(-2, stack.GetStringPlus(-2)+stack.GetStringPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_SUB_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)-stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_SUB_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)-stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_MUL_INT:
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)*stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_MUL_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)*stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_DIV_INT:
			if stack.GetIntPlus(-1) == 0 {
				vmError(DIVISION_BY_ZERO_ERR)
			}
			stack.SetIntPlus(-2, stack.GetIntPlus(-2)/stack.GetIntPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_DIV_FLOAT:
			stack.SetFloatPlus(-2, stack.GetFloatPlus(-2)/stack.GetFloatPlus(-1))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_MINUS_INT:
			stack.SetIntPlus(-1, -stack.GetIntPlus(-1))
			pc++
		case OP_CODE_MINUS_FLOAT:
			stack.SetFloatPlus(-1, -stack.GetFloatPlus(-1))
			pc++
		case OP_CODE_EQ_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) == stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_EQ_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) == stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_EQ_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(!(stack.GetStringPlus(-2) == stack.GetStringPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_EQ_OBJECT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetPlus(-2) == stack.GetPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GT_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) > stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GT_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) > stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GT_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) > stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GE_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) >= stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GE_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) >= stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_GE_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) >= stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LT_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) < stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LT_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) < stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LT_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) < stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LE_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) <= stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LE_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) <= stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LE_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) <= stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_NE_INT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetIntPlus(-2) != stack.GetIntPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_NE_FLOAT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetFloatPlus(-2) != stack.GetFloatPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_NE_OBJECT:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetPlus(-2) != stack.GetPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_NE_STRING:
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) != stack.GetStringPlus(-1)))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LOGICAL_AND:
			stack.SetIntPlus(-2, utils.BoolToInt(utils.IntToBool(stack.GetIntPlus(-2)) && utils.IntToBool(stack.GetIntPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LOGICAL_OR:
			stack.SetIntPlus(-2, utils.BoolToInt(utils.IntToBool(stack.GetIntPlus(-2)) || utils.IntToBool(stack.GetIntPlus(-1))))
			vm.stack.stackPointer--
			pc++
		case OP_CODE_LOGICAL_NOT:
			stack.SetIntPlus(-1, utils.BoolToInt(!utils.IntToBool(stack.GetIntPlus(-1))))
			pc++
		case OP_CODE_POP:
			vm.stack.stackPointer--
			pc++
		case OP_CODE_DUPLICATE:
			stack.Set(vm.stack.stackPointer, stack.Get(vm.stack.stackPointer-1))
			vm.stack.stackPointer++
			pc++
		case OP_CODE_DUPLICATE_OFFSET:
			offset := utils.Get2ByteInt(codeList[pc+1:])
			stack.Set(vm.stack.stackPointer, stack.Get(vm.stack.stackPointer-1-offset))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_JUMP:
			index := utils.Get2ByteInt(codeList[pc+1:])
			pc = index
		case OP_CODE_JUMP_IF_TRUE:
			if utils.IntToBool(stack.GetIntPlus(-1)) {
				index := utils.Get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case OP_CODE_JUMP_IF_FALSE:
			if !utils.IntToBool(stack.GetIntPlus(-1)) {
				index := utils.Get2ByteInt(codeList[pc+1:])
				pc = index
			} else {
				pc += 3
			}
			vm.stack.stackPointer--
		case OP_CODE_PUSH_FUNCTION:
			value := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetIntPlus(0, value)
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_INVOKE:
			funcIdx := stack.GetIntPlus(-1)
			switch f := vm.static.Get(funcIdx).(type) {
			case *GoGoNativeFunction:
				vm.InvokeNativeFunction(f, &vm.stack.stackPointer)
				pc++
			case *GoGoFunction:
				vm.InvokeFunction(&gogoFunc, f, &codeList, &pc, &vm.stack.stackPointer, &base, &exe)
			default:
				panic("TODO")
			}
		case OP_CODE_RETURN:
			vm.ReturnFunction(&gogoFunc, &codeList, &pc, &base, &exe)
		case OP_CODE_NEW_ARRAY:
			size := utils.Get2ByteInt(codeList[pc+1:])
			array := vm.NewObjectArray(size)

			vm.stack.stackPointer -= size
			stack.SetPlus(0, array)
			vm.stack.stackPointer++
			pc += 3
		case OP_CDOE_NEW_MAP:
			size := utils.Get2ByteInt(codeList[pc+1:])
			objectMap := vm.NewObjectMap(size)

			vm.stack.stackPointer -= size * 2
			stack.SetPlus(0, objectMap)
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_NEW_INTERFACE:
			data := stack.GetPlus(-1)
			ifs := vm.NewObjectInterface(data)

			vm.stack.stackPointer -= 1
			stack.SetPlus(0, ifs)
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_NEW_STRUCT:
			size := utils.Get2ByteInt(codeList[pc+1:])
			struct_ := vm.NewObjectStruct(size)

			vm.stack.stackPointer -= size
			stack.SetPlus(0, struct_)
			vm.stack.stackPointer++
			pc += 3
		default:
			panic("TODO")
		}
	}
}

func (vm *VirtualMachine) InitFuncLocalVariables(f *Function, fromSp int) {
	var i, spIdx int

	spIdx = fromSp
	for i = 0; i < len(f.VariableList); i++ {
		vm.stack.Set(spIdx, GetObjectByType(f.VariableList[i].Type))
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
		case OP_CODE_PUSH_STACK, OP_CODE_POP_STACK:

			var parameterCount int

			if f == nil {
				panic("can't find line, need debug!!!")
			}

			parameterCount = f.GetParamCount()
			if f.IsMethod {
				parameterCount += 1 /* for this */
			}

			// 增加返回值的位置
			srcIdx := utils.Get2ByteInt(codeList[i+1:])
			if srcIdx >= parameterCount {
				destIdx = srcIdx + 1
			} else {
				destIdx = srcIdx
			}
			utils.Set2ByteInt(codeList[i+1:], destIdx)
		case OP_CODE_PUSH_STATIC, OP_CODE_POP_STATIC:

			idxInExe := utils.Get2ByteInt(codeList[i+1:])
			packageName := exe.PackageName
			if exe.VariableList[idxInExe].PackageName != "" {
				packageName = exe.VariableList[idxInExe].PackageName
			}
			funcIdx := vm.SearchStatic(packageName, exe.VariableList[idxInExe].Name)
			utils.Set2ByteInt(codeList[i+1:], funcIdx)

		case OP_CODE_PUSH_FUNCTION:
			idxInExe := utils.Get2ByteInt(codeList[i+1:])
			funcIdx := vm.SearchStatic(exe.FunctionList[idxInExe].PackageName, exe.FunctionList[idxInExe].Name)
			utils.Set2ByteInt(codeList[i+1:], funcIdx)
		}

		info := OpcodeInfo[code]
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
// 函数相关 执行原生函数
//
func (vm *VirtualMachine) InvokeNativeFunction(f *GoGoNativeFunction, spP *int) {
	// 5 -- sp
	// 4 funcName
	// 3 arg3
	// 2 arg2
	// 1 arg1
	// 0 base

	sp := *spP

	resultList := f.proc(vm, f.argCount, vm.stack.objectList[sp-f.argCount-1:])

	resultLen := len(resultList)

	for i, value := range resultList {
		vm.stack.Set(sp-f.argCount-1+resultLen-i-1, value)
	}

	*spP = sp - f.argCount - 1 + resultLen
}

// 函数执行
func (vm *VirtualMachine) InvokeFunction(
	caller **GoGoFunction,
	callee *GoGoFunction,
	codeP *[]byte,
	pcP *int,
	spP *int,
	baseP *int,
	exe **Executable,
) {
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
	vm.InitFuncLocalVariables(calleeP, *spP)

	// 设置栈位置
	*spP += len(calleeP.VariableList)
	*pcP = 0

	// 设置字节码为函数的字节码
	*codeP = calleeP.CodeList
}

// 保存返回值,并恢复栈
func (vm *VirtualMachine) ReturnFunction(
	funcP **GoGoFunction,
	codeP *[]byte,
	pcP *int,
	baseP *int,
	exeP **Executable,
) {
	calleeP := (*exeP).FunctionList[(*funcP).Index]
	resultCount := calleeP.GetResultCount()

	bakObjList := make([]Object, resultCount)

	for i := 0; i < resultCount; i++ {
		bakObjList[i] = vm.stack.Get(vm.stack.stackPointer - resultCount + i)
	}
	vm.stack.stackPointer -= resultCount

	// 恢复调用栈
	DoReturn(vm, funcP, codeP, pcP, baseP, exeP)

	for i := 0; i < resultCount; i++ {
		vm.stack.Set(vm.stack.stackPointer, bakObjList[i])
		vm.stack.stackPointer++
	}
}

// 恢复到父调用栈
func DoReturn(vm *VirtualMachine, funcP **GoGoFunction, codeP *[]byte, pcP *int, baseP *int, exeP **Executable) {

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
}

func (vm *VirtualMachine) NewObjectArray(size int) Object {
	obj := NewObjectArray(size)

	vm.AddObject(obj)

	for i := 0; i < size; i++ {
		obj.Set(i, vm.stack.GetPlus(-size+i))
	}

	return obj
}

func (vm *VirtualMachine) NewObjectMap(size int) Object {
	obj := NewObjectMap()

	vm.AddObject(obj)

	for i := 0; i < size; i++ {
		keyIndex := -size + i
		valueIndex := (-size + i) - size

		key := vm.stack.GetPlus(keyIndex)
		value := vm.stack.GetPlus(valueIndex)
		obj.Set(key, value)
	}

	return obj
}

func (vm *VirtualMachine) NewObjectInterface(data Object) Object {
	obj := NewObjectInterface(data)

	vm.AddObject(obj)

	return obj
}

func (vm *VirtualMachine) NewObjectStruct(size int) Object {
	obj := NewObjectStruct(size)

	vm.AddObject(obj)

	// TODO: 倒序入栈, 正序出栈
	for i := 0; i < size; i++ {
		obj.FieldList[i] = vm.stack.GetPlus(-size + i)
	}

	return obj
}

func GetObjectByType(typ *Type) Object {
	var value Object

	if typ.IsReferenceType() {
		value = NilObject
		return value
	}

	switch typ.BasicType {
	case BasicTypeVoid, BasicTypeBool, BasicTypeInt:
		value = NewObjectInt(0)
	case BasicTypeFloat:
		value = NewObjectFloat(0.0)
	case BasicTypeString:
		value = NewObjectString("")
	case BasicTypeStruct:
		structValue := NewObjectStruct(len(typ.StructType.FieldTypeList))
		for i, fieldType := range typ.StructType.FieldTypeList {
			structValue.SetField(i, GetObjectByType(fieldType))
		}

		value = structValue
	case BasicTypeNil:
		fallthrough
	default:
		panic("TODO")
	}

	return value
}
