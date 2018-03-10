package vm

// 虚拟机中的类信息
type ExecClass struct {
	// 具体类实现
	vmClass *Class

	// 所属的exe
	Executable *ExecutableEntry

	// 所属包名
	packageName string
	// 类名
	name        string

	// 虚拟机中的类下标
	classIndex int

	// 父类
	superClass *ExecClass
	// 虚表
	classTable *VTable

	// 所有字段类型表
	fieldTypeList []*TypeSpecifier
}

// 虚表
type VTable struct {
	execClass *ExecClass
	table     []*VTableItem
}

// 虚表字段
type VTableItem struct {
	name  string
	index int
}
