package compiler

import (
	"fmt"
	"os"

	"../vm"
)

var stCompilerList []*Compiler

// Compiler 编译器
type Compiler struct {
	// 词法解析器
	lexer *Lexer

	//
	packageNameList []string

	path string

	requireList []*Require

	// 函数列表
	funcList []*FunctionDefinition

	// vm函数列表
	vmFunctionList []*vm.VmFunction

	// vm类
	vmClassList []*vm.VmClass

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
}

func newCompiler() *Compiler {
	compilerBackup := getCurrentCompiler()
	c := &Compiler{
		statementList:   []Statement{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
	}
	setCurrentCompiler(c)
	// TODO 添加默认函数

	setCurrentCompiler(compilerBackup)
	return c
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

// 开始编译
func (c *Compiler) compile(exeList *vm.ExecutableList, path string, isRequired bool) *vm.Executable {
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

		requireCompiler.addLexerByPath(foundPath)

		// 编译导入的包
		reqExe := requireCompiler.compile(exeList, foundPath, true)
	}

	c.fixTree()

	exe := c.generate()

	exe.Path = path
	exe.IsRequired = isRequired

	exeList.AddExe(exe)

	setCurrentCompiler(compilerBackup)

	return exe
}

func (c *Compiler) Show() {
	fmt.Println("==========")
	fmt.Println("stmt list start\n")
	for _, stmt := range c.statementList {
		stmt.show(0)
	}
	fmt.Println("\nstmt list end")
	fmt.Println("==========\n")
}

// 修正树
func (c *Compiler) fixTree() {
	// fix class
	c.fixClassList()

	for _, func_pos := range c.funcList{
		reserve_function_index(c, func_pos)
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

func (c *Compiler) Generate(exe *vm.Executable) {

	// 添加全局变量声明
	c.addGlobalVariable(exe)
	// 添加函数信息
	c.addFunctions(exe)
	// 添加顶层代码
	c.addTopLevel(exe)
}

func (c *Compiler) generate() *vm.Executable {
	exe := vm.NewExecutable()

	// 添加全局变量声明
	c.addGlobalVariable(exe)
	// 添加函数信息
	c.addFunctions(exe)
	// 添加顶层代码
	c.addTopLevel(exe)

	return exe
}

func (compiler *Compiler) addGlobalVariable(exe *vm.Executable) {
	for _, dl := range compiler.declarationList {
		v := vm.NewVmVariable(dl.name, copyTypeSpecifier(dl.typeSpecifier))

		exe.GlobalVariableList = append(exe.GlobalVariableList, v)
	}
}

// 为每个函数生成所需的信息
func (compiler *Compiler) addFunctions(exe *vm.Executable) {
	for _, srcFd := range compiler.funcList {
		destFd := &vm.VmFunction{}
		copyFunction(srcFd, destFd)

		exe.FunctionList = append(exe.FunctionList, destFd)

		if srcFd.block == nil {
			// 原生函数
			destFd.IsImplemented = false
			continue
		}

		ob := newCodeBuf()
		generateStatementList(exe, srcFd.block, srcFd.block.statementList, ob)

		destFd.IsImplemented = true
		destFd.CodeList = ob.fixOpcodeBuf()
		destFd.LineNumberList = ob.lineNumberList
	}
}

// 生成解释器所需的信息
func (compiler *Compiler) addTopLevel(exe *vm.Executable) {
	ob := newCodeBuf()
	generateStatementList(exe, nil, compiler.statementList, ob)

	exe.CodeList = ob.fixOpcodeBuf()
	exe.LineNumberList = ob.lineNumberList
}

func fixStatementList(currentBlock *Block, statementList []Statement, fd *FunctionDefinition) {
	for _, statement := range statementList {
		statement.fix(currentBlock, fd)
	}
}

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

// ==============================
// utils
// ==============================

func searchDeclaration(name string, currentBlock *Block) *Declaration {

	// 从局部作用域查找
	for b := currentBlock; b != nil; b = b.outerBlock {
		for _, decl := range b.declarationList {
			if decl.name == name {
				return decl
			}
		}
	}

	// 从全局作用域查找
	compiler := getCurrentCompiler()
	for _, decl := range compiler.declarationList {
		if decl.name == name {
			return decl
		}
	}

	return nil
}

func SearchFunction(name string) *FunctionDefinition {
	compiler := getCurrentCompiler()

	for _, pos := range compiler.funcList {
		if pos.name == name {
			return pos
		}
	}
	return nil
}

func compileError(pos Position, errorNumber int, a ...interface{}) {
	fmt.Println("编译错误")
	fmt.Printf("Line: %d:%d\n", pos.Line, pos.Column)
	fmt.Printf(errMessageMap[errorNumber], a...)
	fmt.Println("\n")
	os.Exit(1)
}

// 设置全局compiler
var stCurrentCompiler *Compiler

func getCurrentCompiler() *Compiler {
	return stCurrentCompiler
}

func setCurrentCompiler(c *Compiler) {
	stCurrentCompiler = c
}

// ==============================
// 词法解析，语法解析
// ==============================

func CompileFile(path string) *vm.ExecutableList {
	// 输出yacc错误信息
	yyErrorVerbose = true

	compiler := newCompiler()

	compiler.addLexerByPath(path)

	exeList := &vm.ExecutableList{}

	exe := compiler.compile(exeList, "", false)
	exeList.TopLevel = exe

	return exeList
}

func searchCompiler(list []*Compiler, package_name []string) *Compiler {
	for _, c := range list {
		if comparePackageName(c.packageNameList, package_name) {
			return c
		}
	}
	return nil
}

func comparePackageName(packageNameList1, packageNameList2 []string) bool {
	if packageNameList1 == nil {
		if packageNameList2 == nil {
			return true
		}
		return false
	}

	length1 := len(packageNameList1)
	length2 := len(packageNameList2)

	if length1 != length2 {
		return false
	}

	for i := 0; i < length1; i++ {
		if packageNameList1[i] != packageNameList2[i] {
			return false
		}
	}

	return true
}

func (c *Compiler) fixClassList() {

	for _, classD := range c.classDefinitionList {
		classD.add_class()
		fix_extends(classD)
	}

	//for _, classD := range c.classDefinitionList {
	//    add_super_interfaces(classD)
	//}

	for _, classD := range c.classDefinitionList {
		c.currentClassDefinition = classD
		// TODO 添加默认构造函数
		//add_default_constructor(classD)
		c.currentClassDefinition = nil
	}

	for _, classD := range c.classDefinitionList {

		c.currentClassDefinition = classD

		field_index, method_index := get_super_field_method_count(classD)

		for member_pos := range classD.memberList {
			switch member := member_pos.(type) {
			case *MethodMember:
				member.functionDefinition.fix()

				superMember := search_member_in_super(classD, member.functionDefinition.name)

				if superMember {
					superMethodMember, ok := superMember.(*MethodMember)
					if !ok {
						compileError(member.Position(), FIELD_OVERRIDED_ERR, member.functionDefinition.name)
					}
					member.methodIndex = superMethodMember.methodIndex
				} else {
					member.methodIndex = method_index
					method_index++
				}
			case *FieldMember:
				member.typeSpecifier.fix()

				superMember := search_member_in_super(classD, member.name)

				if superMember {
					compileError(member_pos.Position(), FIELD_NAME_DUPLICATE_ERR, member.name)
				} else {
					member.fieldIndex = field_index
					field_index++
				}
			default:
				panic("TODO")
			}
		}
		c.currentClassDefinition = nil
	}
}

func get_super_field_method_count(ClassDefinition *cd) (int, int) {
	fieldIndex := 0
	methodIndex := 0

	for supderClass := cd.superClass; supderClass != nil; supderClass = supderClass.superClass {
		for member_pos := range supderClass.memberList {
			switch member := member_pos.(type) {
			case *MethodMember:
				if member.methodIndex > methodIndex {
					methodIndex = member.methodIndex
				}
			case *FieldMember:
				if member.fieldIndex > fieldIndex {
					fieldIndex = member.fieldIndex
				}
			default:
				panic("TODO")
			}
		}
	}
	return fieldIndex, methodIndex
}

func search_member_in_super(class_def *ClassDefinition, member_name string) MemberDeclaration {
	var member MemberDeclaration

	if class_def.superClass {
		member := search_member(class_def.superClass, member_name)
		if member != nil {
			return member
		}
	}

	return nil
}

func search_member(cd *ClassDefinition, memberName string) MemberDeclaration {

	for _, md := range cd.memberList {
		switch member := md.(type) {
		case *MethodMember:
			if member.functionDefinition.name == memberName {
				return member
			}
		case *FieldMember:
			if member.name == memberName {
				return member
			}
		default:
			panic("TODO")
		}
	}

	// 递归查找
	if cd.superClass {
		member = search_member(cd.superClass, memberName)
		if member != nil {
			return member
		}
	}

	return nil
}

func search_class_and_add(pos Position, name string, class_index_p *int) *ClassDefinition {

	cd := search_class(name)

	if cd == nil {
		compileError(pos, CLASS_NOT_FOUND_ERR, name)
	}

	*class_index_p = cd.add_class()

	return cd
}

func fix_extends(cd *ClassDefinition) {
	var dummy_class_index int

	for _, extend := range cd.extends {
		super = search_class_and_add(cd.Position(), extend.identifier, &dummy_class_index)

		extend.classDefinition = super

		if cd.superClass {
			compileError(cd.Position(), MULTIPLE_INHERITANCE_ERR, super.name)
		}

		// TODO 只有接口才可以继承
		if !super.is_abstract {
			compileError(cd.Position(), INHERIT_CONCRETE_CLASS_ERR, super.name)
		}

		cd.superClass = super
	}
}

func add_default_constructor(cd *ClassDefinition) {

	for _, member_pos := range cd.memberList {
		methodMember, ok := member_pos.(*MethodMember)
		if !ok {
			continue
		}

		if methodMember.functionDefinition.name == "init" {
			return
		}
	}

	// TODO 行数不对
	typ := &TypeSpecifier{basicType: vm.VoidType}
	block := &Block{}

	if cd.superClass != nil {
		statement := &ExpressionStatement{}

		// TODO super
		// 是否直接使用父类的初始化方法
		member_e = createMemberExpression(nil, DEFAULT_CONSTRUCTOR_NAME, fd.Position())
		func_call_e = createFunctionCallExpression(member_e, nil)
		statement.expression = func_call_e
		block.statementList = []Statement{statement}
	} else {
		block.statementList = nil
	}
	fd := createFunctionDefinition(typ, DEFAULT_CONSTRUCTOR_NAME, nil, block)

	if cd.memberList == nil {
		cd.memberList = []MemberDeclaration{}
	}

	cd.memberList = append(cd.memberList, createMethodMember(fd, cd.Position()))
}

func reserve_function_index(compiler *Compiler,src *FunctionDefinition) int {

    src_package_name = src.getPackageName()

	for i, vmFunction := range compiler.vmFunctionList{
        if (src_package_name == vmFunction.PackageName) && (src.classDefinitio.name == vmFunction.Name) {
            return i
        }
    }

	dest := &vm.VmFunction{}
	compiler.vmFunctionList = append(compiler.vmFunctionList, dest)

    dest.PackageName = src_package_name
    if src.classDefinition {
        dest.name = createMethodFunctionName(src.classDefinition.name, src.name)
    } else {
        dest.name = src.name
    }

    return len(compiler.vmFunctionList) - 1
}
