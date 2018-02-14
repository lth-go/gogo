package vm

//
// VmValue
//
// 虚拟机基本值接口
type VmValue interface {
	getIntValue() int
	setIntValue(int)

	getDoubleValue() float64
	setDoubleValue(float64)

	getObjectValue() VmObject
	setObjectValue(VmObject)

	isPointer() bool
	setPointer(bool)
}

// VmValueImpl

type VmValueImpl struct {
	// 是否是指针
	pointerFlags bool
}

func (v *VmValueImpl) getIntValue() int {
	panic("error")
}

func (v *VmValueImpl) setIntValue(value int) {
	panic("error")
}

func (v *VmValueImpl) getDoubleValue() float64 {
	panic("error")
}

func (v *VmValueImpl) setDoubleValue(value float64) {
	panic("error")
}

func (v *VmValueImpl) getObjectValue() VmObject {
	panic("VmValue: 数据类型错误, 无法获取VmObject")
}

func (v *VmValueImpl) setObjectValue(value VmObject) {
	panic("error")
}

func (v *VmValueImpl) isPointer() bool {
	return v.pointerFlags
}

func (v *VmValueImpl) setPointer(b bool) {
	v.pointerFlags = b
}

// CallInfo
// 函数返回体
type CallInfo struct {
	VmValueImpl

	// 调用的函数
	caller *GFunction
	// 保存执行函数前的pc
	callerAddress int
	// TODO
	base int
}

// VmIntValue
type VmIntValue struct {
	VmValueImpl
	intValue int
}

func NewIntValue(value int) *VmIntValue {
	return &VmIntValue{
		intValue: value,
	}
}

func (v *VmIntValue) getIntValue() int {
	return v.intValue
}

func (v *VmIntValue) setIntValue(value int) {
	v.intValue = value
}

// VmDoubleValue

type VmDoubleValue struct {
	VmValueImpl
	doubleValue float64
}

func NewDoubleValue(value float64) *VmDoubleValue {
	return &VmDoubleValue{
		doubleValue: value,
	}
}

func (v *VmDoubleValue) getDoubleValue() float64 {
	return v.doubleValue
}

func (v *VmDoubleValue) setDoubleValue(value float64) {
	v.doubleValue = value
}

// VmObjectValue

type VmObjectValue struct {
	VmValueImpl

	objectValue VmObject
}

func NewObjectValue(value VmObject) *VmObjectValue {
	return &VmObjectValue{
		objectValue: value,
	}
}

func (v *VmObjectValue) getObjectValue() VmObject {
	return v.objectValue
}

func (v *VmObjectValue) setObjectValue(value VmObject) {
	v.objectValue = value
}

//
// VmObject
//
// 虚拟机对象接口, 包含string,
type VmObject interface {
	isMarked() bool
	setMark(bool)

	getString() string
	setString(string)
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

func (obj *VmObjectImpl) getString() string {
	panic("TODO")
}

func (obj *VmObjectImpl) setString(v string) {
	panic("TODO")
}

type VmObjectString struct {
	VmObjectImpl

	stringValue string
}

func (obj *VmObjectString) getString() string {
	return obj.stringValue
}

func (obj *VmObjectString) setString(v string) {
	obj.stringValue = v
}
