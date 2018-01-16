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

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl
	// 基本类型
	basicType BasicType
	// TODO 派生类型
	//derive
}

// ==============================
// Statement 接口
// ==============================

// Statement 语句接口
type Statement interface {
	// Pos接口
	Pos

	fix(*Block, *FunctionDefinition)
}

// StatementImpl provide commonly implementations for Stmt..
type StatementImpl struct {
	PosImpl    // StmtImpl provide Pos() function.
	lineNumber int
}

// stmt provide restraint interface.
func (s *StatementImpl) stmt() {}

// Block 块接口
type Block struct {
	BlockType       int
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration
	// TODO
	// parent
}

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	typeSpecifier      *TypeSpecifier
	name               string
	parameterList      []*Parameter
	block              *Block
	localVariableCount int

	index int

	declarationList []*Declaration
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

func (s *ForStatement) fix(currentBlock *Block, fd *FunctionDefinition) {
	s.init.fix(currentBlock)
	s.condition.fix(currentBlock)
	s.post.fix(currentBlock)
	fixStatementList(s.block, s.block.statementList, fd)

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

// ==============================
// BreakStatement
// ==============================

// BreakStatement break 语句
type BreakStatement struct {
	StatementImpl
}

func (s *BreakStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

}

// ==============================
// ContinueStatement
// ==============================

// ContinueStatement continue 语句
type ContinueStatement struct {
	StatementImpl
}

func (s *ContinueStatement) fix(currentBlock *Block, fd *FunctionDefinition) {

}

// ==============================
// Declaration
// ==============================

// Declaration 声明语句
type Declaration struct {
	StatementImpl
	typeSpecifier *TypeSpecifier
	name          string
	initializer   Expression

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
