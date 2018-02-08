package vm

//
// Function
//
// 虚拟机全局函数
type Function interface {
	getName() string
}

// 原生函数
type NativeFunction struct {
	Name string

	proc     VmNativeFunctionProc
	argCount int
}

func (f *NativeFunction) getName() string { return f.Name }

type VmNativeFunctionProc func(vm *VmVirtualMachine, argCount int, args VmValue) VmValue

type GFunction struct {
	Name string

	Executable *Executable
	Index      int
}

func (f *GFunction) getName() string { return f.Name }
