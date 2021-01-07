package compiler

import (
	"github.com/lth-go/gogogogo/vm"
)

func createCastExpression(castType CastType, expr Expression) Expression {
	var typ *TypeSpecifier

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case IntToDoubleCast:
		typ = newTypeSpecifier(vm.DoubleType)
	case DoubleToIntCast:
		typ = newTypeSpecifier(vm.IntType)
	case BooleanToStringCast, IntToStringCast, DoubleToStringCast:
		typ = newTypeSpecifier(vm.StringType)
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

	if isObject(destTye) && srcTye.basicType == vm.NullType {
		return src
	}

	if isInt(srcTye) && isDouble(destTye) {
		castExpr = createCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if isDouble(srcTye) && isInt(destTye) {
		castExpr = createCastExpression(DoubleToIntCast, src)
		return castExpr

	} else if isString(destTye) {
		castExpr = createToStringCast(src)
		if castExpr != nil {
			return castExpr
		}
	}

	castMismatchError(src.Position(), srcTye, destTye)
	return nil
}

func createToStringCast(src Expression) Expression {
	var cast Expression

	if isBoolean(src.typeS()) {
		cast = createCastExpression(BooleanToStringCast, src)
	} else if isInt(src.typeS()) {
		cast = createCastExpression(IntToStringCast, src)
	} else if isDouble(src.typeS()) {
		cast = createCastExpression(DoubleToStringCast, src)
	} else {
		panic("TODO")
	}

	return cast
}

func castBinaryExpression(binaryExpr *BinaryExpression) *BinaryExpression {

	leftType := binaryExpr.left.typeS()
	rightType := binaryExpr.right.typeS()

	if isInt(leftType) && isDouble(rightType) {
		binaryExpr.left = createCastExpression(IntToDoubleCast, binaryExpr.left)

	} else if isDouble(leftType) && isInt(rightType) {
		binaryExpr.right = createCastExpression(IntToDoubleCast, binaryExpr.right)

	} else if isString(leftType) && isBoolean(rightType) {
		binaryExpr.right = createCastExpression(BooleanToStringCast, binaryExpr.right)

	} else if isString(leftType) && isInt(rightType) {
		binaryExpr.right = createCastExpression(IntToStringCast, binaryExpr.right)

	} else if isString(leftType) && isDouble(rightType) {
		binaryExpr.right = createCastExpression(DoubleToStringCast, binaryExpr.right)
	}

	return binaryExpr
}

func castMismatchError(pos Position, src, dest *TypeSpecifier) {
	srcName := src.GetTypeName()
	destName := src.GetTypeName()

	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
