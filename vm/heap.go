package vm

const (
	heapThresholdSize = 10240
)

// 虚拟机堆
type Heap struct {
	list             []Object
	currentThreshold int
}

func NewHeap() *Heap {
	h := &Heap{
		list:             make([]Object, 0),
		currentThreshold: heapThresholdSize,
	}
	return h
}

func (h *Heap) Append(value Object) {
	h.list = append(h.list, value)
	h.currentThreshold += value.Len()
}
