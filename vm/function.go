package vm

// 虚拟机全局函数
type ExecFunction interface {
	getName() string
	getPackageName() string
}

// 原生函数
type NativeFunction struct {
	Name        string
	PackageName string

	proc     NativeFunctionProc
	argCount int
}

func (f *NativeFunction) getName() string        { return f.Name }
func (f *NativeFunction) getPackageName() string { return f.PackageName }

type NativeFunctionProc func(vm *VirtualMachine, argCount int, args []Object) Object

// 保存调用函数的索引
type GFunction struct {
	Name        string
	PackageName string

	Executable *Executable
	Index      int
}

func (f *GFunction) getName() string        { return f.Name }
func (f *GFunction) getPackageName() string { return f.PackageName }
