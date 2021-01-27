package compiler

import (
	"github.com/lth-go/gogo/vm"
)

//
// BinaryOperatorKind
//
type BinaryOperatorKind int

const (
	LogicalOrOperator BinaryOperatorKind = iota
	LogicalAndOperator
	EqOperator
	NeOperator
	GtOperator
	GeOperator
	LtOperator
	LeOperator
	AddOperator
	SubOperator
	MulOperator
	DivOperator
)

var operatorCodeMap = map[BinaryOperatorKind]byte{
	EqOperator:  vm.OP_CODE_EQ_INT,
	NeOperator:  vm.OP_CODE_NE_INT,
	GtOperator:  vm.OP_CODE_GT_INT,
	GeOperator:  vm.OP_CODE_GE_INT,
	LtOperator:  vm.OP_CODE_LT_INT,
	LeOperator:  vm.OP_CODE_LE_INT,
	AddOperator: vm.OP_CODE_ADD_INT,
	SubOperator: vm.OP_CODE_SUB_INT,
	MulOperator: vm.OP_CODE_MUL_INT,
	DivOperator: vm.OP_CODE_DIV_INT,
}

type UnaryOperatorKind int

const (
	UnaryOperatorKindMinus UnaryOperatorKind = iota
	UnaryOperatorKindNot
)

//
// Cast
//
type CastType int

const (
	CastTypeIntToString CastType = iota
	CastTypeBoolToString
	CastTypeFloatToString
	CastTypeIntToFloat
	CastTypeFloatToInt
)
