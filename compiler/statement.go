package compiler

import (
	"github.com/lth-go/gogo/vm"
)

// Statement 语句接口
type Statement interface {
	Pos
	fix()
	Generate(ob *OpCodeBuf)
}

type StatementBase struct {
	PosBase
}

//
// ExpressionStatement 表达式语句
//
type ExpressionStatement struct {
	StatementBase
	expression Expression
}

func (stmt *ExpressionStatement) fix() {
	stmt.expression = stmt.expression.fix()
}

func (stmt *ExpressionStatement) Generate(ob *OpCodeBuf) {
	expr := stmt.expression
	expr.generate(ob)

	// TODO: 没有返回值也会pop
	for i := 0; i < expr.GetType().GetResultCount(); i++ {
		ob.generateCode(expr.Position(), vm.VM_POP)
	}
}

func NewExpressionStatement(pos Position, expr Expression) *ExpressionStatement {
	stmt := &ExpressionStatement{
		expression: expr,
	}
	stmt.SetPosition(pos)

	return stmt
}

//
// IfStatement
//
// IfStatement if表达式
type IfStatement struct {
	StatementBase
	condition Expression
	thenBlock *Block
	elifList  []*ElseIf
	elseBlock *Block
}

func (stmt *IfStatement) fix() {
	stmt.condition = stmt.condition.fix()

	if !stmt.condition.GetType().IsBool() {
		compileError(stmt.condition.Position(), IF_CONDITION_NOT_BOOLEAN_ERR)
	}

	if stmt.thenBlock != nil {
		stmt.thenBlock.FixStatementList()
	}

	for _, elif := range stmt.elifList {
		elif.condition = elif.condition.fix()

		if elif.block != nil {
			elif.block.FixStatementList()
		}
	}

	if stmt.elseBlock != nil {
		stmt.elseBlock.FixStatementList()
	}
}

func (stmt *IfStatement) Generate(ob *OpCodeBuf) {

	stmt.condition.generate(ob)

	// 获取false跳转地址
	ifFalseLabel := ob.getLabel()
	ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

	if stmt.thenBlock != nil {
		generateStatementList(stmt.thenBlock.statementList, ob)
	}

	// 获取结束跳转地址
	endLabel := ob.getLabel()

	// 直接跳到最后
	ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

	// 设置false跳转地址,如果false,直接执行这里
	ob.setLabel(ifFalseLabel)

	for _, elif := range stmt.elifList {
		elif.condition.generate(ob)

		// 获取false跳转地址
		ifFalseLabel = ob.getLabel()
		ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(elif.block.statementList, ob)

		// 直接跳到最后
		ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

		// 设置false跳转地址,如果false,直接执行这里
		ob.setLabel(ifFalseLabel)
	}

	if stmt.elseBlock != nil {
		generateStatementList(stmt.elseBlock.statementList, ob)
	}

	// 设置结束地址
	ob.setLabel(endLabel)
}

func NewIfStatement(
	pos Position,
	condition Expression,
	thenBlock *Block,
	elifList []*ElseIf,
	elseBlock *Block,
) *IfStatement {
	stmt := &IfStatement{
		condition: condition,
		thenBlock: thenBlock,
		elifList:  elifList,
		elseBlock: elseBlock,
	}

	stmt.SetPosition(pos)

	return stmt
}

//
// ForStatement
//
type ForStatement struct {
	StatementBase
	init      Statement
	condition Expression
	post      Statement
	block     *Block
}

func (stmt *ForStatement) fix() {
	if stmt.init != nil {
		stmt.init.fix()
	}

	if stmt.condition != nil {
		stmt.condition = stmt.condition.fix()

		if !stmt.condition.GetType().IsBool() {
			compileError(stmt.condition.Position(), FOR_CONDITION_NOT_BOOLEAN_ERR)
		}
	}

	if stmt.post != nil {
		stmt.post.fix()
	}

	if stmt.block != nil {
		stmt.block.FixStatementList()
	}
}

