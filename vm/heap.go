package vm

//
// Heap
//
// 虚拟机堆
type Heap struct {
	// TODO:阈值
	currentThreshold int
	objectList       []VmObject
}
