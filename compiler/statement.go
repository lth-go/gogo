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

type ArrayDerive struct {
}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl
	// 基本类型
	basicType vm.BasicType
	// 派生类型
	deriveList []TypeDerive
}

func (t *TypeSpecifier) appendDerive(derive TypeDerive) {
	if t.deriveList == nil {
		t.deriveList = []TypeDerive{}
	}
	t.deriveList = append(t.deriveList, derive)
}

func (t *TypeSpecifier) isArrayDerive() bool {
	return isArray(t)
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

	for _, decl := range b.declarationList {
		decl.show(subIdent)
	}

	for _, stmt := range b.statementList {
		stmt.show(subIdent)
	}
}

func (b *Block) addDeclaration(decl *Declaration, fd *FunctionDefinition, pos Position) {
	if searchDeclaration(decl.name, b) != nil {
		compileError(pos, VARIABLE_MULTIPLE_DEFINE_ERR, decl.name)
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

func (fd *FunctionDefinition) typeS() *TypeSpecifier {
	return fd.typeSpecifier
}

func (fd *FunctionDefinition) addParameterAsDeclaration() {

	for _, param := range fd.parameterList {
		if searchDeclaration(param.name, fd.block) != nil {
			compileError(param.Position(), PARAMETER_MULTIPLE_DEFINE_ERR, "parameter name: %s\n", param.name)
		}
		decl := &Declaration{name: param.name, typeSpecifier: param.typeSpecifier}

		fd.block.addDeclaration(decl, fd, param.Position())
	}
}

func (fd *FunctionDefinition) addReturnFunction() {

	if fd.block.statementList == nil {
		ret := &ReturnStatement{returnValue: nil}
		ret.fix(fd.block, fd)
		fd.block.statementList = []Statement{ret}
		return
	}

	// TODO return 是否有必要一定最后
	last := fd.block.statementList[len(fd.block.statementList)-1]
	_, ok := last.(*ReturnStatement)
	if ok {
		return
	}

	ret := &ReturnStatement{returnValue: nil}
    ret.SetPosition(fd.typeSpecifier.Position())
    if ret.returnValue != nil {
		ret.returnValue.SetPosition(fd.typeSpecifier.Position())
    }
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
		compileError(expr.Position(), ARGUMENT_COUNT_MISMATCH_ERR, length, len(argumentList))
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
		ob.generateCode(expr.Position(), vm.VM_POP)
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
	ifFalseLabel := ob.getLabel()
	ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

	if stmt.thenBlock != nil {
		generateStatementList(exe, stmt.thenBlock, stmt.thenBlock.statementList, ob)
	}

	// 获取结束跳转地址
	endLabel := ob.getLabel()

	// 直接跳到最后
	ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

	// 设置false跳转地址,如果false,直接执行这里
	ob.setLabel(ifFalseLabel)

	for _, elif := range stmt.elifList {
		elif.condition.generate(exe, currentBlock, ob)

		// 获取false跳转地址
		ifFalseLabel = ob.getLabel()
		ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(exe, elif.block, elif.block.statementList, ob)

		// 直接跳到最后
		ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

		// 设置false跳转地址,如果false,直接执行这里
		ob.setLabel(ifFalseLabel)
	}

	if stmt.elseBlock != nil {
		generateStatementList(exe, stmt.elseBlock, stmt.elseBlock.statementList, ob)
	}

	// 设置结束地址
	ob.setLabel(endLabel)
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

	if stmt.init != nil {
		stmt.init.show(subIdent)
	}
	if stmt.condition != nil {
		stmt.condition.show(subIdent)
	}
	if stmt.post != nil {
		stmt.post.show(subIdent)
	}

	if stmt.block != nil {
		stmt.block.show(subIdent)
	}
}

func (stmt *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	if stmt.init != nil {
		stmt.init.fix(currentBlock)
	}
	if stmt.condition != nil {
		stmt.condition.fix(currentBlock)
	}
	if stmt.post != nil {
		stmt.post.fix(currentBlock)
	}

	if stmt.block != nil {
		fixStatementList(stmt.block, stmt.block.statementList, fd)
	}
}
func (stmt *ForStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {

	if stmt.init != nil {
		stmt.init.generate(exe, currentBlock, ob)
	}

	// 获取循环地址
	loopLabel := ob.getLabel()

	// 设置循环地址
	ob.setLabel(loopLabel)

	if stmt.condition != nil {
		stmt.condition.generate(exe, currentBlock, ob)
	}

	label := ob.getLabel()

	if stmt.condition != nil {
		// 如果条件为否,跳转到break, label = parent.breakLabel
		ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, label)
	}

	if stmt.block != nil {
		parent := stmt.block.parent.(*StatementBlockInfo)
		// 获取break,continue地址
		parent.breakLabel = label
		parent.continueLabel = label

		generateStatementList(exe, stmt.block, stmt.block.statementList, ob)
	}

	// 如果有continue,直接跳过block,从这里执行, label = parent.continueLabel
	ob.setLabel(label)

	if stmt.post != nil {
		stmt.post.generate(exe, currentBlock, ob)
	}

	// 跳回到循环开头
	ob.generateCode(stmt.Position(), vm.VM_JUMP, loopLabel)

	// 设置结束标签, label = parent.breakLabel
	ob.setLabel(label)
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

	if stmt.returnValue != nil {
		stmt.returnValue.fix(currentBlock)
		// 类型转换
		returnValue = createAssignCast(stmt.returnValue, fd.typeSpecifier)
		stmt.returnValue = returnValue
		return
	}

	if fd.typeSpecifier.deriveList != nil {
		_, ok := fd.typeSpecifier.deriveList[0].(*ArrayDerive)
		if !ok {
			panic("TODO")
		}
		returnValue = &NullExpression{}
		stmt.returnValue = returnValue
		return
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
	case vm.NullType:
		fallthrough
	default:
		panic("TODO")
	}

	returnValue.SetPosition(stmt.Position())
	stmt.returnValue = returnValue
}

func (stmt *ReturnStatement) generate(exe *vm.Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.returnValue == nil {
		panic("Return value is nil.")
	}

	stmt.returnValue.generate(exe, currentBlock, ob)

	ob.generateCode(stmt.Position(), vm.VM_RETURN)
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
	// 向外寻找,直到找到for的block
	for block := currentBlock; block != nil; block = block.outerBlock {
		switch block.parent.(type) {
		case *StatementBlockInfo:
			ob.generateCode(stmt.Position(), vm.VM_JUMP, block.parent.(*StatementBlockInfo).breakLabel)
			return
		default:
			continue
		}
	}
	compileError(stmt.Position(), LABEL_NOT_FOUND_ERR)
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
	// 向外寻找,直到找到for的block
	for block := currentBlock; block != nil; block = block.outerBlock {
		switch block.parent.(type) {
		case *StatementBlockInfo:
			ob.generateCode(stmt.Position(), vm.VM_JUMP, block.parent.(*StatementBlockInfo).continueLabel)
			return
		default:
			continue
		}
	}
	compileError(stmt.Position(), LABEL_NOT_FOUND_ERR)

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
	if stmt.initializer != nil {
		stmt.initializer.show(subIdent)
	}
}

func (stmt *Declaration) fix(currentBlock *Block, fd *FunctionDefinition) {
	currentBlock.addDeclaration(stmt, fd, stmt.Position())

	// 类型转换
	if stmt.initializer != nil {
		stmt.initializer.fix(currentBlock)
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
