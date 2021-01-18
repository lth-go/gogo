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

	// TODO: 静态区

	for i := 0; i < vm.stack.stackPointer; i++ {
		mark(vm.stack.Get(i))
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
// 标记，取消标记
//
func mark(obj Object) {
	if obj == nil {
		return
	}
	obj.Mark()
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

func (vm *VirtualMachine) NewObjectArray(size int) Object {
	obj := NewObjectArray(size)

	// add heap
	vm.AddObject(obj)

	// init
	for i := 0; i < size; i++ {
		obj.Set(i, vm.stack.GetPlus(-size+i))
	}

	return obj
}
