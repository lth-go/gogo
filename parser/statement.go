package parser

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
	// 基本类型
	basicType BasicType
	// TODO 派生类型
	//derive
}

// Declaration 变量声明
type Declaration struct {
	name          string
	typeSpecifier *TypeSpecifier

	initializer   Expression
	variableIndex int
	isLocal       Boolean
}

// Statement 语句接口
type Statement interface {
	// Pos
}

// Block 块接口
type Block struct {
	BlockType       int
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration
	// TODO
	// parent
}

// Parameter 形参
type Parameter struct {
	typeSpecifier *TypeSpecifier
	name          string
	lineNumber    int
}

// AssignmentOperator ...
type AssignmentOperator int

const (
	// NormalAssign 赋值操作符 =
	NormalAssign AssignmentOperator = iota
)

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	expression Expression
}

// IfStatement if表达式
type IfStatement struct {
	condition Expression
	thenBlock *Block
	elifList  []*Elif
	elseBlock *Block
}

// Elif ...
type Elif struct {
	condition Expression
	block     *Block
}

// ForStatement for语句
type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *Block
}

// ReturnStatement return 语句
type ReturnStatement struct {
	// 返回值
	returnValue Expression
}

// BreakStatement break 语句
type BreakStatement struct{}

// ContinueStatement continue 语句
type ContinueStatement struct{}

// DeclarationStatement 声明语句
type DeclarationStatement struct {
	typeSpecifier *TypeSpecifier
	name          string
	initializer   Expression
	lineNumber    int
}

// FunctionDefinition 函数定义
type FunctionDefinition struct {
	typeSpecifier      *TypeSpecifier
	name               string
	parameter          *Parameter
	block              *Block
	localVariableCount int

	index int

	declarationList []*Declaration
}
