package compiler

import (
	"../vm"
)

//
// cast func
//

// 声明类型转换
func createAssignCast(src Expression, dest *TypeSpecifier) Expression {
	var castExpr Expression
	if src.typeS().deriveList != nil || dest.deriveList != nil {
		compileError(src.Position(), 0, "")
	}

	if src.typeS().basicType == dest.basicType {
		return src
	}

	if src.typeS().basicType == vm.IntType && dest.basicType == vm.DoubleType {
		castExpr = allocCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if src.typeS().basicType == vm.DoubleType && dest.basicType == vm.IntType {
		castExpr = allocCastExpression(DoubleToIntCast, src)
		return castExpr
	}

	compileError(src.Position(), 0, "")
	return nil
}

func castBinaryExpression(expr Expression) Expression {

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

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
