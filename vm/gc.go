package vm

var HeapThresholdSize = 1024 * 256

//
// Check 判断是否下需要gc
//
func (vm *VirtualMachine) Check() {
	if len(vm.heap.list) > vm.heap.currentThreshold {
		vm.GC()
		vm.heap.currentThreshold += HeapThresholdSize
	}
}

func (vm *VirtualMachine) GC() {
	vm.Mark()
	vm.Sweep()
}

//
// Mark 标记
//
func (vm *VirtualMachine) Mark() {
	for _, obj := range vm.heap.list {
		obj.ResetMark()
	}

	// 静态区
	for _, obj := range vm.static.list {
		if obj != nil {
			obj.Mark()
		}
	}

	for i := 0; i < vm.stack.stackPointer; i++ {
		obj := vm.stack.Get(i)
		if obj != nil {
			obj.Mark()
		}
	}
}

//
// Sweep 清理
//
func (vm *VirtualMachine) Sweep() {
	newObjectList := []Object{}
	for _, obj := range vm.heap.list {
		if !obj.IsMarked() {
			obj.Sweep()
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.list = newObjectList
}

// AddObject 添加对象到堆, 用于垃圾回收
func (vm *VirtualMachine) AddObject(value Object) {
	vm.Check()
	value.ResetMark()
	vm.heap.Append(value)
}
