package compiler

import (
	"github.com/lth-go/gogo/vm"
)

// Statement 语句接口
type Statement interface {
	Pos
	Fix()
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
	Expression Expression
}

func (stmt *ExpressionStatement) Fix() {
	stmt.Expression = stmt.Expression.Fix()
}

func (stmt *ExpressionStatement) Generate(ob *OpCodeBuf) {
	expr := stmt.Expression
	expr.Generate(ob)

	for i := 0; i < expr.GetType().GetResultCount(); i++ {
		ob.GenerateCode(expr.Position(), vm.OP_CODE_POP)
	}
}

func NewExpressionStatement(pos Position, expr Expression) *ExpressionStatement {
	stmt := &ExpressionStatement{
		Expression: expr,
	}
	stmt.SetPosition(pos)

	return stmt
}

//
// IfStatement if表达式
//
type IfStatement struct {
	StatementBase
	Condition  Expression
	ThenBlock  *Block
	ElseIfList []*ElseIf
	ElseBlock  *Block
}

func (stmt *IfStatement) Fix() {
	stmt.Condition = stmt.Condition.Fix()

	if !stmt.Condition.GetType().IsBool() {
		compileError(stmt.Condition.Position(), IF_CONDITION_NOT_BOOLEAN_ERR)
	}

	if stmt.ThenBlock != nil {
		stmt.ThenBlock.Fix()
	}

	for _, elseIf := range stmt.ElseIfList {
		elseIf.Condition = elseIf.Condition.Fix()

		if elseIf.Block != nil {
			elseIf.Block.Fix()
		}
	}

	if stmt.ElseBlock != nil {
		stmt.ElseBlock.Fix()
	}
}

func (stmt *IfStatement) Generate(ob *OpCodeBuf) {
	stmt.Condition.Generate(ob)

	// 获取false跳转地址
	ifFalseLabel := ob.GetLabel()
	ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP_IF_FALSE, ifFalseLabel)

	if stmt.ThenBlock != nil {
		generateStatementList(stmt.ThenBlock.statementList, ob)
	}

	// 获取结束跳转地址
	endLabel := ob.GetLabel()

	// 直接跳到最后
	ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP, endLabel)

	// 设置false跳转地址,如果false,直接执行这里
	ob.SetLabel(ifFalseLabel)

	for _, elif := range stmt.ElseIfList {
		elif.Condition.Generate(ob)

		// 获取false跳转地址
		ifFalseLabel = ob.GetLabel()
		ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(elif.Block.statementList, ob)

		// 直接跳到最后
		ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP, endLabel)

		// 设置false跳转地址,如果false,直接执行这里
		ob.SetLabel(ifFalseLabel)
	}

	if stmt.ElseBlock != nil {
		generateStatementList(stmt.ElseBlock.statementList, ob)
	}

	// 设置结束地址
	ob.SetLabel(endLabel)
}

func NewIfStatement(
	pos Position,
	condition Expression,
	thenBlock *Block,
	elifList []*ElseIf,
	elseBlock *Block,
) *IfStatement {
	stmt := &IfStatement{
		Condition:  condition,
		ThenBlock:  thenBlock,
		ElseIfList: elifList,
		ElseBlock:  elseBlock,
	}

	stmt.SetPosition(pos)

	return stmt
}

type ElseIf struct {
	Condition Expression
	Block     *Block
}

func NewElseIf(condition Expression, block *Block) *ElseIf {
	return &ElseIf{
		Condition: condition,
		Block:     block,
	}
}

//
// ForStatement
//
type ForStatement struct {
	StatementBase
	Init      Statement
	Condition Expression
	Post      Statement
	Block     *Block
}

func (stmt *ForStatement) Fix() {
	if stmt.Init != nil {
		stmt.Init.Fix()
	}

	if stmt.Condition != nil {
		stmt.Condition = stmt.Condition.Fix()

		if !stmt.Condition.GetType().IsBool() {
			compileError(stmt.Condition.Position(), FOR_CONDITION_NOT_BOOLEAN_ERR)
		}
	}

	if stmt.Post != nil {
		stmt.Post.Fix()
	}

	if stmt.Block != nil {
		stmt.Block.Fix()
	}
}

func (stmt *ForStatement) Generate(ob *OpCodeBuf) {
	if stmt.Init != nil {
		stmt.Init.Generate(ob)
	}

	// 获取循环地址
	loopLabel := ob.GetLabel()

	// 设置循环地址
	ob.SetLabel(loopLabel)

	if stmt.Condition != nil {
		stmt.Condition.Generate(ob)
	}

	label := ob.GetLabel()

	if stmt.Condition != nil {
		// 如果条件为否,跳转到break, label = parent.breakLabel
		ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP_IF_FALSE, label)
	}

	if stmt.Block != nil {
		parent := stmt.Block.parent.(*StatementBlockInfo)
		// 获取break,continue地址
		parent.BreakLabel = label
		parent.ContinueLabel = label

		generateStatementList(stmt.Block.statementList, ob)
	}

	// 如果有continue,直接跳过block,从这里执行, label = parent.continueLabel
	ob.SetLabel(label)

	if stmt.Post != nil {
		stmt.Post.Generate(ob)
	}

	// 跳回到循环开头
	ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP, loopLabel)

	// 设置结束标签, label = parent.breakLabel
	ob.SetLabel(label)
}

