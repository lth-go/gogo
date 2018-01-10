package parser

// 基础类型
type BasicType int

const (
	BOOLEAN_TYPE BasicType = iota
	INT_TYPE
	DOUBLE_TYPE
	STRING_TYPE
)

// 类型
type TypeSpecifier struct {
	// 基本类型
	basicType BasicType
	// 派生类型
	//derive
}

// 变量声明
type Declaration struct {
	name string
	Type TypeSpecifier

	initializer   Expression
	variableIndex int
	isLocal       Boolean
}

type Statement interface {
	Pos
}

type Block struct {
	BlockType        int
	outerBlock      *Block
	statementList   []Statement
	declarationList []*Declaration
	// TODO
	// parent
}

type Parameter struct {
	Type       *TypeSpecifier
	name       string
	lineNumber int
}

type AssignmentOperator int

const (
	NORMAL_ASSIGN AssignmentOperator = iota
	ADD_ASSIGN
	SUB_ASSIGN
	MUL_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
)

type ExpressionStatement struct {
	expression_s Expression
}

type IfStatement struct {
	condition  Expression
	then_block *Block
	elsif_list []*Elif
	else_block *Block
}

type Elif struct {
	condition Expression
	block     *Block
}

type WhildStatement struct {
	condition Expression
	block     *Block
}

type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *Block
}

// return 语句
type ReturnStatement struct {
	// 返回值
	returnValue Expression
}

// break 语句
type BreakStatement struct{}

// continue 语句
type ContinueStatement struct{}

// 声明语句
type DeclarationStatement struct {
	Type        *TypeSpecifier
	name        string
	initializer Expression
	lineNumber  int
}

type FunctionDefinition struct {
	Type               *TypeSpecifier
	name               string
	parameter          *Parameter
	block              *Block
	localVariableCount int

	index int

	declarationList []*Declaration
}
