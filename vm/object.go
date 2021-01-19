package vm

// 虚拟机对象接口
type Object interface {
	isMarked() bool // 是否设置标记位
	Mark()          // 设置标记位
	ResetMark()     // 重置标记位
	Sweep()         // 垃圾回收
	Len() int       // 计算堆阈值
}

type ObjectBase struct {
	marked bool
}

func (obj *ObjectBase) isMarked() bool {
	return obj.marked
}

func (obj *ObjectBase) Mark() {
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

type ObjectInt struct {
	ObjectBase
	Value int
}

type ObjectFloat struct {
	ObjectBase
	Value float64
}

type ObjectString struct {
	ObjectBase
	Value string
}

type ObjectNil struct {
	ObjectBase
}

type ObjectArray struct {
	ObjectBase
	// Length int
	// ValueType int
	List []Object
}

func (obj *ObjectArray) Mark() {
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

type ObjectMap struct {
	ObjectBase
	KeyType   int
	ValueType int
	KeyList   []Object
	ValueList []Object
}

type ObjectStruct struct {
	ObjectBase
	FieldList []Object
}

type ObjectPointer struct {
	ObjectBase
}

func NewObjectInt(value int) *ObjectInt {
	return &ObjectInt{
		Value: value,
	}
}

func NewObjectFloat(value float64) *ObjectFloat {
	return &ObjectFloat{
		Value: value,
	}
}

func NewObjectString(value string) *ObjectString {
	return &ObjectString{
		Value: value,
	}
}

//
// ObjectCallInfo 函数返回体
// TODO: 临时定义为对象
//
type ObjectCallInfo struct {
	ObjectBase                  // TODO: 兼容
	caller        *GoGoFunction // 调用的函数
	callerAddress int           // 保存执行函数前的pc
	base          int
}
