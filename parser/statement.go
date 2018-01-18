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
	parameterList *[]Parameter
}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl
	// 基本类型
	basicType BasicType
	// 派生类型
	derive TypeDerive
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
func (s *StatementImpl) stmt() {}

//
// Block 块接口
//

type BlockInfo interface{}

type StatementBlockInfo struct {
	statement      Statement
	continue_label int
	break_label    int
}

type FunctionBlockInfo struct {
	function  *FunctionDefinition
	end_label int
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

func (s *ExpressionStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	s.expression.fix(currentBlock)
}

func (stmt *ExpressionStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	expr := stmt.expression
	switch expr.(type) {
	case *AssignExpression:
		expr.generate(exe, currentBlock, ob, true)
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

func (s *IfStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

	s.condition.fix(currentBlock)

	fixStatementList(s.thenBlock, s.thenBlock.statementList, fd)

	for _, elifPos := range s.elifList {
		elifPos.condition.fix(currentBlock)

		if elifPos.block != nil {
			fixStatementList(elifPos.block, elifPos.block.statementList, fd)
		}
	}

	if s.elifList != nil {
		fixStatementList(s.elseBlock, s.elseBlock.statementList, fd)
	}

}

func (stmt *IfStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	stmt.condition.generate(exe, currentBlock, ob)

	if_false_label = get_label(ob)

	generate_code(ob, stmt.Position(), JUMP_IF_FALSE, if_false_label)

	generate_statement_list(exe, stmt.thenBlock, stmt.thenBlock.statementList, ob)
	end_label = get_label(ob)
	generate_code(ob, statement.Position(), DVM_JUMP, end_label)

	set_label(ob, if_false_label)

	for _, elif := range stmt.elifList {
		elif.condition.generate(exe, currentBlock, ob)
		if_false_label = get_label(ob)
		generate_code(ob, statement.Position(), DVM_JUMP_IF_FALSE, if_false_label)

		generate_statement_list(exe, elif.block, elif.block.statementList, ob)
		generate_code(ob, statement.Position(), DVM_JUMP, end_label)
		set_label(ob, if_false_label)
	}
	if stmt.elseBlock != nil {
		generate_statement_list(exe, if_s.else_block, if_s.else_block.statement_list, ob)
	}
	set_label(ob, end_label)

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

func (s *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	s.init.fix(currentBlock)
	s.condition.fix(currentBlock)
	s.post.fix(currentBlock)
	fixStatementList(s.block, s.block.statementList, fd)

}
func (stmt *ForStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {

	if stmt.init != nil {
		stmt.generate(exe, currentBlock, ob)
	}
	loop_label = get_label(ob)
	set_label(ob, loop_label)

	if stmt.condition != nil {
		stmt.condition.generate(exe, currentBlock, ob)
	}

	parent, ok := stmt.block.parent.(*StatementBlockInfo)
	if !ok {
		compileError(stmt.Position(), 0, "")
	}

	parent.break_label = get_label(ob)
	parent.continue_label = get_label(ob)

	if stmt.condition != nil {
		generate_code(ob, stmt.Position(), DVM_JUMP_IF_FALSE, parent.break_label)
	}

	generate_statement_list(exe, stmt.block, stmt.block.statementList, ob)
	set_label(ob, parent.continue_label)

	if stmt.post != nil {
		stmt.post.generate(exe, currentBlock, ob)
	}

	generate_code(ob, statement.Position(), DVM_JUMP, loop_label)
	set_label(ob, parent.break_label)

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

func (s *ReturnStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	s.returnValue.fix(currentBlock)

	if s.returnValue == nil {
		switch fd.typeSpecifier.basicType {
		case BooleanType:
			s.returnValue = &BooleanExpression{
				booleanValue:  false,
				typeSpecifier: &TypeSpecifier{basicType: BooleanType},
			}
			return
		case NumberType:
			s.returnValue = &NumberExpression{
				numberValue:   0.0,
				typeSpecifier: &TypeSpecifier{basicType: NumberType},
			}
			return
		case StringType:
			s.returnValue = &StringExpression{
				stringValue:   "",
				typeSpecifier: &TypeSpecifier{basicType: StringType},
			}
			return
		}

	}
	// 类型转换
	createAssignCast(s.returnValue, fd.typeSpecifier)
}

func (stmt *ReturnStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.returnValue == nil {
		compileError(stmt.Position(), 0, "")
	}

	stmt.returnValue.generate(exe, currentBlock, ob)
	generate_code(ob, stmt.Position(), DVM_RETURN)

}

// ==============================
// BreakStatement
// ==============================

// BreakStatement break 语句
type BreakStatement struct {
	StatementImpl

	label string
}

func (s *BreakStatement) fix(currentBlock *Block, fd *FunctionDefinition) {}

func (stmt *BreakStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	var parent *StatementBlockInfo

	for block_p := currentBlock; block_p != nil; block_p = block_p.outerBlock {
		parent, ok := block_p.parent.(*StatementBlockInfo)
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

	if block_p == nil {
		compileError(stmt.Position(), 0, "")
	}

	generate_code(ob, statement.Position(), DVM_JUMP, parent.break_label)

}

// ==============================
// ContinueStatement
// ==============================

// ContinueStatement continue 语句
type ContinueStatement struct {
	StatementImpl

	label string
}

func (s *ContinueStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

}
func (stmt *ContinueStatement) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	var parent *StatementBlockInfo

	for block_p := currentBlock; block_p != nil; block_p = block_p.outerBlock {
		parent, ok := block_p.parent.(*StatementBlockInfo)
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

	if block_p == nil {
		dkc_compile_error(statement.Position(), 0, "")
	}
	generate_code(ob, statement.Position(), DVM_JUMP, block_p.parent.statement.continue_label)

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

	isLocal bool
}

func (s *Declaration) fix(currentBlock *Block, fd *FunctionDefinition) {
	if searchDeclaration(s.name, currentBlock) != nil {
		compileError(s.Position(), 0, "")
	}

	if currentBlock != nil {
		currentBlock.declarationList = append(currentBlock.declarationList, s)
		addLocalVariable(fd, s)
		s.isLocal = true
	} else {
		compiler := getCurrentCompiler()
		compiler.declarationList = append(compiler.declarationList, s)
		s.isLocal = false
	}

	s.initializer.fix(currentBlock)

	// 类型转换
	if s.initializer != nil {
		s.initializer = createAssignCast(s.initializer, s.typeSpecifier)
	}
}
func (stmt *Declaration) generate(exe *Executable, currentBlock *Block, ob *OpcodeBuf) {
	if stmt.initializer == nil {
		return
	}

	stmt.initializer.generate(exe, currentBlock, ob)
	generate_pop_to_identifier(decl, statement.Position(), ob)
}
