package vm

//
// VmValue
//
// 虚拟机基本值接口
type VmValue interface {
	isPointer() bool
	setPointer(bool)
}

//
// VmValueImpl
//
type VmValueImpl struct {
	// 是否是指针
	pointerFlags bool
}

func (v *VmValueImpl) isPointer() bool {
	return v.pointerFlags
}

func (v *VmValueImpl) setPointer(b bool) {
	v.pointerFlags = b
}

//
// CallInfo 函数返回体
//
type CallInfo struct {
	VmValueImpl

	// 调用的函数
	caller *GFunction
	// 保存执行函数前的pc
	callerAddress int
	// TODO
	base int
}

//
// VmIntValue
//
type VmIntValue struct {
	VmValueImpl
	intValue int
}

func NewIntValue(value int) *VmIntValue {
	return &VmIntValue{
		intValue: value,
	}
}

//
// VmDoubleValue
//
type VmDoubleValue struct {
	VmValueImpl
	doubleValue float64
}

func NewDoubleValue(value float64) *VmDoubleValue {
	return &VmDoubleValue{
		doubleValue: value,
	}
}

//
// VmObjectValue
//
type VmObjectValue struct {
	VmValueImpl

	objectValue VmObject
}

func NewObjectValue(value VmObject) *VmObjectValue {
	return &VmObjectValue{
		objectValue: value,
	}
}

//
// VmObject
//
// 虚拟机对象接口, 包含string,
type VmObject interface {
	isMarked() bool
	setMark(bool)
}

type VmObjectImpl struct {
	// gc用
	marked bool
}

func (obj *VmObjectImpl) isMarked() bool {
	return obj.marked
}

func (obj *VmObjectImpl) setMark(m bool) {
	obj.marked = m
}

//
// object string
//
type VmObjectString struct {
	VmObjectImpl

	stringValue string
}

//
// object array interface
//
type VmObjectArray interface {
	getArraySize() int
}

// array int
type VmObjectArrayInt struct {
	VmObjectImpl
	intArray []int
}

func (obj *VmObjectArrayInt) getArraySize() int {
	if obj.intArray == nil {
		return -1
	}
	return len(obj.intArray)
}

func (array *VmObjectArrayInt) getInt(index int) int {
	checkArray(array, index)

	return array.intArray[index]
}

func (array *VmObjectArrayInt) setInt(index int, value int) {
	checkArray(array, index)

	array.intArray[index] = value
}

// array double
type VmObjectArrayDouble struct {
	VmObjectImpl
	doubleArray []float64
}

func (obj *VmObjectArrayDouble) getArraySize() int {
	if obj.doubleArray == nil {
		return -1
	}
	return len(obj.doubleArray)
}

func (array *VmObjectArrayDouble) getDouble(index int) float64 {
	checkArray(array, index)

	return array.doubleArray[index]
}
func (array *VmObjectArrayDouble) setDouble(index int, value float64) {
	checkArray(array, index)

	array.doubleArray[index] = value
}

// array object
type VmObjectArrayObject struct {
	VmObjectImpl
	objectArray []VmObject
}

func (obj *VmObjectArrayObject) getArraySize() int {
	if obj.objectArray == nil {
		return -1
	}
	return len(obj.objectArray)
}

func (array *VmObjectArrayObject) getObject(index int) VmObject {
	checkArray(array, index)

	return array.objectArray[index]
}

func (array *VmObjectArrayObject) setObject(index int, value VmObject) {
	checkArray(array, index)

	array.objectArray[index] = value
}

// utils
func checkArray(array VmObjectArray, index int) {

	if array == nil {
		vmError(NULL_POINTER_ERR)
		return
	}

	arraySize := array.getArraySize()
	if arraySize < 0 || index >= arraySize {
		vmError(INDEX_OUT_OF_BOUNDS_ERR, index, arraySize)
	}
}
