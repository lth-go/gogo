package vm

func check_array(array VmObject, index int, exe *Executable, function *GFunction, pc int) {
	if array == nil {
		vmError(exe, function, pc, NULL_POINTER_ERR)
		return
	}

	arraySize := array.getArraySize()
	if index < 0 || index >= arraySize {
		vmError(exe, function, pc, INDEX_OUT_OF_BOUNDS_ERR, index, arraySize)
	}
}

func (vm *VmVirtualMachine) array_get_int(array *VmObjectArrayInt, index int) int {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	return array.intArray[index]
}

func (vm *VmVirtualMachine) array_get_double(array *VmObjectArrayDouble, index int) float64 {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	return array.doubleArray[index]
}

func (vm *VmVirtualMachine) array_get_object(array *VmObjectArrayObject, index int) VmObject {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	return array.objectArray[index]
}

func (vm *VmVirtualMachine) array_set_int(array *VmObjectArrayInt, index int, value int) {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	array.intArray[index] = value
}

func (vm *VmVirtualMachine) array_set_double(array *VmObjectArrayDouble, index int, value float64) {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	array.doubleArray[index] = value
}

func (vm *VmVirtualMachine) array_set_object(array *VmObjectArrayObject, index int, value VmObject) {
	check_array(array, index, vm.currentExecutable, vm.currentFunction, vm.pc)

	array.objectArray[index] = value
}
