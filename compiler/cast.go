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
		typ = newTypeSpecifier(vm.BasicTypeFloat)
	case DoubleToIntCast:
		typ = newTypeSpecifier(vm.BasicTypeInt)
	case BooleanToStringCast, IntToStringCast, DoubleToStringCast:
		typ = newTypeSpecifier(vm.BasicTypeString)
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

	if destTye.IsObject() && srcTye.IsNil() {
		return src
	}

	if srcTye.IsInt() && destTye.IsFloat() {
		castExpr = createCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if srcTye.IsFloat() && destTye.IsInt() {
		castExpr = createCastExpression(DoubleToIntCast, src)
		return castExpr

	} else if destTye.IsString() {
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

	if src.typeS().IsBool() {
		cast = createCastExpression(BooleanToStringCast, src)
	} else if src.typeS().IsInt() {
		cast = createCastExpression(IntToStringCast, src)
	} else if src.typeS().IsFloat() {
		cast = createCastExpression(DoubleToStringCast, src)
	} else {
		panic("TODO")
	}

	return cast
}

func castBinaryExpression(binaryExpr *BinaryExpression) *BinaryExpression {

	leftType := binaryExpr.left.typeS()
	rightType := binaryExpr.right.typeS()

	if leftType.IsInt() && rightType.IsFloat() {
		binaryExpr.left = createCastExpression(IntToDoubleCast, binaryExpr.left)

	} else if leftType.IsFloat() && rightType.IsInt() {
		binaryExpr.right = createCastExpression(IntToDoubleCast, binaryExpr.right)

	} else if leftType.IsString() && rightType.IsBool() {
		binaryExpr.right = createCastExpression(BooleanToStringCast, binaryExpr.right)

	} else if leftType.IsString() && rightType.IsInt() {
		binaryExpr.right = createCastExpression(IntToStringCast, binaryExpr.right)

	} else if leftType.IsString() && rightType.IsFloat() {
		binaryExpr.right = createCastExpression(DoubleToStringCast, binaryExpr.right)
	}

	return binaryExpr
}

func castMismatchError(pos Position, src, dest *TypeSpecifier) {
	srcName := src.GetTypeName()
	destName := src.GetTypeName()

	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
