package vm

import (
	"fmt"
)

//
// 虚拟机
//
type VmVirtualMachine struct {
	// 栈
	stack *Stack
	// 堆
	heap *Heap

	// 全局函数列表
	function []Function
	// 解释器
	executable *ExecutableEntry

	// 类
	classList []*ExecClass

	// 程序计数器
	pc int

	// 当前exe
	currentExecutable *ExecutableEntry
	// 当前函数
	currentFunction *GFunction

	// 从exe传来的数据
	executableList *ExecutableList

	// exe列表
	executableEntryList []*ExecutableEntry

	// 顶层exe
	topLevel *ExecutableEntry
}

func NewVirtualMachine() *VmVirtualMachine {
	vm := &VmVirtualMachine{
		stack:             NewStack(),
		heap:              NewHeap(),
		static:            NewStatic(),
		function:          []Function{},
		executable:        nil,
		currentExecutable: nil,
	}
	vm.AddNativeFunctions()

	StVirtualMachine = vm

	return vm
}

var StVirtualMachine *VmVirtualMachine

func getVirtualMachine() *VmVirtualMachine {
	return StVirtualMachine
}

//
// 一些初始化操作
//
// TODO
func (vm *VmVirtualMachine) SetExecutableList(exelist *ExecutableList) {
	vm.executableList = exelist

	for _, exe := range exelist.List {
		vm.AddExecutable(vm, exe, (exe == exeList.TopLevel))
	}

}

// 虚拟机添加解释器
func (vm *VmVirtualMachine) AddExecutable(exe *Executable, isTopLevel bool) {

	newEntry := &ExecutableEntry{executable: exe}

	vm.executableEntryList = append(vm.executableEntryList, newEntry)

	vm.addFunctions(exe)

	vm.addClasses(newEntry)

	vm.convertCode(exe, exe.CodeList, nil)

	for _, f := range exe.FunctionList {
		vm.convertCode(exe, f.CodeList, f)
	}

    addStaticVariables(newEntry, exe)

	if isTopLevel {
		vm.topLevel = newEntry
	}
}

