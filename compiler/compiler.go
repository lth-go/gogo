package compiler

import (
	"fmt"
	"strings"

	"../vm"
)

// 全局compiler列表
var stCompilerList []*Compiler

// 设置全局compiler
var stCurrentCompiler *Compiler

func getCurrentCompiler() *Compiler { return stCurrentCompiler }

func setCurrentCompiler(c *Compiler) { stCurrentCompiler = c }

// Compiler 编译器
type Compiler struct {
	// 词法解析器
	lexer *Lexer

	// 包名
	packageNameList []string
	// 源文件路径
	path string
	// 依赖的包
	requireList []*Require

	// 函数列表
	funcList []*FunctionDefinition
	// 声明列表
	declarationList []*Declaration
	// 语句列表
	statementList []Statement
	// 类定义列表
	classDefinitionList []*ClassDefinition

	// 当前块
	currentBlock *Block
	// 当前类
	currentClassDefinition *ClassDefinition

	//current_try_statement *TryStatement
	//current_catch_clause *CatchClause

	// 已加载compiler列表
	requiredList []*Compiler

	// arrayMethodList  []*FunctionDefinition
	// stringMethodList []*FunctionDefinition

	// TODO 能否去掉
	// vm函数列表
	vmFunctionList []*vm.VmFunction
	// vm类
	vmClassList []*vm.VmClass
}

func newCompiler() *Compiler {
	compilerBackup := getCurrentCompiler()
	c := &Compiler{
		requireList:         []*Require{},
		funcList:            []*FunctionDefinition{},
		vmFunctionList:      []*vm.VmFunction{},
		vmClassList:         []*vm.VmClass{},
		declarationList:     []*Declaration{},
		statementList:       []Statement{},
		classDefinitionList: []*ClassDefinition{},
		requiredList:        []*Compiler{},
	}
	setCurrentCompiler(c)
	// TODO 添加默认函数

	setCurrentCompiler(compilerBackup)
	return c
}
func (c *Compiler) getPackageName() string {
	return strings.Join(c.packageNameList, ".")
}
func (c *Compiler) addLexer(lexer *Lexer) {
	lexer.compiler = c
	c.lexer = lexer
}

func (c *Compiler) addLexerByPath(path string) {
	lexer := newLexerByFilePath(path)
	c.addLexer(lexer)

	c.path = path
}

