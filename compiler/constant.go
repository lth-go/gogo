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
	EqOperator:  vm.VM_EQ_INT,
	NeOperator:  vm.VM_NE_INT,
	GtOperator:  vm.VM_GT_INT,
	GeOperator:  vm.VM_GE_INT,
	LtOperator:  vm.VM_LT_INT,
	LeOperator:  vm.VM_LE_INT,
	AddOperator: vm.VM_ADD_INT,
	SubOperator: vm.VM_SUB_INT,
	MulOperator: vm.VM_MUL_INT,
	DivOperator: vm.VM_DIV_INT,
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