func NewForStatement(pos Position, init Statement, condition Expression, post Statement, block *Block) *ForStatement {
	stmt := &ForStatement{
		Init:      init,
		Condition: condition,
		Post:      post,
		Block:     block,
	}

	stmt.SetPosition(pos)

	return stmt
}

//
// ReturnStatement
//
type ReturnStatement struct {
	StatementBase
	ValueList []Expression
	Block     *Block
}

func (stmt *ReturnStatement) Fix() {
	fd := stmt.Block.GetCurrentFunction()

	resultCount := len(fd.GetType().funcType.Results)
	valueCount := len(stmt.ValueList)

	if resultCount == 0 && valueCount == 0 {
		return
	} else if resultCount == 0 && valueCount > 0 {
		// 函数没有定义返回值,却返回了
		compileError(stmt.Position(), RETURN_IN_VOID_FUNCTION_ERR)
	} else if resultCount != 0 && valueCount == 0 {
		// 函数定义了返回值,却没返回
		compileError(stmt.Position(), BAD_RETURN_TYPE_ERR)
	} else {
		// 只定义了单个返回值
		stmt.ValueList = []Expression{stmt.ValueList[0].Fix()}

		if !fd.GetType().funcType.Results[0].Type.Equal(stmt.ValueList[0].GetType()) {
			compileError(stmt.Position(), BAD_RETURN_TYPE_ERR)
		}
	}
}

func (stmt *ReturnStatement) Generate(ob *OpCodeBuf) {
	for _, value := range stmt.ValueList {
		value.Generate(ob)
	}
	ob.GenerateCode(stmt.Position(), vm.OP_CODE_RETURN)
}

func NewReturnStatement(pos Position, valueList []Expression) *ReturnStatement {
	stmt := &ReturnStatement{
		ValueList: valueList,
	}
	stmt.SetPosition(pos)

	stmt.Block = GetCurrentPackage().currentBlock

	return stmt
}

//
// BreakStatement
//
type BreakStatement struct {
	StatementBase
	Block *Block
}

func (stmt *BreakStatement) Fix() {
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
			ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP, block.parent.(*StatementBlockInfo).BreakLabel)
			return
		}
	}
	panic("TODO")
}

func NewBreakStatement(pos Position) *BreakStatement {
	stmt := &BreakStatement{}
	stmt.SetPosition(pos)
	stmt.Block = GetCurrentPackage().currentBlock

	return stmt
}

//
// ContinueStatement
//
type ContinueStatement struct {
	StatementBase
	Block *Block
}

func (stmt *ContinueStatement) Fix() {}