//////////////////////////////
// 函数定义
//////////////////////////////
func (c *Compiler) functionDefine(typ *TypeSpecifier, identifier string, parameterList []*Parameter, block *Block) {
	// 定义重复
	if SearchFunction(identifier) != nil || searchDeclaration(identifier, nil) != nil {
		compileError(typ.Position(), FUNCTION_MULTIPLE_DEFINE_ERR, identifier)
	}

	fd := &FunctionDefinition{
		typeSpecifier:     typ,
		name:              identifier,
		parameterList:     parameterList,
		block:             block,
		index:             len(c.funcList),
		localVariableList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

//////////////////////////////
// 编译
//////////////////////////////
func (c *Compiler) compile(exeList *vm.ExecutableList, isRequired bool) *vm.Executable {
	compilerBackup := getCurrentCompiler()
	setCurrentCompiler(c)

	// 开始解析文件
	if yyParse(c.lexer) != 0 {
		panic(c.lexer.e)
	}

	for _, require := range c.requireList {
		// 判断是否已经被解析过
		requireCompiler := searchCompiler(stCompilerList, require.packageNameList)
		if requireCompiler != nil {
			c.requiredList = append(c.requiredList, requireCompiler)
			continue
		}

		requireCompiler = newCompiler()

		requireCompiler.packageNameList = require.packageNameList

		c.requiredList = append(c.requiredList, requireCompiler)
		stCompilerList = append(stCompilerList, requireCompiler)

		// 获取要导入的全路径
		foundPath := require.getFullPath()

		// 编译导入的包
		requireCompiler.addLexerByPath(foundPath)
		requireCompiler.compile(exeList, true)
	}

	// fix and generate
	c.fixTree()
	exe := c.generate()

	exe.Path = c.path
	exe.IsRequired = isRequired

	exeList.AddExe(exe)

	setCurrentCompiler(compilerBackup)

	return exe
}

//////////////////////////////
// 打印语法树
//////////////////////////////
func (c *Compiler) Show() {
	fmt.Println("==========")
	fmt.Println("stmt list start\n")
	for _, stmt := range c.statementList {
		stmt.show(0)
	}
	fmt.Println("\nstmt list end")
	fmt.Println("==========\n")
}

//////////////////////////////
// 修正树
//////////////////////////////
func (c *Compiler) fixTree() {
	// class
	c.fixClassList()

	// TODO remove
	// add function
	for _, fd := range c.funcList {
		c.addToVmFunctionList(fd)
	}

	// 修正表达式列表
	fixStatementList(nil, c.statementList, nil)

	// 修正函数
	for _, fd := range c.funcList {
		fd.fix()
	}

	// 修正全局声明
	for varCount, decl := range c.declarationList {
		decl.variableIndex = varCount
	}

}

func (c *Compiler) fixClassList() {

	// 修正继承
	for _, cd := range c.classDefinitionList {
		cd.addToCurrentCompiler()
		cd.fixExtends()
	}

	// 修正方法和属性, 搜索父类, 设置索引
	for _, cd := range c.classDefinitionList {

		c.currentClassDefinition = cd

		fieldIndex, methodIndex := cd.getSuperFieldMethodCount()

		for memberIfs := range cd.memberList {
			switch member := memberIfs.(type) {
			case *MethodMember:
				member.functionDefinition.fix()

				superMember := cd.searchMemberInSuper(member.functionDefinition.name)

				if superMember != nil {
					superMethodMember, ok := superMember.(*MethodMember)
					if !ok {
						compileError(member.Position(), FIELD_OVERRIDED_ERR, member.functionDefinition.name)
					}
					member.methodIndex = superMethodMember.methodIndex
				} else {
					member.methodIndex = methodIndex
					methodIndex++
				}
			case *FieldMember:
				member.typeSpecifier.fix()

				superMember := cd.searchMemberInSuper(member.name)

				// TODO 只有方法能够继承?
				if superMember != nil {
					compileError(member.Position(), FIELD_NAME_DUPLICATE_ERR, member.name)
				} else {
					member.fieldIndex = fieldIndex
					fieldIndex++
				}
			default:
				panic("TODO")
			}
		}
		c.currentClassDefinition = nil
	}
}

// 添加VmFuntion
func (c *Compiler) addToVmFunctionList(src *FunctionDefinition) int {

	srcPackageName := src.getPackageName()

	for i, vmFunction := range c.vmFunctionList {
		if (srcPackageName == vmFunction.PackageName) && (src.name == vmFunction.Name) {
			return i
		}
	}

	dest := &vm.VmFunction{}
	c.vmFunctionList = append(c.vmFunctionList, dest)

	dest.PackageName = srcPackageName
	if src.classDefinition != nil {
		dest.Name = createMethodFunctionName(src.classDefinition.name, src.name)
	} else {
		dest.Name = src.name
	}

	return len(c.vmFunctionList) - 1
}

//////////////////////////////
// 生成字节码
//////////////////////////////
func (c *Compiler) generate() *vm.Executable {
	exe := vm.NewExecutable()
	exe.PackageName = c.getPackageName()

	exe.FunctionList = c.vmFunctionList
	exe.ClassDefinitionList = c.vmClassList

	// 添加全局变量声明
	c.addGlobalVariable(exe)
	// 添加类信息
	c.addClasses(exe)
	// 添加函数信息
	c.addFunctions(exe)
	// 添加顶层代码
	c.addTopLevel(exe)

	return exe
}

// 添加全局变量
func (compiler *Compiler) addGlobalVariable(exe *vm.Executable) {
	for _, dl := range compiler.declarationList {

		newValue := vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))
		exe.GlobalVariableList = append(exe.GlobalVariableList, newValue)
	}
}

