package vm

//
// Function
//
// 虚拟机全局函数
type Function interface {
	getName() string
	getPackageName() string
}

// TODO
type FunctionImpl struct {
	Name string
	PackageName string
}

// 原生函数
type NativeFunction struct {
	Name string
	PackageName string

	proc     VmNativeFunctionProc
	argCount int
}

func (f *NativeFunction) getName() string { return f.Name }
func (f *NativeFunction) getPackageName() string { return f.PackageName }

type VmNativeFunctionProc func(vm *VmVirtualMachine, argCount int, args []VmValue) VmValue

// 保存调用函数的索引
type GFunction struct {
	Name string
	PackageName string

	Executable *ExecutableEntry
	Index      int
}

func (f *GFunction) getName() string { return f.Name }
func (f *GFunction) getPackageName() string { return f.PackageName }