func (stmt *ForStatement) Generate(ob *OpCodeBuf) {
	if stmt.init != nil {
		stmt.init.Generate(ob)
	}

	// 获取循环地址
	loopLabel := ob.getLabel()

	// 设置循环地址
	ob.setLabel(loopLabel)

	if stmt.condition != nil {
		stmt.condition.generate(ob)
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

		generateStatementList(stmt.block.statementList, ob)
	}

	// 如果有continue,直接跳过block,从这里执行, label = parent.continueLabel
	ob.setLabel(label)

	if stmt.post != nil {
		stmt.post.Generate(ob)
	}

	// 跳回到循环开头
	ob.generateCode(stmt.Position(), vm.VM_JUMP, loopLabel)

	// 设置结束标签, label = parent.breakLabel
	ob.setLabel(label)
}

func NewForStatement(pos Position, init Statement, condition Expression, post Statement, block *Block) *ForStatement {
	stmt := &ForStatement{
		init:      init,
		condition: condition,
		post:      post,
		block:     block,
	}

	stmt.SetPosition(pos)

	return stmt
}

type ElseIf struct {
	condition Expression
	block     *Block
}

func NewElseIf(condition Expression, block *Block) *ElseIf {
	return &ElseIf{
		condition: condition,
		block:     block,
	}
}

//
// ReturnStatement
//
type ReturnStatement struct {
	StatementBase
	ValueList []Expression
	Block     *Block
}

func (stmt *ReturnStatement) fix() {
	fd := stmt.Block.getCurrentFunction()

	if len(fd.GetType().funcType.Results) == 0 && len(stmt.ValueList) == 0 {
		return
	} else if len(fd.GetType().funcType.Results) == 0 && len(stmt.ValueList) > 0 {
		// 函数没有定义返回值,却返回了
		compileError(stmt.Position(), RETURN_IN_VOID_FUNCTION_ERR)
	} else if len(fd.GetType().funcType.Results) != 0 && len(stmt.ValueList) == 0 {
		// 函数定义了返回值,却没返回
		compileError(stmt.Position(), BAD_RETURN_TYPE_ERR)
	} else {
		// 只定义了单个返回值
		stmt.ValueList = []Expression{stmt.ValueList[0].fix()}

		if !fd.GetType().funcType.Results[0].Type.Equal(stmt.ValueList[0].GetType()) {
			compileError(stmt.Position(), BAD_RETURN_TYPE_ERR)
		}
	}
}

func (stmt *ReturnStatement) Generate(ob *OpCodeBuf) {
	for _, value := range stmt.ValueList {
		value.generate(ob)
	}
	ob.generateCode(stmt.Position(), vm.VM_RETURN)
}

func NewReturnStatement(pos Position, valueList []Expression) *ReturnStatement {
	stmt := &ReturnStatement{
		ValueList: valueList,
	}
	stmt.SetPosition(pos)

	stmt.Block = GetCurrentCompiler().currentBlock

	return stmt
}

//
// BreakStatement
//
type BreakStatement struct {
	StatementBase
	Block *Block
}

func (stmt *BreakStatement) fix() {
	for block := stmt.Block; block != nil; block = block.outerBlock {
		switch block.parent.(type) {
		case *StatementBlockInfo:
			return
		}
	}
	compileError(stmt.Position(), LABEL_NOT_FOUND_ERR)
}

func (stmt *BreakStatement) Generate(ob *OpCodeBuf) {
	// 向外寻找,直到找到for的block
	for block := stmt.Block; block != nil; block = block.outerBlock {
		switch block.parent.(type) {
		case *StatementBlockInfo:
			ob.generateCode(stmt.Position(), vm.VM_JUMP, block.parent.(*StatementBlockInfo).breakLabel)
			return
		}
	}
	panic("TODO")
}

func NewBreakStatement(pos Position) *BreakStatement {
	stmt := &BreakStatement{}
	stmt.SetPosition(pos)
	stmt.Block = GetCurrentCompiler().currentBlock

	return stmt
}

//
// ContinueStatement
//
type ContinueStatement struct {
	StatementBase
	Block *Block
}

func (stmt *ContinueStatement) fix() {}