func (stmt *ContinueStatement) Generate(ob *OpCodeBuf) {
	// 向外寻找,直到找到for的block
	for block := stmt.Block; block != nil; block = block.outerBlock {
		switch block.parent.(type) {
		case *StatementBlockInfo:
			ob.GenerateCode(stmt.Position(), vm.OP_CODE_JUMP, block.parent.(*StatementBlockInfo).ContinueLabel)
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
	stmt.Block = GetCurrentPackage().currentBlock

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
	Value       Expression
	Index       int    // 下标
	IsLocal     bool   // 是否本地声明
	Block       *Block // 所属块
}

func (stmt *Declaration) Fix() {
	//
	// 向父块添加
	//
	block := stmt.Block
	block.declarationList = append(block.declarationList, stmt)

	// 向父函数添加
	fd := block.GetCurrentFunction()
	stmt.Index = len(fd.DeclarationList)
	fd.DeclarationList = append(fd.DeclarationList, stmt)

	stmt.Type.Fix()

	//
	// 设置类型默认值
	//
	if stmt.Value == nil {
		stmt.Value = GetTypeDefaultValue(stmt.Type, stmt.Position())
	}

	stmt.Value = stmt.Value.Fix()
	stmt.Value = CreateAssignCast(stmt.Value, stmt.Type)
}

func (stmt *Declaration) Generate(ob *OpCodeBuf) {
	stmt.Value.Generate(ob)
	generatePopToIdentifier(stmt, stmt.Position(), ob)
}

func NewDeclaration(pos Position, typ *Type, name string, value Expression) *Declaration {
	decl := &Declaration{
		Type:        typ,
		PackageName: "",
		Name:        name,
		Value:       value,
		Index:       -1,
	}
	decl.SetPosition(pos)

	return decl
}

func GetTypeDefaultValue(typ *Type, pos Position) Expression {
	if typ.IsArray() || typ.IsMap() || typ.IsInterface() {
		return CreateNilExpression(pos)
	}

	switch typ.GetBasicType() {
	case BasicTypeBool:
		return CreateBooleanExpression(pos, false)
	case BasicTypeInt:
		return CreateIntExpression(pos, 0)
	case BasicTypeFloat:
		return CreateFloatExpression(pos, 0.0)
	case BasicTypeString:
		return CreateStringExpression(pos, "")
	case BasicTypeStruct:
		// TODO: 结构体默认值
		// value := CreateStructExpression(typ, nil)

		// return value
		fallthrough
	default:
		panic("TODO")
	}
}

//
// AssignStatement
//
type AssignStatement struct {
	StatementBase
	Left  []Expression
	Right []Expression
}

func (stmt *AssignStatement) Fix() {
	//
	// 检查左值类型
	//
	for _, expr := range stmt.Left {
		switch expr.(type) {
		case *IdentifierExpression, *IndexExpression, *SelectorExpression:
		default:
			compileError(expr.Position(), NOT_LVALUE_ERR, "")
		}
	}

	// 校验右边是否有函数调用,如果有取函数返回值为长度
	if stmt.isFuncCall() {
		leftLen := len(stmt.Left)
		rightLen := len(stmt.Right)

		if rightLen != 1 {
			panic("TODO")
		}

		stmt.Right[0] = stmt.Right[0].Fix()
		rightLen = stmt.Right[0].GetType().GetResultCount()

		if leftLen != rightLen {
			panic("TODO")
		}

		for i := 0; i < len(stmt.Left); i++ {
			stmt.Left[i] = stmt.Left[i].Fix()
		}
	} else {
		leftLen := len(stmt.Left)
		rightLen := len(stmt.Right)

		if leftLen != rightLen {
			panic("TODO")
		}

		for i := 0; i < len(stmt.Left); i++ {
			stmt.Left[i] = stmt.Left[i].Fix()
			stmt.Right[i] = stmt.Right[i].Fix()
			stmt.Right[i] = CreateAssignCast(stmt.Right[i], stmt.Left[i].GetType())
		}
	}
}

func (stmt *AssignStatement) Generate(ob *OpCodeBuf) {
	isCall := stmt.isFuncCall()

	if isCall {
		for _, expr := range stmt.Right {
			expr.Generate(ob)
		}

		count := len(stmt.Left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.Left[count-i-1]
			ob.GenerateCode(stmt.Position(), vm.OP_CODE_DUPLICATE)
			generatePopToLvalue(leftExpr, ob)
		}
	} else {
		count := len(stmt.Left)
		for i := 0; i < count; i++ {
			leftExpr := stmt.Left[i]
			rightExpr := stmt.Right[i]

			rightExpr.Generate(ob)
			ob.GenerateCode(stmt.Position(), vm.OP_CODE_DUPLICATE)
			generatePopToLvalue(leftExpr, ob)
		}
	}
}

func (stmt *AssignStatement) isFuncCall() bool {
	for _, expr := range stmt.Right {
		_, ok := expr.(*CallExpression)
		if ok {
			return true
		}
	}
	return false
}

func NewAssignStatement(pos Position, left []Expression, right []Expression) *AssignStatement {
	stmt := &AssignStatement{
		Left:  left,
		Right: right,
	}

	stmt.SetPosition(pos)

	return stmt
}

func generateStatementList(statementList []Statement, ob *OpCodeBuf) {
	for _, stmt := range statementList {
		stmt.Generate(ob)
	}
}

func generatePopToLvalue(expr Expression, ob *OpCodeBuf) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		generatePopToIdentifier(e.Obj.(*Declaration), expr.Position(), ob)
	case *IndexExpression:
		if e.X.GetType().IsArray() {
			e.X.Generate(ob)
			e.Index.Generate(ob)
			ob.GenerateCode(expr.Position(), vm.OP_CODE_POP_ARRAY)
		} else if e.X.GetType().IsMap() {
			e.X.Generate(ob)
			e.Index.Generate(ob)
			ob.GenerateCode(expr.Position(), vm.OP_CODE_POP_MAP)
		} else {
			panic("TODO")
		}
	case *SelectorExpression:
		if e.X.GetType().IsStruct() {
			e.X.Generate(ob)
			ob.GenerateCode(expr.Position(), vm.OP_CODE_PUSH_INT_2BYTE, e.Index)
			ob.GenerateCode(expr.Position(), vm.OP_CODE_POP_STRUCT)
		} else {
			panic("TODO")
		}

	default:
		panic("TODO")
	}
}

func generatePopToIdentifier(decl *Declaration, pos Position, ob *OpCodeBuf) {
	var code byte

	if decl.IsLocal {
		code = vm.OP_CODE_POP_STACK
	} else {
		code = vm.OP_CODE_POP_STATIC
	}

	ob.GenerateCode(pos, code, decl.Index)
}
