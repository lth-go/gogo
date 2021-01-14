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

type _ObjectArray struct {
	ObjectBase
	// Length int
	// ValueType int
	List []Object
}

func (obj *_ObjectArray) Mark() {
	obj.ObjectBase.Mark()

	for _, subObj := range obj.List {
		subObj.Mark()
	}
}

func (obj *_ObjectArray) ResetMark() {
	obj.ObjectBase.ResetMark()

	for _, subObj := range obj.List {
		subObj.ResetMark()
	}
}

func (obj *_ObjectArray) Sweep() {
	for _, subObj := range obj.List {
		subObj.Sweep()
	}
}

func (obj *_ObjectArray) Len() int {
	return len(obj.List)
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
