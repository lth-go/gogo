package vm
//
// Heap
//
type Heap struct {
	// TODO:阈值
	currentThreshold int
	objectList       []VmObject
}

var HeapThresholdSize = 1024 * 256

func (vm *VmVirtualMachine) check_gc() {
	if len(vm.heap.objectList) > vm.heap.currentThreshold {
		vm.garbage_collect()

		vm.heap.currentThreshold += HeapThresholdSize
	}
}

func (vm *VmVirtualMachine) alloc_object() VmObject{
	ret := &VmObjectString{}

	vm.check_gc()

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

func (vm *VmVirtualMachine) create_vm_string_i(str string) *VmObjectString{
	ret := vm.alloc_object_string()
	ret.stringValue = str
	ret.isLiteral = false

	return ret
}

func gc_mark(obj VmObject) { obj.setMark(true) }

func gc_reset_mark(obj VmObject) { obj.setMark(false) }

//
// 标记
//
// TODO
func (vm *VmVirtualMachine) gc_mark_objects() {
	for _, obj := range vm.heap.objectList {
		gc_reset_mark(obj)
	}

		for i, v := range vm.static.variableList {
		if vm.executable.globalVariableList[i].typeSpecifier.basicType == StringType {
			gc_mark(v)
		}
	}

	for i := 0; i < vm.stack.stackPointer; i++ {
		if (vm.stack.pointer_flags[i]) {
			gc_mark(vm.stack.stack[i].object)
		}
	}
}


//
// 删除对象
//
func  (vm *VmVirtualMachine)gc_dispose_object(obj VmObject) {
	delete(obj)
}

//
// 清理
//
func (vm *VmVirtualMachine) gc_sweep_objects() {
	vm_Object *obj;
	vm_Object *tmp;

	newObjectList := []VmObject{}
	for _, obj := range vm.heap.objectList {
		if !obj.isMarked {
			vm.gc_dispose_object(obj)
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.objectList = newObjectList
}

func (vm *VmVirtualMachine) garbage_collect() {
	vm.gc_mark_objects();
	vm.gc_sweep_objects();
}
