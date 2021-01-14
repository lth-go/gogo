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

	for _, exe := range vm.executableList {
		for _, variable := range exe.VariableList.VariableList {
			if variable.IsReferenceType() {
				mark(variable.Value.(*ObjectRef))
			}
		}
	}

	for i := 0; i < vm.stack.stackPointer; i++ {
		if vm.stack.Get(i).isPointer() {
			mark(vm.stack.GetObject(i))
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
			// TODO: 对象自身sweep
			obj.Sweep()
			vm.disposeObject(obj)
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.objectList = newObjectList
}

//
// 标记，取消标记
//
func mark(ref *ObjectRef) {
	obj := ref.data
	if obj == nil {
		return
	}
	obj.Mark()

	switch o := obj.(type) {
	case *ObjectArrayObject:
		for _, subObj := range o.objectArray {
			mark(subObj)
		}
	}
}

//
// 删除对象
//
func (vm *VirtualMachine) disposeObject(obj Object) {
	switch o := obj.(type) {
	case *ObjectString:
		//
	case *ObjectArrayInt:
		o.intArray = nil
	case *ObjectArrayDouble:
		o.doubleArray = nil
	case *ObjectArrayObject:
		o.objectArray = nil
	default:
		panic("TODO")
	}

	obj = nil
}

//
// 创建对象
//

//
// 添加对象到堆, 用于垃圾回收
//
func (vm *VirtualMachine) AddObject(value Object) {
	vm.Check()
	value.ResetMark()
	vm.heap.Append(value)
}

//
// string object
//
func (vm *VirtualMachine) createStringObject(str string) *ObjectRef {
	ret := &ObjectString{}
	vm.AddObject(ret)

	ret.stringValue = str

	ref := &ObjectRef{data: ret}

	return ref
}

//
// Array object
//
func (vm *VirtualMachine) createArrayInt(size int) *ObjectRef {
	obj := &ObjectArrayInt{intArray: make([]int, size)}
	vm.AddObject(obj)

	ref := &ObjectRef{data: obj}

	return ref
}

func (vm *VirtualMachine) createArrayDouble(size int) *ObjectRef {
	obj := &ObjectArrayDouble{doubleArray: make([]float64, size)}
	vm.AddObject(obj)

	ref := &ObjectRef{data: obj}

	return ref
}

func (vm *VirtualMachine) createArrayObject(size int) *ObjectRef {
	obj := &ObjectArrayObject{objectArray: make([]*ObjectRef, size)}
	vm.AddObject(obj)

	ref := &ObjectRef{data: obj}

	return ref
}

// 连接字符对象
func (vm *VirtualMachine) chainStringObject(str1, str2 *ObjectRef) *ObjectRef {
	var left, right string
	if str1.data == nil {
		left = "null"
	} else {
		left = str1.data.(*ObjectString).stringValue
	}

	if str2.data == nil {
		right = "null"
	} else {
		right = str2.data.(*ObjectString).stringValue
	}

	str := left + right
	ret := vm.createStringObject(str)
	return ret
}
