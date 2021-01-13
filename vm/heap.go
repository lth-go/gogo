package vm

const (
	heapThresholdSize = 10240
)

// 虚拟机堆
type Heap struct {
	// TODO:阈值
	currentThreshold int
	objectList       []Object
}

func NewHeap() *Heap {
	h := &Heap{
		currentThreshold: heapThresholdSize,
		objectList:       []Object{},
	}
	return h
}

func (h *Heap) Append(value Object) {
	h.objectList = append(h.objectList, value)
	h.currentThreshold += value.Len()
}
