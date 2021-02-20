package vm

//
// 字节码解释器
//
type Executable struct {
	PackageName  string        // 包名
	Constant     []interface{} // 常量池
	VariableList []*Variable   // 全局变量
	FunctionList []*Function   // 函数列表
}

func NewExecutable() *Executable {
	return &Executable{}
}

//
// Variable 全局变量
//
type Variable struct {
	PackageName string
	Name        string
	Type        *Type
	Value       interface{}
}

func (v *Variable) Init() {
	if v.Value != nil {
		return
	}

	var value interface{}

	if v.Type.IsReferenceType() {
		value = NilObject
	} else {
		switch v.Type.GetBasicType() {
		case BasicTypeBool, BasicTypeInt:
			value = 0
		case BasicTypeFloat:
			value = 0.0
		case BasicTypeString:
			value = ""
		case BasicTypeNil:
			fallthrough
		default:
			panic("TODO")
		}
	}

	v.Value = value
}

func NewVmVariable(packageName string, name string, typ *Type) *Variable {
	return &Variable{
		PackageName: packageName,
		Name:        name,
		Type:        typ,
	}
}

//
// Function 函数
//
type Function struct {
	IsImplemented  bool          // 是否在当前包实现
	PackageName    string        // 包名
	Name           string        // 函数名
	ArgCount       int           // 参数数量
	ResultCount    int           // 返回值数量
	VariableList   []Object      // 局部变量列表
	CodeList       []byte        // 字节码类表
	LineNumberList []*LineNumber // 行号对应表
}

//
// 行号对应表
//
type LineNumber struct {
	// 源代码行号
	LineNumber int
	// 字节码开始的位置
	StartPc int
	// 接下来有多少字节码对应相同的行号
	PcCount int
}
