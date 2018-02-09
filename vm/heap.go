package vm

const (
	HEAP_THRESHOLD_SIZE = 10240
)
//
// Heap
//
// 虚拟机堆
type Heap struct {
	// TODO:阈值
	currentThreshold int
	objectList       []VmObject
}

func NewHeap() Heap {
	h := Heap{
		currentThreshold: HEAP_THRESHOLD_SIZE,
		objectList: []VmObject{},
	}
	return h
}
