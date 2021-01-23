package vm

const (
	heapThresholdSize = 10240
)

// 虚拟机堆
type Heap struct {
	currentThreshold int
	objectList       []Object
}

func NewHeap() *Heap {
	h := &Heap{
		currentThreshold: heapThresholdSize,
		objectList:       make([]Object, 0),
	}
	return h
}

func (h *Heap) Append(value Object) {
	h.objectList = append(h.objectList, value)
	h.currentThreshold += value.Len()
}
