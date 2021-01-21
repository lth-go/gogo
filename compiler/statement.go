package compiler

import (
	"github.com/lth-go/gogo/vm"
)

// Statement 语句接口
type Statement interface {
	Pos
	fix(*Block, *FunctionDefinition)
	generate(currentBlock *Block, ob *OpCodeBuf)
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

func (stmt *ExpressionStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.expression = stmt.expression.fix(currentBlock)
}

func (stmt *ExpressionStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
	expr := stmt.expression
	expr.generate(currentBlock, ob)

	// 处理函数多返回值
	funcExpr, ok := expr.(*FunctionCallExpression)
	if ok {
		for i := 0; i < len(funcExpr.Type.funcType.Results); i++ {
			ob.generateCode(expr.Position(), vm.VM_POP)
		}
	} else {
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
type ElseIf struct {
	condition Expression
	block     *Block
}

// IfStatement if表达式
type IfStatement struct {
	StatementBase
	condition Expression
	thenBlock *Block
	elifList  []*ElseIf
	elseBlock *Block
}

func (stmt *IfStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.condition = stmt.condition.fix(currentBlock)

	if !stmt.condition.GetType().IsBool() {
		compileError(stmt.condition.Position(), IF_CONDITION_NOT_BOOLEAN_ERR)
	}

	if stmt.thenBlock != nil {
		stmt.thenBlock.FixStatementList(fd)
	}

	for _, elif := range stmt.elifList {
		elif.condition = elif.condition.fix(currentBlock)

		if elif.block != nil {
			elif.block.FixStatementList(fd)
		}
	}

	if stmt.elseBlock != nil {
		stmt.elseBlock.FixStatementList(fd)
	}
}

func (stmt *IfStatement) generate(currentBlock *Block, ob *OpCodeBuf) {

	stmt.condition.generate(currentBlock, ob)

	// 获取false跳转地址
	ifFalseLabel := ob.getLabel()
	ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

	if stmt.thenBlock != nil {
		generateStatementList(stmt.thenBlock, stmt.thenBlock.statementList, ob)
	}

	// 获取结束跳转地址
	endLabel := ob.getLabel()

	// 直接跳到最后
	ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

	// 设置false跳转地址,如果false,直接执行这里
	ob.setLabel(ifFalseLabel)

	for _, elif := range stmt.elifList {
		elif.condition.generate(currentBlock, ob)

		// 获取false跳转地址
		ifFalseLabel = ob.getLabel()
		ob.generateCode(stmt.Position(), vm.VM_JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(elif.block, elif.block.statementList, ob)

		// 直接跳到最后
		ob.generateCode(stmt.Position(), vm.VM_JUMP, endLabel)

		// 设置false跳转地址,如果false,直接执行这里
		ob.setLabel(ifFalseLabel)
	}

	if stmt.elseBlock != nil {
		generateStatementList(stmt.elseBlock, stmt.elseBlock.statementList, ob)
	}

	// 设置结束地址
	ob.setLabel(endLabel)
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

func (stmt *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	if stmt.init != nil {
		stmt.init.fix(currentBlock, fd)
	}

	if stmt.condition != nil {
		stmt.condition = stmt.condition.fix(currentBlock)

		if !stmt.condition.GetType().IsBool() {
			compileError(stmt.condition.Position(), FOR_CONDITION_NOT_BOOLEAN_ERR)
		}
	}

	if stmt.post != nil {
		stmt.post.fix(currentBlock, fd)
	}

	if stmt.block != nil {
		stmt.block.FixStatementList(fd)
	}
}

func (stmt *ForStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
	if stmt.init != nil {
		stmt.init.generate(currentBlock, ob)
	}

	// 获取循环地址
	loopLabel := ob.getLabel()

	// 设置循环地址
	ob.setLabel(loopLabel)

	if stmt.condition != nil {
		stmt.condition.generate(currentBlock, ob)
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

		generateStatementList(stmt.block, stmt.block.statementList, ob)
	}

	// 如果有continue,直接跳过block,从这里执行, label = parent.continueLabel
	ob.setLabel(label)

	if stmt.post != nil {
		stmt.post.generate(currentBlock, ob)
	}

	// 跳回到循环开头
	ob.generateCode(stmt.Position(), vm.VM_JUMP, loopLabel)

	// 设置结束标签, label = parent.breakLabel
	ob.setLabel(label)
}

//
// ReturnStatement
//
type ReturnStatement struct {
	StatementBase
	ValueList []Expression
}

func (stmt *ReturnStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	if len(fd.GetType().funcType.Results) == 0 && len(stmt.ValueList) == 0 {
		return
	} else if len(fd.GetType().funcType.Results) == 0 && len(stmt.ValueList) > 0 {
		// 函数没有定义返回值,却返回了
		compileError(stmt.Position(), RETURN_IN_VOID_FUNCTION_ERR)
	} else if len(fd.GetType().funcType.Results) != 0 && len(stmt.ValueList) == 0 {
		// 函数定义了返回值,却没返回
		compileError(stmt.Position(), BAD_RETURN_TYPE_ERR)
	} else {
		stmt.ValueList = []Expression{
			CreateAssignCast(stmt.ValueList[0].fix(currentBlock), fd.GetType().funcType.Results[0].Type),
		}
	}
}

func (stmt *ReturnStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
	for _, value := range stmt.ValueList {
		value.generate(currentBlock, ob)
	}
	ob.generateCode(stmt.Position(), vm.VM_RETURN)
}

func NewReturnStatement(pos Position, valueList []Expression) *ReturnStatement {
	stmt := &ReturnStatement{
		ValueList: valueList,
	}
	stmt.SetPosition(pos)

	return stmt
}

//
// BreakStatement
//
type BreakStatement struct {
	StatementBase
}

func (stmt *BreakStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *BreakStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
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

//
// ContinueStatement
//
type ContinueStatement struct {
	StatementBase
}

func (stmt *ContinueStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *ContinueStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
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

//
// Declaration 声明语句
//
type Declaration struct {
	StatementBase
	Type      *Type
	Name      string
	InitValue Expression
	Index     int
	IsLocal   bool
}

func (stmt *Declaration) fix(currentBlock *Block, fd *FunctionDefinition) {
	currentBlock.AddDeclaration(stmt, fd)

	stmt.Type.Fix()

	// 类型转换
	if stmt.InitValue != nil {
		stmt.InitValue = stmt.InitValue.fix(currentBlock)
		stmt.InitValue = CreateAssignCast(stmt.InitValue, stmt.Type)
	}
}

func (stmt *Declaration) generate(currentBlock *Block, ob *OpCodeBuf) {
	if stmt.InitValue == nil {
		return
	}

	stmt.InitValue.generate(currentBlock, ob)
	generatePopToIdentifier(stmt, stmt.Position(), ob)
}

func NewDeclaration(pos Position, typ *Type, name string, value Expression) *Declaration {
	decl := &Declaration{
		Type:      typ,
		Name:      name,
		InitValue: value,
		Index:     -1,
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

func (stmt *AssignStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	leftLen := len(stmt.left)
	rightLen := len(stmt.right)

	// 校验右边是否有函数调用,如果有取函数返回值为长度
	isCall := false
	for _, expr := range stmt.right {
		callExpr, ok := expr.(*FunctionCallExpression)
		if ok {
			if rightLen > 1 {
				panic("TODO")
			}
			// 先fix,否则type不对
			callExpr.fix(currentBlock)
			rightLen = len(callExpr.Type.funcType.Results)
			isCall = true
			break
		}
	}

	if leftLen != rightLen {
		panic("TODO")
	}

	for _, expr := range stmt.left {
		switch expr.(type) {
		case *IdentifierExpression, *IndexExpression, *SelectorExpression:
		default:
			compileError(expr.Position(), NOT_LVALUE_ERR, "")
		}
	}

	if isCall {
		for i := 0; i < len(stmt.left); i++ {
			leftExpr := stmt.left[i]
			leftExpr.fix(currentBlock)

		}
	} else {
		for i := 0; i < len(stmt.left); i++ {
			leftExpr := stmt.left[i]
			leftExpr.fix(currentBlock)

			rightExpr := stmt.right[i]
			rightExpr.fix(currentBlock)
			stmt.right[i] = CreateAssignCast(stmt.right[i], leftExpr.GetType())
		}
	}

}

func (stmt *AssignStatement) generate(currentBlock *Block, ob *OpCodeBuf) {
	isCall := false
	for _, expr := range stmt.right {
		_, ok := expr.(*FunctionCallExpression)
		if ok {
			isCall = true
			break
		}
	}

	if isCall {
		for _, expr := range stmt.right {
			expr.generate(currentBlock, ob)
		}

		count := len(stmt.left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.left[count-i-1]
			ob.generateCode(stmt.Position(), vm.VM_DUPLICATE)
			generatePopToLvalue(currentBlock, leftExpr, ob)
		}
	} else {
		count := len(stmt.left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.left[i]
			rightExpr := stmt.right[i]

			rightExpr.generate(currentBlock, ob)
			ob.generateCode(stmt.Position(), vm.VM_DUPLICATE)
			generatePopToLvalue(currentBlock, leftExpr, ob)
		}
	}
}
