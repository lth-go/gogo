package vm

type ExecClass struct {
	vmClass *VmClass

	Executable *ExecutableEntry

	packageName string
	name string

	// TODO remove
	isImplemented bool

	classIndex int

	superClass *ExecClass
	classTable *VmVTable

	fieldTypeList []*VmTypeSpecifier
}

type VmVTable struct  {
	execClass *ExecClass
	table []*VTableItem
}

type VTableItem struct {
	name string
	index int
} 
