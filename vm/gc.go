package vm

var HeapThresholdSize = 1024 * 256

type objectType int

const (
	stringObjectType objectType = iota
)

//
// 判断是否下需要gc
//
func (vm *VmVirtualMachine) checkGC() {
	if len(vm.heap.objectList) > vm.heap.currentThreshold {
		vm.garbageCollect()

		vm.heap.currentThreshold += HeapThresholdSize
	}
}

//
// 创建对象
//
func (vm *VmVirtualMachine) addObject(value VmObject) {
	vm.checkGC()
	value.setMark(false)
	vm.heap.append(value)
}

func (vm *VmVirtualMachine) newStringObject() *VmObjectString {
	ret := &VmObjectString{}
	vm.addObject(ret)

	return ret
}

// 新建字符对象
func (vm *VmVirtualMachine) createStringObject(str string) VmObject {
	ret := vm.newStringObject()
	ret.stringValue = str

	return ret
}

// 连接字符对象
func (vm *VmVirtualMachine) chainStringObject(str1 VmObject, str2 VmObject) VmObject {
	var left, right string
	if str1 == nil {
		left = "null"
	} else {
		left = str1.getString()
	}

	if str2 == nil {
		right = "null"
	} else {
		right = str2.getString()
	}
	
	str := left + right
	ret := vm.createStringObject(str)
	return ret
}

// array
func (vm *VmVirtualMachine) createArrayInt(size int) *VmObjectArrayInt {
	ret := &VmObjectArrayInt{intArray: make([]int, size)}
	vm.addObject(ret)

	return ret
}

func (vm *VmVirtualMachine) createArrayDouble(size int) *VmObjectArrayDouble {
	ret := &VmObjectArrayDouble{doubleArray: make([]float64, size)}
	vm.addObject(ret)

	return ret
}

func (vm *VmVirtualMachine) createArrayObject(size int) *VmObjectArrayObject {
	ret := &VmObjectArrayObject{objectArray: make([]VmObject, size)}
	vm.addObject(ret)

	return ret
}

//
// 标记，取消标记
//
func mark(obj VmObject) {
	obj.setMark(true)

	arrayObj, ok := obj.(*VmObjectArrayObject)
	if ok {
		for _, subObj := range arrayObj.objectArray {
			mark(subObj)
		}
	}
}

func resetMark(obj VmObject) {
	obj.setMark(false)
}

//
// 标记
//
// TODO
func (vm *VmVirtualMachine) markObjects() {
	for _, obj := range vm.heap.objectList {
		resetMark(obj)
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
func (vm *VmVirtualMachine) disposeObject(obj VmObject) {
	switch o := obj.(type) {
	case *VmObjectString:
		//
	case *VmObjectArrayInt:
		o.intArray = nil
	case *VmObjectArrayDouble:
		o.doubleArray = nil
	case *VmObjectArrayObject:
		o.objectArray = nil
	default:
		panic("TODO")
	}

	obj = nil
}

//
// 清理
//
func (vm *VmVirtualMachine) sweepObjects() {
	newObjectList := []VmObject{}
	for _, obj := range vm.heap.objectList {
		if !obj.isMarked() {
			vm.disposeObject(obj)
		} else {
			newObjectList = append(newObjectList, obj)
		}
	}
	vm.heap.objectList = newObjectList
}

func (vm *VmVirtualMachine) garbageCollect() {
	vm.markObjects()
	vm.sweepObjects()
}
