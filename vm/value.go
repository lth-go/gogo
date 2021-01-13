package vm

// 虚拟机基本值接口
type Value interface {
	isPointer() bool
	setPointer(bool)
}

type ValueImpl struct {
	// 是否是指针
	pointerFlags bool
}

func (v *ValueImpl) isPointer() bool {
	return v.pointerFlags
}

func (v *ValueImpl) setPointer(b bool) {
	v.pointerFlags = b
}

//
// CallInfo 函数返回体
//
type CallInfo struct {
	ValueImpl
	// 调用的函数
	caller *GFunction
	// 保存执行函数前的pc
	callerAddress int
	// TODO
	base int
}

//
// IntValue
//
type IntValue struct {
	ValueImpl
	intValue int
}

func NewIntValue(value int) *IntValue {
	return &IntValue{
		intValue: value,
	}
}

//
// DoubleValue
//
type DoubleValue struct {
	ValueImpl
	doubleValue float64
}

func NewDoubleValue(value float64) *DoubleValue {
	return &DoubleValue{
		doubleValue: value,
	}
}

// 引用对象
type ObjectRef struct {
	ValueImpl
	data Object
}

type ObjectImpl struct {
	// gc用
	marked bool
}

func (obj *ObjectImpl) isMarked() bool {
	return obj.marked
}

func (obj *ObjectImpl) Mark() {
	obj.marked = true
}

func (obj *ObjectImpl) ResetMark() {
	obj.marked = false
}

func (obj *ObjectImpl) Sweep() {
	// TODO:
}

func (obj *ObjectImpl) Len() int {
	return 1
}

//
// object string
//
type ObjectString struct {
	ObjectImpl
	stringValue string
}

//
// object array interface
//
type ObjectArray interface {
	getArraySize() int
}

// array int
type ObjectArrayInt struct {
	ObjectImpl
	intArray []int
}

func (array *ObjectArrayInt) getArraySize() int {
	if array.intArray == nil {
		return -1
	}
	return len(array.intArray)
}

func (array *ObjectArrayInt) getInt(index int) int {
	checkArray(array, index)
	return array.intArray[index]
}

func (array *ObjectArrayInt) setInt(index int, value int) {
	checkArray(array, index)
	array.intArray[index] = value
}

// array double
type ObjectArrayDouble struct {
	ObjectImpl
	doubleArray []float64
}

func (obj *ObjectArrayDouble) getArraySize() int {
	if obj.doubleArray == nil {
		return -1
	}
	return len(obj.doubleArray)
}

func (obj *ObjectArrayDouble) getDouble(index int) float64 {
	checkArray(obj, index)

	return obj.doubleArray[index]
}
func (obj *ObjectArrayDouble) setDouble(index int, value float64) {
	checkArray(obj, index)

	obj.doubleArray[index] = value
}

// array object
type ObjectArrayObject struct {
	ObjectImpl
	objectArray []*ObjectRef
}

func (obj *ObjectArrayObject) getArraySize() int {
	if obj.objectArray == nil {
		return -1
	}
	return len(obj.objectArray)
}

func (obj *ObjectArrayObject) getObject(index int) *ObjectRef {
	checkArray(obj, index)
	return obj.objectArray[index]
}

func (obj *ObjectArrayObject) setObject(index int, value *ObjectRef) {
	checkArray(obj, index)
	obj.objectArray[index] = value
}

// utils
func checkArray(array ObjectArray, index int) {
	if array == nil {
		vmError(NULL_POINTER_ERR)
		return
	}

	arraySize := array.getArraySize()
	if arraySize < 0 || index >= arraySize {
		vmError(INDEX_OUT_OF_BOUNDS_ERR, index, arraySize)
	}
}
