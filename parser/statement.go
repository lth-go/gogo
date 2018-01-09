package parser

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
	basic_type BasicType
	// 派生类型
	//derive
}

type Stmt interface {
	Pos
}

type Block struct {
	BlockType        int
	outer_block      *Block
	statement_list   []Statement
	declaration_list []Declaration
	// TODO
	// parent
}

type Parameter struct {
	Type        *TypeSpecifier
	name        string
	line_number int
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
	expression_s *Expression
}

type IfStatement struct {
	condition  Expression
	then_block *Block
	elsif_list []Elif
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

type ReturnStatement struct {
	// 返回值
	return_value Expression
}

type BreakStatement struct{}

type ContinueStatement struct{}

type DeclarationStatement struct {
	Type        *TypeSpecifier
	name        string
	initializer Expression
	line_number int
}
