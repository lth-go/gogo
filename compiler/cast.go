package compiler

import (
	"../vm"
	"fmt"
)

// 声明类型转换
func createAssignCast(src Expression, destTye *TypeSpecifier) Expression {
	var castExpr Expression

	srcTye := src.typeS()

	if srcTye.deriveList != nil || destTye.deriveList != nil {
		compileError(src.Position(), DERIVE_TYPE_CAST_ERR)
	}

	if srcTye.basicType == destTye.basicType {
		return src
	}

	if isInt(srcTye) && isDouble(destTye) {
		castExpr = allocCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if isDouble(srcTye) && isInt(destTye) {
		castExpr = allocCastExpression(DoubleToIntCast, src)
		return castExpr
	}

	castMismatchError(src.Position(), srcTye.basicType, destTye.basicType)
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

func castMismatchError(pos Position, src, dest vm.BasicType) {
	src_name := getBasicTypeName(src)
	dest_name := getBasicTypeName(dest)

	compileError(pos, CAST_MISMATCH_ERR, src_name, dest_name)
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
