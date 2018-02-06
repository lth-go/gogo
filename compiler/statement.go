package compiler

import (
	"../vm"
)

// ==============================
// 衍生类型
// ==============================

type TypeDerive interface {
}

type FunctionDerive struct {
	parameterList []*Parameter
}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl
	// 基本类型
	basicType vm.BasicType
	// 派生类型
	deriveList []TypeDerive
}

// ==============================
// Statement 接口
// ==============================

// Statement 语句接口
type Statement interface {
	// Pos接口
	Pos

	fix(*Block, *FunctionDefinition)
	generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf)

	show(ident int)
}

type StatementImpl struct {
	PosImpl
}

//
// Block ...
//
type Block struct {
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration

	// 块信息，函数块，还是条件语句
	parent BlockInfo
}

func (b *Block) show(ident int) {
	printWithIdent("Block", ident)
	subIdent := ident + 2

	for _, stmt := range b.statementList {
		stmt.show(subIdent)
	}

	for _, decl := range b.declarationList {
		decl.show(subIdent)
	}
}

func addDeclaration(b *Block, decl *Declaration, fd *FunctionDefinition, pos Position) {
	if searchDeclaration(decl.name, b) != nil {
		compileError(pos, VARIABLE_MULTIPLE_DEFINE_ERR, "Declaration name: %s\n", decl.name)
	}

	if b != nil {
		b.declarationList = append(b.declarationList, decl)
		fd.addLocalVariable(decl)
		decl.isLocal = true
	} else {
		compiler := getCurrentCompiler()
		compiler.declarationList = append(compiler.declarationList, decl)
		decl.isLocal = false
	}
}

type BlockInfo interface{}

type StatementBlockInfo struct {
	statement     Statement
	continueLabel int
	breakLabel    int
}

type FunctionBlockInfo struct {
	function *FunctionDefinition
	endLabel int
}

//
// FunctionDefinition 函数定义
//
type FunctionDefinition struct {
	typeSpecifier *TypeSpecifier
	name          string
	parameterList []*Parameter
	block         *Block

	index int

	localVariableList []*Declaration
}

// TODO 是否可以去掉
func (fd *FunctionDefinition) typeS() *TypeSpecifier {
	return fd.typeSpecifier
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {

	for _, param := range fd.parameterList {
		if searchDeclaration(param.name, fd.block) != nil {
			compileError(param.Position(), PARAMETER_MULTIPLE_DEFINE_ERR, "parameter name: %s\n", param.name)
		}
		decl := &Declaration{name: param.name, typeSpecifier: param.typeSpecifier}

		addDeclaration(fd.block, decl, fd, param.Position())
	}
}

func (fd *FunctionDefinition) addReturnFunction() {

	if fd.block.statementList == nil {
		ret := &ReturnStatement{returnValue: nil}
		ret.fix(fd.block, fd)
		fd.block.statementList = []Statement{ret}
		return
	}

	last := fd.block.statementList[len(fd.block.statementList)-1]
	_, ok := last.(*ReturnStatement)
	if ok {
		return
	}

	ret := &ReturnStatement{returnValue: nil}
	ret.fix(fd.block, fd)
	fd.block.statementList = append(fd.block.statementList, ret)
}

func (fd *FunctionDefinition) addLocalVariable(decl *Declaration) {
	decl.variableIndex = len(fd.localVariableList)
	fd.localVariableList = append(fd.localVariableList, decl)
}

func (fd *FunctionDefinition) checkArgument(currentBlock *Block, expr Expression) {
	functionCallExpr := expr.(*FunctionCallExpression)

	ParameterList := fd.parameterList
	argumentList := functionCallExpr.argumentList

	length := len(ParameterList)
	if len(argumentList) != length {
		compileError(expr.Position(), ARGUMENT_COUNT_MISMATCH_ERR, "Need: %d, Give: %d\n", length, len(argumentList))
	}

	for i := 0; i < length; i++ {
		argumentList[i] = argumentList[i].fix(currentBlock)
		argumentList[i] = createAssignCast(argumentList[i], ParameterList[i].typeSpecifier)
	}
}

//
// Parameter 形参
//
type Parameter struct {
	PosImpl
	typeSpecifier *TypeSpecifier
	name          string
}

// ==============================
// ExpressionStatement
// ==============================

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	StatementImpl
	expression Expression
}

