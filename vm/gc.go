package vm

var HeapThresholdSize = 1024 * 256

//
// 判断是否下需要gc
//
func (vm *VmVirtualMachine) check_gc() {
	if len(vm.heap.objectList) > vm.heap.currentThreshold {
		vm.garbage_collect()

		vm.heap.currentThreshold += HeapThresholdSize
	}
}

//
// 创建对象
//
func (vm *VmVirtualMachine) alloc_object() VmObject {
	ret := &VmObjectString{}

	vm.check_gc()

	ret.marked = false

	vm.heap.objectList = append(vm.heap.objectList, ret)

	return ret
}

func (vm *VmVirtualMachine) alloc_object_string() *VmObjectString {

	//check_gc(vm)
	ret := &VmObjectString{}

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

func (vm *VmVirtualMachine) create_vm_string_i(str string) *VmObjectString {
	ret := vm.alloc_object_string()
	ret.stringValue = str
	ret.isLiteral = false

	return ret
}

//
// 标记，取消标记
//
func mark(obj VmObject) { obj.setMark(true) }

func reset_mark(obj VmObject) { obj.setMark(false) }

//
// 标记
//
// TODO
func (vm *VmVirtualMachine) mark_objects() {
	for _, obj := range vm.heap.objectList {
		reset_mark(obj)
	}

	for _, v := range vm.static.variableList {
		if o, ok := v.(VmObject); ok {
			mark(o)
		}
	}

	for i := 0; i < vm.stack.stackPointer; i++ {
		if vm.stack.stack[i].isPointer() {
			o := vm.stack.stack[i].(VmObject)
			mark(o)
		}
	}
}

//
// 删除对象
//
func (vm *VmVirtualMachine) dispose_object(obj VmObject) {
	switch o := obj.(type) {
	case *VmObjectString:
		if !o.isLiteral {
			//
		}
	}
	obj = nil
}

//
// 清理
//
func (vm *VmVirtualMachine) sweep_objects() {
	newObjectList := []VmObject{}
	for _, obj := range vm.heap.objectList {
		if !obj.isMarked() {
			vm.dispose_object(obj)
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.objectList = newObjectList
}

func (vm *VmVirtualMachine) garbage_collect() {
	vm.mark_objects()
	vm.sweep_objects()
}
