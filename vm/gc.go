package vm

var HeapThresholdSize = 1024 * 256

//
// 判断是否下需要gc
//
func (vm *VirtualMachine) Check() {
	if len(vm.heap.objectList) > vm.heap.currentThreshold {
		vm.GC()
		vm.heap.currentThreshold += HeapThresholdSize
	}
}

func (vm *VirtualMachine) GC() {
	vm.Mark()
	vm.Sweep()
}

//
// 标记
//
func (vm *VirtualMachine) Mark() {
	for _, obj := range vm.heap.objectList {
		obj.ResetMark()
	}

	// 静态区
	for _, v := range vm.static.list {
		staticValue, ok := v.(*StaticVariable)
		if ok {
			obj, ok := staticValue.Value.(Object)
			if ok && obj != nil {
				obj.Mark()
			}
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
// 清理
//
func (vm *VirtualMachine) Sweep() {
	newObjectList := []Object{}
	for _, obj := range vm.heap.objectList {
		if !obj.isMarked() {
			obj.Sweep()
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.objectList = newObjectList
}

//
// 创建对象
//

// 添加对象到堆, 用于垃圾回收
func (vm *VirtualMachine) AddObject(value Object) {
	vm.Check()
	value.ResetMark()
	vm.heap.Append(value)
}
