package vm

import (
	"github.com/lth-go/gogo/utils"
)

// 虚拟机对象接口
type Object interface {
	IsMarked() bool // 是否设置标记位
	Mark()          // 设置标记位
	ResetMark()     // 重置标记位
	Sweep()         // 垃圾回收
	Len() int       // 计算堆阈值
	Hash() int      // 哈希
}

//
// ObjectBase
//
type ObjectBase struct {
	marked bool
}

func (obj *ObjectBase) IsMarked() bool {
	return obj.marked
}

func (obj *ObjectBase) Mark() {
	if obj == nil {
		return
	}
	obj.marked = true
}

func (obj *ObjectBase) ResetMark() {
	obj.marked = false
}

func (obj *ObjectBase) Sweep() {
	// TODO:
}

func (obj *ObjectBase) Len() int {
	return 1
}

func (obj *ObjectBase) Hash() int {
	return 0
}

//
// ObjectInt
//
type ObjectInt struct {
	ObjectBase
	Value int
}

func NewObjectInt(value int) *ObjectInt {
	return &ObjectInt{
		Value: value,
	}
}

//
// ObjectFloat
//
type ObjectFloat struct {
	ObjectBase
	Value float64
}

func NewObjectFloat(value float64) *ObjectFloat {
	return &ObjectFloat{
		Value: value,
	}
}

//
// ObjectString
//
type ObjectString struct {
	ObjectBase
	Value string
}

func NewObjectString(value string) *ObjectString {
	return &ObjectString{
		Value: value,
	}
}

//
// ObjectNil
//
type ObjectNil struct {
	ObjectBase
}

var NilObject = &ObjectNil{}

//
// ObjectArray
//
type ObjectArray struct {
	ObjectBase
	List []Object
}

func (obj *ObjectArray) Mark() {
	if obj == nil {
		return
	}

	obj.ObjectBase.Mark()

	for _, subObj := range obj.List {
		if subObj == nil {
			continue
		}
		subObj.Mark()
	}
}

func (obj *ObjectArray) ResetMark() {
	obj.ObjectBase.ResetMark()

	for _, subObj := range obj.List {
		if subObj == nil {
			continue
		}
		subObj.ResetMark()
	}
}

func (obj *ObjectArray) Sweep() {
	for _, subObj := range obj.List {
		subObj.Sweep()
	}
}

func (obj *ObjectArray) Len() int {
	return len(obj.List)
}

func (obj *ObjectArray) Set(index int, value Object) {
	obj.Check(index)
	obj.List[index] = value
}

func (obj *ObjectArray) SetInt(index int, value int) {
	obj.Set(index, NewObjectInt(value))
}

func (obj *ObjectArray) SetFloat(index int, value float64) {
	obj.Set(index, NewObjectFloat(value))
}

func (obj *ObjectArray) Get(index int) Object {
	obj.Check(index)
	return obj.List[index]
}

func (obj *ObjectArray) GetInt(index int) int {
	return obj.Get(index).(*ObjectInt).Value
}

func (obj *ObjectArray) GetFloat(index int) float64 {
	return obj.Get(index).(*ObjectFloat).Value
}

func (obj *ObjectArray) Check(index int) {
	if obj.List == nil {
		vmError(NULL_POINTER_ERR)
		return
	}

	length := obj.Len()
	if length < 0 || index >= length {
		vmError(INDEX_OUT_OF_BOUNDS_ERR, index, length)
	}
}

func NewObjectArray(size int) *ObjectArray {
	return &ObjectArray{
		List: make([]Object, size),
	}
}

//
// ObjectMap
//
type ObjectMap struct {
	ObjectBase
	// TODO: 临时简单处理, 键为对象hash,值为key,
	Map map[string][2]Object
}

func (obj *ObjectMap) Get(key Object) Object {
	hash := utils.Hash(key)
	v, ok := obj.Map[hash]
	if !ok {
		return nil
	}

	return v[1]
}

func (obj *ObjectMap) Set(key Object, value Object) {
	obj.Map[utils.Hash(key)] = [2]Object{key, value}
}

func (obj *ObjectMap) Delete(key Object) {
	hash := utils.Hash(key)
	delete(obj.Map, hash)
}

func NewObjectMap() *ObjectMap {
	return &ObjectMap{
		Map: make(map[string][2]Object),
	}
}

//
// ObjectInterface
//
type ObjectInterface struct {
	ObjectBase
	Data Object
}

func NewObjectInterface(data Object) *ObjectInterface {
	return &ObjectInterface{
		Data: data,
	}
}

//
// ObjectStruct
//
type ObjectStruct struct {
	ObjectBase
	FieldList []Object
}

func (obj *ObjectStruct) GetField(i int) Object {
	return obj.FieldList[i]
}

func (obj *ObjectStruct) SetField(i int, value Object) {
	obj.FieldList[i] = value
}

func NewObjectStruct(size int) *ObjectStruct {
	return &ObjectStruct{
		FieldList: make([]Object, size),
	}
}

//
// ObjectPointer
//
type ObjectPointer struct {
	ObjectBase
}

//
// ObjectCallInfo 函数返回体 TODO: 临时定义为对象
//
type ObjectCallInfo struct {
	ObjectBase                  // TODO: 兼容
	caller        *GoGoFunction // 调用的函数
	callerAddress int           // 保存执行函数前的pc
	bp            int           // 栈基
}
