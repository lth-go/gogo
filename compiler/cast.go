package compiler

import (
	"fmt"

	"../vm"
)

func allocCastExpression(castType CastType, expr Expression) Expression {
	var typ *TypeSpecifier

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case IntToDoubleCast:
		typ = &TypeSpecifier{basicType: vm.DoubleType}
	case DoubleToIntCast:
		typ = &TypeSpecifier{basicType: vm.IntType}
	case BooleanToStringCast, IntToStringCast, DoubleToStringCast:
		typ = &TypeSpecifier{basicType: vm.StringType}
	}
	castExpr.setType(typ)

	return castExpr
}

// 声明类型转换
func createAssignCast(src Expression, destTye *TypeSpecifier) Expression {
	var castExpr Expression

	srcTye := src.typeS()

	if compareType(src.typeS(), destTye) {
		return src
	}

	if isObject(destTye) && src.typeS().basicType == vm.NullType {
		if src.typeS().deriveList != nil {
			panic("derive != NULL")
		}
		return src
	}


	if isInt(srcTye) && isDouble(destTye) {
		castExpr = allocCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if isDouble(srcTye) && isInt(destTye) {
		castExpr = allocCastExpression(DoubleToIntCast, src)
		return castExpr
	}

	castMismatchError(src.Position(), srcTye, destTye)
	return nil
}

func castBinaryExpression(expr Expression) Expression {

	binaryExpr := expr.(*BinaryExpression)

	leftType := binaryExpr.left.typeS()
	rightType := binaryExpr.right.typeS()

	if isInt(leftType) && isDouble(rightType) {
		binaryExpr.left = allocCastExpression(IntToDoubleCast, binaryExpr.left)

	} else if isDouble(leftType) && isInt(rightType) {
		binaryExpr.right = allocCastExpression(IntToDoubleCast, binaryExpr.right)

	} else if isString(leftType) && isBoolean(rightType) {
		binaryExpr.right = allocCastExpression(BooleanToStringCast, binaryExpr.right)

	} else if isString(leftType) && isInt(rightType) {
		binaryExpr.right = allocCastExpression(IntToStringCast, binaryExpr.right)

	} else if isString(leftType) && isDouble(rightType) {
		binaryExpr.right = allocCastExpression(DoubleToStringCast, binaryExpr.right)

	}
	return binaryExpr
}

func castMismatchError(pos Position, src, dest *TypeSpecifier) {
	// TODO v2
	srcName := getTypeName(src)
	destName := getTypeName(dest)

	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}

func getBasicTypeName(typ vm.BasicType) string {
	switch typ {
	case vm.BooleanType:
		return "boolean"
	case vm.IntType:
		return "int"
	case vm.DoubleType:
		return "double"
	case vm.StringType:
		return "string"
	default:
		panic(fmt.Sprintf("bad case. type..%d\n", typ))
	}
}

func getTypeName(typ *TypeSpecifier )string {
	typeName := getBasicTypeName(typ.basicType)

	for _, derive := range typ.deriveList {
        switch derive.(type) {
        case *FunctionDerive:
			panic("TODO:derive_tag")
        case *ArrayDerive:
			typeName = typeName + "[]"
        default:
			panic("TODO:derive_tag")
        }
    }

    return typeName
}
