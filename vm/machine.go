package vm

import (
	"github.com/lth-go/gogo/utils"
)

//
// 虚拟机
//
type VirtualMachine struct {
	pc       int           // 程序计数器
	stack    *Stack        // 栈
	heap     *Heap         // 堆
	static   *Static       // 静态区
	constant []interface{} // 常量池
	funcList []Function    // 函数引用列表
	codeList []byte        // 字节码
}

func NewVirtualMachine(
	constant []interface{},
	variableList []Object,
	functionList []*GoGoFunction,
	codeList []byte,
) *VirtualMachine {
	vm := &VirtualMachine{
		stack:    NewStack(),
		heap:     NewHeap(),
		static:   NewStatic(),
		constant: constant,
		funcList: make([]Function, 0),
		codeList: codeList,
	}

	for _, value := range variableList {
		vm.static.Append(value)
	}

	// 添加原生函数
	vm.AddNativeFunctions()

	for _, f := range functionList {
		vm.funcList = append(vm.funcList, f)
	}

	return vm
}

//
// 虚拟机执行入口
//
func (vm *VirtualMachine) Execute() {
	vm.pc = 0

	codeList := vm.codeList

	vm.stack.Expand(codeList)
	vm.execute(nil, codeList)
}

func (vm *VirtualMachine) execute(caller *GoGoFunction, codeList []byte) {
	pc := vm.pc
	base := 0
	stack := vm.stack
	static := vm.static
	constant := vm.constant

	for pc < len(codeList) {
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
			stack.SetIntPlus(0, constant[index].(int))
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
			stack.SetFloatPlus(0, constant[index].(float64))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_PUSH_STRING:
			index := utils.Get2ByteInt(codeList[pc+1:])
			stack.SetStringPlus(0, constant[index].(string))
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
			stack.SetPlus(0, static.Get(index))
			vm.stack.stackPointer++
			pc += 3
		case OP_CODE_POP_STATIC:
			index := utils.Get2ByteInt(codeList[pc+1:])
			static.Set(index, stack.GetPlus(-1))
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
			stack.SetIntPlus(-2, utils.BoolToInt(stack.GetStringPlus(-2) != stack.GetStringPlus(-1)))
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
			switch callee := vm.funcList[funcIdx].(type) {
			case *GoGoNativeFunction:
				vm.InvokeNativeFunction(callee, &vm.stack.stackPointer)
				pc++
			case *GoGoFunction:
				vm.InvokeFunction(&caller, callee, &codeList, &pc, &vm.stack.stackPointer, &base)
			default:
				panic("TODO")
			}
		case OP_CODE_RETURN:
			vm.ReturnFunction(&caller, &codeList, &pc, &vm.stack.stackPointer, &base)
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

//
// 函数相关 执行原生函数
//
func (vm *VirtualMachine) InvokeNativeFunction(f *GoGoNativeFunction, spP *int) {
	resultList := f.Proc(vm, f.ParamCount, vm.stack.list[*spP-f.ParamCount-1:])

	*spP = *spP - f.ParamCount - f.ResultCount - 1

	for _, value := range resultList {
		vm.stack.Set(*spP, value)
		*spP++
	}
}

// 函数执行
// caller 调用者, 当前所属的函数调用域
// callee 调用函数
// codeListP 字节码指针,用于设置新的字节码
func (vm *VirtualMachine) InvokeFunction(
	caller **GoGoFunction,
	callee *GoGoFunction,
	codeListP *[]byte,
	pcP *int,
	spP *int,
	bpP *int,
) {
	// 拓展栈大小
	vm.stack.Expand(callee.CodeList)

	// 设置返回值信息
	callInfo := &ObjectCallInfo{
		caller:        *caller,
		callerAddress: *pcP,
		bp:            *bpP,
	}

	// 栈上保存返回信息
	vm.stack.Set(*spP-1, callInfo)

	//
	// 设置新环境
	//
	*caller = callee
	*codeListP = callee.CodeList
	*pcP = 0
	*bpP = *spP - 1

	// 初始化局部变量
	for _, v := range callee.VariableList {
		vm.stack.Set(*spP, v)
		*spP++
	}
}

// 返回值入栈,恢复调用栈
func (vm *VirtualMachine) ReturnFunction(
	callerP **GoGoFunction,
	codeListP *[]byte,
	pcP *int,
	spP *int,
	bpP *int,
) {
	caller := *callerP

	paramCount := caller.ParamCount
	resultCount := caller.ResultCount

	for i := 0; i < resultCount; i++ {
		vm.stack.Set(*bpP-paramCount-resultCount+i, vm.stack.Get(*spP-resultCount+i))
		*spP++
	}

	// 恢复调用栈
	callInfo := vm.stack.Get(*bpP).(*ObjectCallInfo)

	if callInfo.caller != nil {
		*codeListP = callInfo.caller.CodeList
	} else {
		*codeListP = vm.codeList
	}

	*callerP = callInfo.caller
	*spP = *bpP - paramCount
	*pcP = callInfo.callerAddress + 1
	*bpP = callInfo.bp
}

//
// New
//

func (vm *VirtualMachine) NewObjectArray(size int) Object {
	obj := NewObjectArray(size)

	vm.AddObject(obj)

	for i := 0; i < size; i++ {
		obj.Set(i, vm.stack.GetPlus(-i-1))
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

	for i := 0; i < size; i++ {
		obj.FieldList[i] = vm.stack.GetPlus(-i - 1)
	}

	return obj
}