func (stmt *ExpressionStatement) show(ident int) {
	printWithIdent("ExprStmt", ident)

	subIdent := ident + 2

	stmt.expression.show(subIdent)
}

func (stmt *ExpressionStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.expression.fix(currentBlock)
}

func (stmt *ExpressionStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr := stmt.expression
	switch assignExpr := expr.(type) {
	case *AssignExpression:
		// TODO
		assignExpr.generateEx(exe, currentBlock, ob)
	default:
		expr.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), vm.VM_POP)
	}
}

// ==============================
// IfStatement
// ==============================

// IfStatement if表达式
type IfStatement struct {
	StatementImpl

	condition Expression
	thenBlock *Block
	elifList  []*Elif
	elseBlock *Block
}

func (stmt *IfStatement) show(ident int) {
	printWithIdent("IfStmt", ident)

	subIdent := ident + 2
	stmt.condition.show(subIdent)
	if stmt.thenBlock != nil {
		stmt.thenBlock.show(subIdent)
	}
	for _, elif := range stmt.elifList {
		printWithIdent("Elif", subIdent)
		elif.condition.show(subIdent + 2)
		elif.block.show(subIdent + 2)
	}

	if stmt.elseBlock != nil {
		stmt.elseBlock.show(subIdent)
	}
}

func (stmt *IfStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

	stmt.condition.fix(currentBlock)

	if stmt.thenBlock != nil {
		fixStatementList(stmt.thenBlock, stmt.thenBlock.statementList, fd)
	}

	for _, elif := range stmt.elifList {
		elif.condition.fix(currentBlock)

		if elif.block != nil {
			fixStatementList(elif.block, elif.block.statementList, fd)
		}
	}

	if stmt.elseBlock != nil {
		fixStatementList(stmt.elseBlock, stmt.elseBlock.statementList, fd)
	}
}

