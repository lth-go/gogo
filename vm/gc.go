package vm

var HeapThresholdSize = 1024 * 256

type objectType int

const (
	stringObjectType objectType = iota
)

//////////////////////////////
// 垃圾回收
//////////////////////////////

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

//////////////////////////////
// 创建对象
//////////////////////////////

//
// 添加对象到堆, 用于垃圾回收
//
func (vm *VmVirtualMachine) addObject(value VmObject) {
	vm.checkGC()
	value.setMark(false)
	vm.heap.append(value)
}

//////////////////////////////

//
// string object
//
func (vm *VmVirtualMachine) createStringObject(str string) *VmObjectRef {
	ret := &VmObjectString{}
	vm.addObject(ret)

	ret.stringValue = str

	ref := &VmObjectRef{data: ret}

	return ref
}

//
// Array object
//
func (vm *VmVirtualMachine) createArrayInt(size int) *VmObjectRef {
	obj := &VmObjectArrayInt{intArray: make([]int, size)}
	vm.addObject(obj)

	ref := &VmObjectRef{data: obj}

	return ref
}

func (vm *VmVirtualMachine) createArrayDouble(size int) *VmObjectRef {
	obj := &VmObjectArrayDouble{doubleArray: make([]float64, size)}
	vm.addObject(obj)

	ref := &VmObjectRef{data: obj}

	return ref
}

func (vm *VmVirtualMachine) createArrayObject(size int) *VmObjectRef {
	obj := &VmObjectArrayObject{objectArray: make([]*VmObjectRef, size)}
	vm.addObject(obj)

	ref := &VmObjectRef{data: obj}

	return ref
}

//
// class object
//
func (vm *VmVirtualMachine) createClassObject(classIndex int) *VmObjectRef {
	obj := &VmObjectClassObject{}
	vm.addObject(obj)

	execClass := vm.classList[classIndex]

	obj.fieldList = []VmValue{}
	for _, typ := range execClass.fieldTypeList {
		obj.fieldList = append(obj.fieldList, initializeValue(typ))
	}

	ref := &VmObjectRef{
		vTable: execClass.classTable,
		data:   obj,
	}

	return ref
}

// utils

// 连接字符对象
func (vm *VmVirtualMachine) chainStringObject(str1, str2 *VmObjectRef) *VmObjectRef {
	var left, right string
	if str1.data == nil {
		left = "null"
	} else {
		left = str1.data.(*VmObjectString).stringValue
	}

	if str2.data == nil {
		right = "null"
	} else {
		right = str2.data.(*VmObjectString).stringValue
	}

	str := left + right
	ret := vm.createStringObject(str)
	return ret
}