func addStaticVariables(entry *ExecutableEntry, exe *Executable) {

	entry.static = NewStatic()

	for _, value := range exe.GlobalVariableList {
		// TODO
        //if value.VmTypeSpecifier.BasicType == StringType {
        //    entry.static_v.variable[i].object = dvm_null_object_ref;
        //}

		entry.static.append(initializeValue(value.typeSpecifier))
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

func (vm *VmVirtualMachine) addClasses(ee *ExecutableEntry) {
	exe := ee.executable

	// 循环判断有无重复定义
	for _, exeClass := range exe.ClassDefinitionList {
		for dest_idx, vmClass := range vm.classList {
			// 有重复定义
			if vmClass.name == exeClass.Name && vmClass.getPackageName() == exeClass.getPackageName() {
				var package_name string

				if vmClass.packageName {
					package_name = vmClass.packageName
				} else {
					package_name = ""
				}
				vmError(CLASS_MULTIPLE_DEFINE_ERR, package_name, vmClass.name)
			}
		}
	}

	old_class_count = len(vm.classList)

	dest_idx := len(vm.classList)
	for _, exeClass := range exe.ClassDefinitionList {

		newVmClass := &ExecClass{}
		vm.classList = append(vm.classList, newVmClass)

		newVmClass.vmClass = exeClass
		newVmClass.executable = ee
		newVmClass.name = exeClass.name
		newVmClass.classIndex = dest_idx

		if exeClass.PackageName {
			newVmClass.PackageName = exeClass.PackageName
		} else {
			newVmClass.PackageName = ""
		}

		vm.addClass(exe, exeClass, newVmClass)

		dest_idx++
	}

	// 设置父类
	set_super_class(vm, exe, old_class_count)
}

func (vm *VmVirtualMachine) addClass(exe *Executable, src *VmClass, dest *ExecClass) {
	add_fields(exe, src, dest)
	add_methods(vm, exe, src, dest)
}

func add_fields(exe *Executable, src *VmClass, dest *ExecClass) {
	field_count := 0

	for pos := src; ; {
		field_count += len(pos.FieldList)
		if pos.SuperClass == nil {
			break
		}
		pos = search_class_from_executable(exe, pos.SuperClass.Name)
	}

	dest.fieldTypeList = make([]*VmTypeSpecifier, field_count)
	set_field_types(exe, src, dest.fieldTypeList, 0)
}

func set_field_types(exe *Executable, pos *VmClass, field_type []*VmTypeSpecifier, index int) int {

	if pos.SuperClass != nil {
		next := search_class_from_executable(exe, pos.SuperClass.Name)
		index = set_field_types(exe, next, field_type, index)
	}
	for _, field := range pos.FieldList {
		field_type[index] = field.Typ
		index++
	}

	return index
}

func add_methods(vm *VmVirtualMachine, exe *Executable, src *VmClass, dest *ExecClass) {

	v_table := alloc_v_table(dest)
	add_method(vm, exe, src, v_table)
	// TODO
	dest.class_table = v_table
}

func add_method(vm *VmVirtualMachine, exe *Executable, pos *VmClass, v_table *VmVTable) int {

	if pos.SuperClass != nil {
		next := search_class_from_executable(exe, pos.SuperClass.Name)
		super_method_count = add_method(vm, exe, next, v_table)
	}

	method_count := super_method_count
	for i = 0; i < pos.method_count; i++ {
		for j = 0; j < super_method_count; j++ {
			if pos.method[i].name == v_table.table[j].name {
				set_v_table(vm, pos, &pos.method[i], &v_table.table[j], false)
				break
			}
		}
		if j == super_method_count {
			v_table.table = MEM_realloc(v_table.table, sizeof(VTableItem)*(method_count+1))
			set_v_table(vm, pos, &pos.method[i], &v_table.table[method_count], true)
			method_count++
			v_table.table_size = method_count
		}
	}

	return method_count
}

func set_v_table(vm *VmVirtualMachine, class_p *VmClass, src *VmMethod, dest *VTableItem, set_name bool) {

	if set_name {
		dest.name = src.Name
	}

	func_name := vm_create_method_function_name(class_p.name, src.name)
	func_idx := search_function(dvm, class_p.package_name, func_name)

	if func_idx == FUNCTION_NOT_FOUND {
		vmError(FUNCTION_NOT_FOUND_ERR, func_name)
	}

	dest.index = func_idx
}

func set_super_class(vm *VmVirtualMachine, exe *Executable, old_class_count int) {
	for class_idx := old_class_count; class_idx < len(vm.classList); class_idx++ {
		vmClass := vm.classList[class_idx]

		exeClass := search_class_from_executable(exe, vmClass.name)
		if exeClass.super_class == nil {
			vmClass.super_class = nil
		} else {
			super_class_index = search_class(vm, exeClass.SuperClass.PackageName, exeClass.SuperClass.Name)
			vmClass.superClass = vm.class[super_class_index]
		}
	}
}

func search_class_from_executable(exe *Executable, name string) *VmClass {

	for _, exeClass := range exe.ClassDefinitionList {
		if exeClass.Name == name {
			return exeClass
		}
	}

	panic("TODO")
}

func search_class(vm *VmVirtualMachine, package_name string, name string) int {

	for i, vmClass := range vm.classList {
		if vmClass.packageName == package_name && vmClass.name == name {
			return i
		}
	}
	vmError(CLASS_NOT_FOUND_ERR, name)
	return 0
}

//
// 虚拟机执行入口
//
func (vm *VmVirtualMachine) Execute() {
	vm.pc = 0
	vm.currentExecutable = vm.executable
	vm.currentFunction = nil

	vm.stack.expand(vm.executable.CodeList)

	vm.execute(nil, vm.executable.CodeList)
}

func (vm *VmVirtualMachine) execute(gFunc *GFunction, codeList []byte) {
	var base int

	stack := vm.stack
	static := vm.static
	exe := vm.currentExecutable

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
		case VM_PUSH_NULL:
			stack.writeObject(0, nil)
			vm.stack.stackPointer++
			pc++
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
		case VM_PUSH_STACK_OBJECT:
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
		case VM_POP_STACK_OBJECT:
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
		case VM_PUSH_STATIC_OBJECT:
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
		case VM_POP_STATIC_OBJECT:
			index := get2ByteInt(codeList[pc+1:])
			static.setObject(index, stack.getObject(-1))
			vm.stack.stackPointer--
			pc += 3
		case VM_PUSH_ARRAY_INT:
			array := stack.getObject(-2).(*VmObjectArrayInt)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			intValue := array.getInt(index)

			stack.writeInt(-2, intValue)
			vm.stack.stackPointer--
			pc++
		case VM_PUSH_ARRAY_DOUBLE:
			array := stack.getObject(-2).(*VmObjectArrayDouble)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			doubleValue := array.getDouble(index)

			stack.writeDouble(-2, doubleValue)
			vm.stack.stackPointer--
			pc++
		case VM_PUSH_ARRAY_OBJECT:
			array := stack.getObject(-2).(*VmObjectArrayObject)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			object := array.getObject(index)

			stack.writeObject(-2, object)
			vm.stack.stackPointer--
			pc++
		case VM_POP_ARRAY_INT:
			value := stack.getInt(-3)
			array := stack.getObject(-2).(*VmObjectArrayInt)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			array.setInt(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_POP_ARRAY_DOUBLE:
			value := stack.getDouble(-3)
			array := stack.getObject(-2).(*VmObjectArrayDouble)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			array.setDouble(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_POP_ARRAY_OBJECT:
			value := stack.getObject(-3)
			array := stack.getObject(-2).(*VmObjectArrayObject)
			index := stack.getInt(-1)

			vm.restore_pc(exe, gFunc, pc)
			array.setObject(index, value)
			vm.stack.stackPointer -= 3
			pc++
		case VM_ADD_INT:
			stack.writeInt(-2, stack.getInt(-2)+stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_DOUBLE:
			stack.writeDouble(-2, stack.getDouble(-2)+stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_ADD_STRING:
			stack.writeObject(-2, vm.chainStringObject(stack.getObject(-2), stack.getObject(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_INT:
			stack.writeInt(-2, stack.getInt(-2)-stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_SUB_DOUBLE:
			stack.writeDouble(-2, stack.getDouble(-2)-stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_INT:
			stack.writeInt(-2, stack.getInt(-2)*stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MUL_DOUBLE:
			stack.writeDouble(-2, stack.getDouble(-2)*stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_INT:
			if stack.getInt(-1) == 0 {
				vmError(DIVISION_BY_ZERO_ERR)
			}
			stack.writeInt(-2, stack.getInt(-2)/stack.getInt(-1))
			vm.stack.stackPointer--
			pc++
		case VM_DIV_DOUBLE:
			stack.writeDouble(-2, stack.getDouble(-2)/stack.getDouble(-1))
			vm.stack.stackPointer--
			pc++
		case VM_MINUS_INT:
			stack.writeInt(-1, -stack.getInt(-1))
			pc++
		case VM_MINUS_DOUBLE:
			stack.writeDouble(-1, -stack.getDouble(-1))
			pc++
		case VM_CAST_INT_TO_DOUBLE:
			stack.writeDouble(-1, float64(stack.getInt(-1)))
			pc++
		case VM_CAST_DOUBLE_TO_INT:
			stack.writeInt(-1, int(stack.getDouble(-1)))
			pc++
		case VM_CAST_BOOLEAN_TO_STRING:
			if stack.getInt(-1) != 0 {
				stack.writeObject(-1, vm.createStringObject("true"))
			} else {
				stack.writeObject(-1, vm.createStringObject("false"))
			}
			pc++
		case VM_CAST_INT_TO_STRING:
			// TODO 啥意思
			vm.restore_pc(exe, gFunc, pc)
			buf := fmt.Sprintf("%d", stack.getInt(-1))
			stack.writeObject(-1, vm.createStringObject(buf))
			pc++
		case VM_CAST_DOUBLE_TO_STRING:
			// TODO 啥意思
			vm.restore_pc(exe, gFunc, pc)
			buf := fmt.Sprintf("%f", stack.getDouble(-1))
			stack.writeObject(-1, vm.createStringObject(buf))
			pc++
		case VM_EQ_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) == stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) == stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_OBJECT:
			stack.writeInt(-2, boolToInt(stack.getObject(-2) == stack.getObject(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_EQ_STRING:
			stack.writeInt(-2, boolToInt(!(stack.getString(-2) == stack.getString(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_GT_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) > stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) > stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GT_STRING:
			stack.writeInt(-2, boolToInt(stack.getString(-2) > stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) >= stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) >= stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_GE_STRING:
			stack.writeInt(-2, boolToInt(stack.getString(-2) >= stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) < stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) < stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LT_STRING:
			stack.writeInt(-2, boolToInt(stack.getString(-2) < stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) <= stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) <= stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LE_STRING:
			stack.writeInt(-2, boolToInt(stack.getString(-2) <= stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_INT:
			stack.writeInt(-2, boolToInt(stack.getInt(-2) != stack.getInt(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_DOUBLE:
			stack.writeInt(-2, boolToInt(stack.getDouble(-2) != stack.getDouble(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_OBJECT:
			stack.writeInt(-2, boolToInt(stack.getObject(-2) != stack.getObject(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_NE_STRING:
			stack.writeInt(-2, boolToInt(stack.getString(-2) != stack.getString(-1)))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_AND:
			stack.writeInt(-2, boolToInt(intToBool(stack.getInt(-2)) && intToBool(stack.getInt(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_OR:
			stack.writeInt(-2, boolToInt(intToBool(stack.getInt(-2)) || intToBool(stack.getInt(-1))))
			vm.stack.stackPointer--
			pc++
		case VM_LOGICAL_NOT:
			stack.writeInt(-1, boolToInt(!intToBool(stack.getInt(-1))))
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
			value := get2ByteInt(codeList[pc+1:])
			stack.writeInt(0, value)
			vm.stack.stackPointer++
			pc += 3
		case VM_INVOKE:
			funcIdx := stack.getInt(-1)
			switch f := vm.function[funcIdx].(type) {
			case *NativeFunction:
				vm.invokeNativeFunction(f, &vm.stack.stackPointer)
				pc++
			case *GFunction:
				vm.invokeGFunction(&gFunc, f, &codeList, &pc, &vm.stack.stackPointer, &base, exe)
			default:
				panic("TODO")
			}
		case VM_RETURN:
			vm.returnFunction(&gFunc, &codeList, &pc, &vm.stack.stackPointer, &base, exe)
		case VM_NEW_ARRAY:
			dim := int(codeList[pc+1])
			typ := exe.TypeSpecifierList[get2ByteInt(codeList[pc+2:])]

			vm.restore_pc(exe, gFunc, pc)
			array := vm.create_array(dim, typ)
			vm.stack.stackPointer -= dim

			stack.writeObject(0, array)
			vm.stack.stackPointer++
			pc += 4
		case VM_NEW_ARRAY_LITERAL_INT:
			size := get2ByteInt(codeList[pc+1:])

			vm.restore_pc(exe, gFunc, pc)
			array := vm.create_array_literal_int(size)
			vm.stack.stackPointer -= size
			stack.writeObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		case VM_NEW_ARRAY_LITERAL_DOUBLE:
			size := get2ByteInt(codeList[pc+1:])

			vm.restore_pc(exe, gFunc, pc)
			array := vm.create_array_literal_double(size)
			vm.stack.stackPointer -= size
			stack.writeObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		case VM_NEW_ARRAY_LITERAL_OBJECT:
			size := get2ByteInt(codeList[pc+1:])

			vm.restore_pc(exe, gFunc, pc)
			array := vm.create_array_literal_object(size)
			vm.stack.stackPointer -= size
			stack.writeObject(0, array)
			vm.stack.stackPointer++
			pc += 3
		default:
			panic("TODO")
		}
	}
}

func (vm *VmVirtualMachine) initializeValue(typ *VmTypeSpecifier) VmValue {
	var value VmValue

	if typ.DeriveList != nil && len(typ.DeriveList) > 0 {
		_, ok := typ.DeriveList[0].(*VmArrayDerive)
		if !ok {
			panic("TODO")
		}
		value = &VmObjectValue{objectValue: nil}
		return value
	}

	switch typ.BasicType {
	case BooleanType:
		fallthrough
	case IntType:
		value = &VmIntValue{intValue: 0}
	case DoubleType:
		value = &VmDoubleValue{doubleValue: 0.0}
	case StringType:
		value = &VmObjectValue{objectValue: vm.createStringObject("")}
	case NullType:
		fallthrough
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
		vm.stack.stack[sp_idx] = vm.initializeValue(f.LocalVariableList[i].TypeSpecifier)

		if f.LocalVariableList[i].TypeSpecifier.BasicType == StringType {
			vm.stack.stack[i].setPointer(true)
		}
		sp_idx++
	}
}

// 修正转换code
func (vm *VmVirtualMachine) convertCode(exe *Executable, codeList []byte, f *VmFunction) {
	var dest_idx int

	for i := 0; i < len(codeList); i++ {
		code := codeList[i]
		switch code {
		// 函数内的本地声明
		case VM_PUSH_STACK_INT, VM_POP_STACK_INT,
			VM_PUSH_STACK_DOUBLE, VM_POP_STACK_DOUBLE,
			VM_PUSH_STACK_OBJECT, VM_POP_STACK_OBJECT:

			// TODO
			if f == nil {
				for _, lineNumber := range exe.LineNumberList {
					if lineNumber.StartPc == i {
						panic(fmt.Sprintf("Line: %d, func == nil", lineNumber.LineNumber))
					}
				}
				panic("can't find line, need debug!!!")
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
			funcIdx := vm.SearchFunction(exe.FunctionList[idx_in_exe].Name)
			set2ByteInt(codeList[i+1:], funcIdx)
		}

		info := &OpcodeInfo[code]
		for _, p := range []byte(info.Parameter) {
			switch p {
			case 'b':
				i++
			case 's':
				fallthrough
			case 'p':
				i += 2
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
	vmError(FUNCTION_NOT_FOUND_ERR, name)
	return 0
}

//
// 函数相关
//
// 执行原生函数
func (vm *VmVirtualMachine) invokeNativeFunction(f *NativeFunction, sp_p *int) {

	stack := vm.stack.stack
	sp := *sp_p

	ret := f.proc(vm, f.argCount, stack[sp-f.argCount-1:])

	stack[sp-f.argCount-1] = ret

	*sp_p = sp - f.argCount
}

// 函数执行
func (vm *VmVirtualMachine) invokeGFunction(caller **GFunction, callee *GFunction,
	code_p *[]byte, pc_p *int, sp_p *int, base_p *int,
	exe *Executable) {
	// caller 调用者, 当前所属的函数调用域

	// callee 要调用的函数的基本信息

	exe = callee.Executable
	// 包含调用函数的全部信息
	callee_p := exe.FunctionList[callee.Index]

	// 拓展栈大小
	vm.stack.expand(callee_p.CodeList)

	// 设置返回值信息
	callInfo := &CallInfo{
		caller:        *caller,
		callerAddress: *pc_p,
		base:          *base_p,
	}

	// 栈上保存返回信息
	vm.stack.stack[*sp_p-1] = callInfo

	// 设置base
	*base_p = *sp_p - len(callee_p.ParameterList) - 1

	// 设置调用者
	*caller = callee

	// 初始化参数
	vm.initializeLocalVariables(callee_p, *sp_p)

	// 设置栈位置
	*sp_p += len(callee_p.LocalVariableList)
	*pc_p = 0

	// 设置字节码为函数的字节码
	*code_p = callee_p.CodeList
}

func (vm *VmVirtualMachine) returnFunction(func_p **GFunction, code_p *[]byte, pc_p *int, sp_p *int, base_p *int, exe *Executable) {

	returnValue := vm.stack.stack[(*sp_p)-1]

	callee_p := exe.FunctionList[(*func_p).Index]
	callInfo := vm.stack.stack[*sp_p-1-len(callee_p.LocalVariableList)-1].(*CallInfo)

	if callInfo.caller != nil {
		exe = callInfo.caller.Executable
		caller_p := exe.FunctionList[callInfo.caller.Index]
		*code_p = caller_p.CodeList
	} else {
		// TODO为什么没有返回值
		exe = vm.executable
		*code_p = vm.executable.CodeList
	}
	*func_p = callInfo.caller

	*pc_p = callInfo.callerAddress + 1
	*base_p = callInfo.base

	*sp_p -= (len(callee_p.LocalVariableList) + 1 + len(callee_p.ParameterList))

	vm.stack.stack[*sp_p-1] = returnValue
}

func (vm *VmVirtualMachine) create_array_sub(dim int, dim_index int, typ *VmTypeSpecifier) VmObject {
	var ret VmObject

	size := vm.stack.getInt(-dim)

	if dim_index == (len(typ.DeriveList) - 1) {
		switch typ.BasicType {
		case BooleanType:
			fallthrough
		case IntType:
			ret = vm.createArrayInt(size)
		case DoubleType:
			ret = vm.createArrayDouble(size)
		case StringType:
			ret = vm.createArrayObject(size)
		case NullType:
			fallthrough
		default:
			panic("TODO")
		}
	} else if _, ok := typ.DeriveList[dim_index].(*VmFunctionDerive); ok {
		// BUG ?
		ret = nil
	} else {
		ret = vm.createArrayObject(size)
		if dim_index < dim-1 {
			vm.stack.writeObject(0, ret)
			vm.stack.stackPointer++

			for i := 0; i < size; i++ {
				child := vm.create_array_sub(dim, dim_index+1, typ)
				ret.(*VmObjectArrayObject).setObject(i, child)
			}
			vm.stack.stackPointer--
		}
	}
	return ret
}

func (vm *VmVirtualMachine) create_array(dim int, typ *VmTypeSpecifier) VmObject {
	return vm.create_array_sub(dim, 0, typ)
}

func (vm *VmVirtualMachine) create_array_literal_int(size int) VmObject {

	array := vm.createArrayInt(size)
	for i := 0; i < size; i++ {
		array.intArray[i] = vm.stack.getInt(-size + i)
	}

	return array
}

func (vm *VmVirtualMachine) create_array_literal_double(size int) VmObject {

	array := vm.createArrayDouble(size)
	for i := 0; i < size; i++ {
		array.doubleArray[i] = vm.stack.getDouble(-size + i)
	}

	return array
}

func (vm *VmVirtualMachine) create_array_literal_object(size int) VmObject {
	array := vm.createArrayObject(size)
	for i := 0; i < size; i++ {
		array.objectArray[i] = vm.stack.getObject(-size + i)
	}

	return array
}

func (vm *VmVirtualMachine) restore_pc(exe *Executable, function *GFunction, pc int) {
	vm.currentExecutable = exe
	vm.currentFunction = function
	vm.pc = pc
}