func (stmt *IfStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	stmt.condition.generate(exe, currentBlock, ob)

	// 获取false跳转地址
	ifFalseLabel := getLabel(ob)
	generateCode(ob, stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

	generateStatementList(exe, stmt.thenBlock, stmt.thenBlock.statementList, ob)

	// 获取结束跳转地址
	endLabel := getLabel(ob)

	// 直接跳到最后
	generateCode(ob, stmt.Position(), vm.VM_JUMP, endLabel)

	// 设置false跳转地址,如果false,直接执行这里
	setLabel(ob, ifFalseLabel)

	for _, elif := range stmt.elifList {
		elif.condition.generate(exe, currentBlock, ob)

		// 获取false跳转地址
		ifFalseLabel = getLabel(ob)
		generateCode(ob, stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(exe, elif.block, elif.block.statementList, ob)

		// 直接跳到最后
		generateCode(ob, stmt.Position(), vm.VM_JUMP, endLabel)

		// 设置false跳转地址,如果false,直接执行这里
		setLabel(ob, ifFalseLabel)
	}
	if stmt.elseBlock != nil {
		generateStatementList(exe, stmt.elseBlock, stmt.elseBlock.statementList, ob)
	}

	// 设置结束地址
	setLabel(ob, endLabel)
}

// Elif ...
type Elif struct {
	condition Expression
	block     *Block
}

// ==============================
// ForStatement
// ==============================

// ForStatement for语句
type ForStatement struct {
	StatementImpl

	init      Expression
	condition Expression
	post      Expression
	block     *Block
}

func (stmt *ForStatement) show(ident int) {
	printWithIdent("ForStmt", ident)
	subIdent := ident + 2

	stmt.init.show(subIdent)
	stmt.condition.show(subIdent)
	stmt.post.show(subIdent)

	if stmt.block != nil {
		stmt.block.show(subIdent)
	}
}

func (stmt *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.init.fix(currentBlock)
	stmt.condition.fix(currentBlock)
	stmt.post.fix(currentBlock)

	if stmt.block != nil {
		fixStatementList(stmt.block, stmt.block.statementList, fd)
	}
}
func (stmt *ForStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	if stmt.init != nil {
		stmt.generate(exe, currentBlock, ob)
	}

	// 获取循环地址
	loopLabel := getLabel(ob)

	// 设置循环地址
	setLabel(ob, loopLabel)

	if stmt.condition != nil {
		stmt.condition.generate(exe, currentBlock, ob)
	}

	parent := stmt.block.parent.(*StatementBlockInfo)

	// 获取break,continue地址
	parent.breakLabel = getLabel(ob)
	parent.continueLabel = getLabel(ob)

	if stmt.condition != nil {
		// 如果条件为否,跳转到break
		generateCode(ob, stmt.Position(), vm.VM_JUMP_IF_FALSE, parent.breakLabel)
	}

	generateStatementList(exe, stmt.block, stmt.block.statementList, ob)

	// 如果有continue,直接跳过block,从这里执行
	setLabel(ob, parent.continueLabel)

	if stmt.post != nil {
		stmt.post.generate(exe, currentBlock, ob)
	}

	// 跳回到循环开头
	generateCode(ob, stmt.Position(), vm.VM_JUMP, loopLabel)

	// 设置结束标签
	setLabel(ob, parent.breakLabel)
}

// ==============================
// ReturnStatement
// ==============================

// ReturnStatement return 语句
type ReturnStatement struct {
	StatementImpl
	// 返回值
	returnValue Expression
}

func (stmt *ReturnStatement) show(ident int) {
	printWithIdent("ReturnStmt", ident)
	subIdent := ident + 2

	stmt.returnValue.show(subIdent)
}

func (stmt *ReturnStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	var returnValue Expression

	stmt.returnValue.fix(currentBlock)

	if stmt.returnValue != nil {
		// 类型转换
		returnValue = createAssignCast(stmt.returnValue, fd.typeSpecifier)
		stmt.returnValue = returnValue
		return
	}

	if fd.typeSpecifier.deriveList != nil {
		panic("TODO")
	}

	switch fd.typeSpecifier.basicType {
	case vm.BooleanType:
		returnValue = &BooleanExpression{booleanValue: false}
	case vm.IntType:
		returnValue = &IntExpression{intValue: 0}
	case vm.DoubleType:
		stmt.returnValue = &DoubleExpression{doubleValue: 0.0}
	case vm.StringType:
		stmt.returnValue = &StringExpression{stringValue: ""}
	}

	returnValue.SetPosition(stmt.Position())
	stmt.returnValue = returnValue
}

func (stmt *ReturnStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.returnValue == nil {
		panic("Return value is nil.")
	}

	stmt.returnValue.generate(exe, currentBlock, ob)

	generateCode(ob, stmt.Position(), vm.VM_RETURN)
}

// ==============================
// BreakStatement
// ==============================

// BreakStatement break 语句
type BreakStatement struct {
	StatementImpl
}

func (stmt *BreakStatement) show(ident int) {
	printWithIdent("BreakStmt", ident)
}

func (stmt *BreakStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *BreakStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	generateCode(ob, stmt.Position(), vm.VM_JUMP, currentBlock.parent.(*StatementBlockInfo).breakLabel)
}

// ==============================
// ContinueStatement
// ==============================

// ContinueStatement continue 语句
type ContinueStatement struct {
	StatementImpl
}

func (stmt *ContinueStatement) show(ident int) {
	printWithIdent("ContinueStmt", ident)
}

func (stmt *ContinueStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *ContinueStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	generateCode(ob, stmt.Position(), vm.VM_JUMP, currentBlock.parent.(*StatementBlockInfo).continueLabel)
}

// ==============================
// Declaration
// ==============================

// Declaration 声明语句
type Declaration struct {
	StatementImpl

	typeSpecifier *TypeSpecifier

	name        string
	initializer Expression

	variableIndex int

	isLocal bool
}

func (stmt *Declaration) show(ident int) {
	printWithIdent("DeclStmt", ident)

	subIdent := ident + 2
	stmt.initializer.show(subIdent)
}

func (stmt *Declaration) fix(currentBlock *Block, fd *FunctionDefinition) {
	addDeclaration(currentBlock, stmt, fd, stmt.Position())

	stmt.initializer.fix(currentBlock)

	// 类型转换
	if stmt.initializer != nil {
		stmt.initializer = createAssignCast(stmt.initializer, stmt.typeSpecifier)
	}
}

func (stmt *Declaration) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.initializer == nil {
		return
	}

	stmt.initializer.generate(exe, currentBlock, ob)
	generatePopToIdentifier(stmt, stmt.Position(), ob)
}