func (stmt *ContinueStatement) Generate(ob *OpCodeBuf) {
	// 向外寻找,直到找到for的block
	for block := stmt.Block; block != nil; block = block.outerBlock {
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

func NewContinueStatement(pos Position) *ContinueStatement {
	stmt := &ContinueStatement{}
	stmt.SetPosition(pos)
	stmt.Block = GetCurrentCompiler().currentBlock

	return stmt
}

//
// Declaration 声明语句
//
type Declaration struct {
	StatementBase
	Type        *Type
	PackageName string
	Name        string
	InitValue   Expression
	Index       int
	IsLocal     bool
	Block       *Block
}

func (stmt *Declaration) fix() {
	fd := stmt.Block.getCurrentFunction()
	stmt.IsLocal = true
	stmt.Block.declarationList = append(stmt.Block.declarationList, stmt)
	fd.AddDeclarationList(stmt)

	stmt.Type.Fix()

	// 类型转换
	if stmt.InitValue != nil {
		stmt.InitValue = stmt.InitValue.fix()
		stmt.InitValue = CreateAssignCast(stmt.InitValue, stmt.Type)
		stmt.InitValue = stmt.InitValue.fix()
	}
}

func (stmt *Declaration) Generate(ob *OpCodeBuf) {
	if stmt.InitValue == nil {
		return
	}

	stmt.InitValue.generate(ob)
	generatePopToIdentifier(stmt, stmt.Position(), ob)
}

func NewDeclaration(pos Position, typ *Type, name string, value Expression) *Declaration {
	decl := &Declaration{
		Type:        typ,
		PackageName: "",
		Name:        name,
		InitValue:   value,
		Index:       -1,
		Block:       GetCurrentCompiler().currentBlock,
	}
	decl.SetPosition(pos)

	return decl
}

//
// AssignStatement
//
type AssignStatement struct {
	StatementBase
	left  []Expression
	right []Expression
}

func (stmt *AssignStatement) fix() {
	for _, expr := range stmt.left {
		switch expr.(type) {
		case *IdentifierExpression, *IndexExpression, *SelectorExpression:
		default:
			compileError(expr.Position(), NOT_LVALUE_ERR, "")
		}
	}

	leftLen := len(stmt.left)
	rightLen := len(stmt.right)

	// 校验右边是否有函数调用,如果有取函数返回值为长度
	isCall := stmt.isFuncCall()
	if isCall {
		if rightLen > 1 {
			panic("TODO")
		}

		stmt.right[0] = stmt.right[0].fix()
		rightLen = stmt.right[0].GetType().GetResultCount()
	}

	if leftLen != rightLen {
		panic("TODO")
	}

	if isCall {
		for i := 0; i < len(stmt.left); i++ {
			stmt.left[i] = stmt.left[i].fix()
		}
	} else {
		for i := 0; i < len(stmt.left); i++ {
			stmt.left[i] = stmt.left[i].fix()
			stmt.right[i] = stmt.right[i].fix()
			stmt.right[i] = CreateAssignCast(stmt.right[i], stmt.left[i].GetType())
			// TODO:
			stmt.right[i] = stmt.right[i].fix()
		}
	}

}

func (stmt *AssignStatement) Generate(ob *OpCodeBuf) {
	isCall := stmt.isFuncCall()

	if isCall {
		for _, expr := range stmt.right {
			expr.generate(ob)
		}

		count := len(stmt.left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.left[count-i-1]
			ob.generateCode(stmt.Position(), vm.VM_DUPLICATE)
			generatePopToLvalue(leftExpr, ob)
		}
	} else {
		count := len(stmt.left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.left[i]
			rightExpr := stmt.right[i]

			rightExpr.generate(ob)
			ob.generateCode(stmt.Position(), vm.VM_DUPLICATE)
			generatePopToLvalue(leftExpr, ob)
		}
	}
}

func (stmt *AssignStatement) isFuncCall() bool {
	for _, expr := range stmt.right {
		_, ok := expr.(*FunctionCallExpression)
		if ok {
			return true
		}
	}
	return false
}

func NewAssignStatement(pos Position, left []Expression, right []Expression) *AssignStatement {
	stmt := &AssignStatement{
		left:  left,
		right: right,
	}
	stmt.SetPosition(pos)

	return stmt
}
