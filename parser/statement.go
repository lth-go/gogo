package parser

// ==============================
// 基本类型
// ==============================

// BasicType 基础类型
type BasicType int

const (
	// BooleanType 布尔类型
	BooleanType BasicType = iota
	// NumberType 数字类型
	NumberType
	// StringType 字符串类型
	StringType
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
	basicType BasicType
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
	generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf)
}

// StatementImpl provide commonly implementations for Stmt..
type StatementImpl struct {
	PosImpl    // StmtImpl provide Pos() function.
	lineNumber int
}

// stmt provide restraint interface.
func (stmt *StatementImpl) stmt() {}

//
// Block 块接口
//

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

type Block struct {
	BlockType       int
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration

	// 块信息，函数块，还是条件语句
	parent BlockInfo
}

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	typeSpecifier *TypeSpecifier
	name          string
	parameterList []*Parameter
	block         *Block

	index int

	localVariable []*Declaration
}

func (fd *FunctionDefinition) typeS() *TypeSpecifier {
	return fd.typeSpecifier
}

// Parameter 形参
type Parameter struct {
	PosImpl
	typeSpecifier *TypeSpecifier
	name          string
}

// AssignmentOperator ...
type AssignmentOperator int

const (
	// NormalAssign 赋值操作符 =
	NormalAssign AssignmentOperator = iota
)

// ==============================
// ExpressionStatement
// ==============================

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	StatementImpl
	expression Expression
}

func (stmt *ExpressionStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.expression.fix(currentBlock)
}