// 添加类
func (compiler *Compiler) addClasses(exe *vm.Executable) {
	for _, cd := range compiler.classDefinitionList {
		vmClass := compiler.searchVmClass(cd)
		vmClass.IsImplemented = true
	}

	for _, vmClass := range compiler.vmClassList {
		cd := searchClass(vmClass.Name)
		addClass(cd, vmClass)
	}
}

// TODO 改名
func addClass(cd *ClassDefinition, dest *vm.VmClass) {

	if cd.superClass != nil {
		dest.SuperClass = &vm.VmClassIdentifier{}
		dest.Name = cd.name
		dest.PackageName = cd.getPackageName()
	} else {
		dest.SuperClass = nil
	}

	for _, memberIfs := range cd.memberList {
		switch member := memberIfs.(type) {
		case *MethodMember:
			newMethod := &vm.VmMethod{
				Name: member.functionDefinition.name,
			}
			dest.MethodList = append(dest.MethodList, newMethod)
		case *FieldMember:
			newField := &vm.VmField{
				Name: member.name,
				Typ:  copyTypeSpecifier(member.typeSpecifier),
			}
			dest.FieldList = append(dest.FieldList, newField)
		default:
			panic("TODO")
		}
	}
}

// 添加函数
func (c *Compiler) addFunctions(exe *vm.Executable) {

	in_this_exe := make([]bool, len(c.vmFunctionList))

	for _, fd := range c.funcList {
		// TODO 为什么block是空
		if fd.classDefinition != nil && fd.block == nil {
			continue
		}

		dest_idx := c.getFunctionIndex(fd)
		in_this_exe[dest_idx] = true

		add_function(exe, fd, c.vmFunctionList[dest_idx], true)
	}

	for i, vmFunc := range c.vmFunctionList {
		if in_this_exe[i] {
			continue
		}

		fd := SearchFunction(vmFunc.Name)
		add_function(exe, fd, vmFunc, false)
	}
}

func add_function(exe *vm.Executable, src *FunctionDefinition, dest *vm.VmFunction, in_this_exe bool) {
	ob := newCodeBuf()

	dest.TypeSpecifier = copyTypeSpecifier(src.typeS())
	dest.ParameterList = copyParameterList(src.parameterList)

	if src.block != nil && in_this_exe {
		generateStatementList(exe, src.block, src.block.statementList, ob)

		dest.IsImplemented = true
		dest.CodeList = ob.fixOpcodeBuf()
		dest.LineNumberList = ob.lineNumberList
		dest.LocalVariableList = copyLocalVariables(src)
	} else {
		dest.IsImplemented = false
		dest.LocalVariableList = nil
	}

	if src.classDefinition != nil {
		dest.IsMethod = true
	} else {
		dest.IsMethod = false
	}
}

func (c *Compiler) getFunctionIndex(src *FunctionDefinition) int {
	var func_name string

	srcPackageName := src.getPackageName()

	if src.classDefinition != nil {
		func_name = createMethodFunctionName(src.classDefinition.name, src.name)
	} else {
		func_name = src.name
	}

	for i, vmFunc := range c.vmFunctionList {
		if srcPackageName == vmFunc.PackageName && func_name == vmFunc.Name {
			return i
		}
	}

	panic("TODO")
}

// 添加字节码
func (compiler *Compiler) addTopLevel(exe *vm.Executable) {
	ob := newCodeBuf()
	generateStatementList(exe, nil, compiler.statementList, ob)

	exe.CodeList = ob.fixOpcodeBuf()
	exe.LineNumberList = ob.lineNumberList
}

// other
func (c *Compiler) searchVmClass(src *ClassDefinition) *vm.VmClass {

	srcPackageName := src.getPackageName()

	for _, vmClass := range c.vmClassList {
		if srcPackageName == vmClass.PackageName && src.name == vmClass.Name {
			return vmClass
		}
	}
	panic("TODO")
}

// ==============================
// 编译文件
// ==============================

func CompileFile(path string) *vm.ExecutableList {
	// 输出yacc错误信息
	yyErrorVerbose = true

	compiler := newCompiler()

	compiler.addLexerByPath(path)

	exeList := vm.NewExecutableList()

	exe := compiler.compile(exeList, false)
	exeList.TopLevel = exe

	return exeList
}
