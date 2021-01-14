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

// array object
type ObjectArrayObject struct {
	ObjectBase
	objectArray []*ObjectRef
	List        []Object
}

func (obj *ObjectArrayObject) getArraySize() int {
	if obj.objectArray == nil {
		return -1
	}
	return len(obj.objectArray)
}

func (obj *ObjectArrayObject) getObject(index int) *ObjectRef {
	return obj.objectArray[index]
}

func (obj *ObjectArrayObject) setObject(index int, value *ObjectRef) {
	obj.objectArray[index] = value
}