func (stmt *ExpressionStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr := stmt.expression
	switch assignExpr := expr.(type) {
	case *AssignExpression:
		// TODO
		assignExpr.generateEx(exe, currentBlock, ob)
	default:
		expr.generate(exe, currentBlock, ob)
		generateCode(ob, expr.Position(), POP)
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

func (stmt *IfStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

	stmt.condition.fix(currentBlock)

	fixStatementList(stmt.thenBlock, stmt.thenBlock.statementList, fd)

	for _, elifPos := range stmt.elifList {
		elifPos.condition.fix(currentBlock)

		if elifPos.block != nil {
			fixStatementList(elifPos.block, elifPos.block.statementList, fd)
		}
	}

	if stmt.elifList != nil {
		fixStatementList(stmt.elseBlock, stmt.elseBlock.statementList, fd)
	}

}

func (stmt *IfStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	stmt.condition.generate(exe, currentBlock, ob)

	ifFalseLabel = getLabel(ob)

	generateCode(ob, stmt.Position(), JUMP_IF_FALSE, ifFalseLabel)

	generateStatementList(exe, stmt.thenBlock, stmt.thenBlock.statementList, ob)

	endLabel = getLabel(ob)

	generateCode(ob, statement.Position(), JUMP, endLabel)

	setLabel(ob, ifFalseLabel)

	for _, elif := range stmt.elifList {
		elif.condition.generate(exe, currentBlock, ob)
		ifFalseLabel = getLabel(ob)

		generateCode(ob, statement.Position(), JUMP_IF_FALSE, ifFalseLabel)

		generateStatementList(exe, elif.block, elif.block.statementList, ob)
		generateCode(ob, statement.Position(), JUMP, endLabel)

		setLabel(ob, ifFalseLabel)
	}
	if stmt.elseBlock != nil {
		generateStatementList(exe, if_s.else_block, if_s.else_block.statement_list, ob)
	}

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

	label     string
	init      Expression
	condition Expression
	post      Expression
	block     *Block
}

func (stmt *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.init.fix(currentBlock)
	stmt.condition.fix(currentBlock)
	stmt.post.fix(currentBlock)
	fixStatementList(stmt.block, stmt.block.statementList, fd)

}
func (stmt *ForStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	if stmt.init != nil {
		stmt.generate(exe, currentBlock, ob)
	}
	loop_label = getLabel(ob)

	setLabel(ob, loop_label)

	if stmt.condition != nil {
		stmt.condition.generate(exe, currentBlock, ob)
	}

	parent, ok := stmt.block.parent.(*StatementBlockInfo)
	if !ok {
		compileError(stmt.Position(), 0, "")
	}

	parent.breakLabel = getLabel(ob)
	parent.continueLabel = getLabel(ob)

	if stmt.condition != nil {
		generateCode(ob, stmt.Position(), JUMP_IF_FALSE, parent.breakLabel)
	}

	generateStatementList(exe, stmt.block, stmt.block.statementList, ob)

	setLabel(ob, parent.continueLabel)

	if stmt.post != nil {
		stmt.post.generate(exe, currentBlock, ob)
	}

	generateCode(ob, statement.Position(), JUMP, loop_label)

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

func (stmt *ReturnStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	stmt.returnValue.fix(currentBlock)

	if stmt.returnValue == nil {
		switch fd.typeSpecifier.basicType {
		case BooleanType:
			stmt.returnValue = &BooleanExpression{
				booleanValue:  false,
				typeSpecifier: &TypeSpecifier{basicType: BooleanType},
			}
			return
		case NumberType:
			stmt.returnValue = &NumberExpression{
				numberValue:   0.0,
				typeSpecifier: &TypeSpecifier{basicType: NumberType},
			}
			return
		case StringType:
			stmt.returnValue = &StringExpression{
				stringValue:   "",
				typeSpecifier: &TypeSpecifier{basicType: StringType},
			}
			return
		}

	}
	// 类型转换
	createAssignCast(stmt.returnValue, fd.typeSpecifier)
}

func (stmt *ReturnStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.returnValue == nil {
		compileError(stmt.Position(), 0, "")
	}

	stmt.returnValue.generate(exe, currentBlock, ob)

	generateCode(ob, stmt.Position(), RETURN)
}

// ==============================
// BreakStatement
// ==============================

// BreakStatement break 语句
type BreakStatement struct {
	StatementImpl

	label string
}

func (stmt *BreakStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *BreakStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	var parent *StatementBlockInfo

	for block := currentBlock; block != nil; block = block.outerBlock {
		parent, ok := block.parent.(*StatementBlockInfo)
		if !ok {
			continue
		}

		if stmt.label == "" {
			break
		}

		parentFor, ok := parent.statement.(*ForStatement)
		if !ok {
			compileError(stmt.Position(), 0, "")
		}
		if parentFor.label == "" {
			continue
		}

		if stmt.label != parentFor.label {
			break
		}
	}

	if block == nil {
		compileError(stmt.Position(), 0, "")
	}

	generateCode(ob, statement.Position(), JUMP, parent.breakLabel)

}

// ==============================
// ContinueStatement
// ==============================

// ContinueStatement continue 语句
type ContinueStatement struct {
	StatementImpl

	label string
}

func (stmt *ContinueStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

}
func (stmt *ContinueStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	var parent *StatementBlockInfo

	for block := currentBlock; block != nil; block = block.outerBlock {
		if parent, ok := block.parent.(*StatementBlockInfo); !ok {
			continue
		}

		if stmt.label == "" {
			break
		}

		if parentFor, ok := parent.statement.(*ForStatement); !ok {
			compileError(stmt.Position(), 0, "")
		}

		if parentFor.label == "" {
			continue
		}

		if stmt.label != parentFor.label {
			break
		}
	}

	if block == nil {
		dkc_compile_error(statement.Position(), 0, "")
	}

	generateCode(ob, statement.Position(), JUMP, block.parent.statement.continueLabel)

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

func (stmt *Declaration) fix(currentBlock *Block, fd *FunctionDefinition) {
	if searchDeclaration(stmt.name, currentBlock) != nil {
		compileError(stmt.Position(), 0, "")
	}

	if currentBlock != nil {
		currentBlock.declarationList = append(currentBlock.declarationList, stmt)
		addLocalVariable(fd, stmt)
		stmt.isLocal = true
	} else {
		compiler := getCurrentCompiler()
		compiler.declarationList = append(compiler.declarationList, stmt)
		stmt.isLocal = false
	}

	stmt.initializer.fix(currentBlock)

	// 类型转换
	if stmt.initializer != nil {
		stmt.initializer = createAssignCast(stmt.initializer, stmt.typeSpecifier)
	}
}
func (stmt *Declaration) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.initializer == nil {
		return
	}

	stmt.initializer.generate(exe, currentBlock, ob)
	generate_pop_to_identifier(decl, statement.Position(), ob)
}
