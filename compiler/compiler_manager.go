package compiler

import (
	"github.com/lth-go/gogo/vm"
)

type CompilerManager struct {
	doingList []*Compiler
	doneList  []*Compiler
}

var compilerManager = &CompilerManager{
	doingList: []*Compiler{},
	doneList:  []*Compiler{},
}

// GetCurrentCompiler
func GetCurrentCompiler() *Compiler {
	length := len(compilerManager.doingList)
	if length == 0 {
		return nil
	}

	return compilerManager.doingList[length-1]
}

func PushCurrentCompiler(c *Compiler) {
	compilerManager.doingList = append(compilerManager.doingList, c)
}

func PopCurrentCompiler() {
	compilerManager.doingList = compilerManager.doingList[:len(compilerManager.doingList)-1]
}

func IsCompiling(packageName string) bool {
	for _, c := range compilerManager.doingList {
		if c.GetPackageName() == packageName {
			return true
		}
	}

	return false
}

func AddDoneCompilerList(c *Compiler) {
	compilerManager.doneList = append(compilerManager.doneList, c)
}

func GetDoneCompilerList() []*Compiler {
	return compilerManager.doneList
}

func GetDoneCompiler(packageName string) *Compiler {
	for _, c := range GetDoneCompilerList() {
		if c.GetPackageName() == packageName {
			return c
		}
	}
	return nil
}

func (cm *CompilerManager) Parse(path string) {
	c := NewCompiler(path)

	PushCurrentCompiler(c)
	AddDoneCompilerList(c)

	// 生成语法树
	c.Parse()

	for _, imp := range c.importList {
		if IsCompiling(imp.packageName) {
			panic("TODO")
		}

		// 判断是否已经被解析过
		if GetDoneCompiler(imp.packageName) != nil {
			continue
		}

		cm.Parse(imp.GetPath())
	}
	PopCurrentCompiler()
}

func (cm *CompilerManager) Compile() []*vm.Executable {
	doneCompilerList := GetDoneCompilerList()

	// 倒序编译,防止依赖问题
	for i := len(doneCompilerList) - 1; i >= 0; i-- {
		c := doneCompilerList[i]

		PushCurrentCompiler(c)
		c.FixTree()
		PopCurrentCompiler()
	}

	exeList := make([]*vm.Executable, 0)
	for i := len(doneCompilerList) - 1; i >= 0; i-- {
		c := doneCompilerList[i]

		exe := c.Generate()
		exeList = append(exeList, exe)
	}

	return exeList
}

func Parse(path string) {
	compilerManager.Parse(path)
}

func Compile() []*vm.Executable {
	return compilerManager.Compile()
}

//
// 编译文件
//
func CompileFile(path string) *vm.ExecutableList {
	// 输出yacc错误信息
	if true {
		yyErrorVerbose = true
	}

	Parse(path)
	exeList := Compile()

	return vm.NewExecutableList(exeList)
}
