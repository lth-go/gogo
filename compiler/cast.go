package compiler

import (
	"github.com/lth-go/gogogogo/vm"
)

// 声明类型转换
func CreateAssignCast(src Expression, destTye *TypeSpecifier) Expression {
	var castExpr Expression

	srcTye := src.typeS()

	if compareType(src.typeS(), destTye) {
		return src
	}

	if destTye.IsObject() && srcTye.IsNil() {
		return src
	}

	if srcTye.IsInt() && destTye.IsFloat() {
		castExpr = createCastExpression(CastTypeIntToFloat, src)
		return castExpr
	} else if srcTye.IsFloat() && destTye.IsInt() {
		castExpr = createCastExpression(CastTypeFloatToInt, src)
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

func CastBinaryExpression(binaryExpr *BinaryExpression) *BinaryExpression {
	leftType := binaryExpr.left.typeS()
	rightType := binaryExpr.right.typeS()

	if leftType.IsInt() && rightType.IsFloat() {
		binaryExpr.left = createCastExpression(CastTypeIntToFloat, binaryExpr.left)
	} else if leftType.IsFloat() && rightType.IsInt() {
		binaryExpr.right = createCastExpression(CastTypeIntToFloat, binaryExpr.right)
	} else if leftType.IsString() && rightType.IsBool() {
		binaryExpr.right = createCastExpression(CastTypeBoolToString, binaryExpr.right)
	} else if leftType.IsString() && rightType.IsInt() {
		binaryExpr.right = createCastExpression(CastTypeIntToString, binaryExpr.right)
	} else if leftType.IsString() && rightType.IsFloat() {
		binaryExpr.right = createCastExpression(CastTypeFloatToString, binaryExpr.right)
	}

	return binaryExpr
}

func createCastExpression(castType CastType, expr Expression) Expression {
	var typ *TypeSpecifier

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case CastTypeIntToFloat:
		typ = newTypeSpecifier(vm.BasicTypeFloat)
	case CastTypeFloatToInt:
		typ = newTypeSpecifier(vm.BasicTypeInt)
	case CastTypeBoolToString, CastTypeIntToString, CastTypeFloatToString:
		typ = newTypeSpecifier(vm.BasicTypeString)
	}
	castExpr.setType(typ)

	return castExpr
}

func createToStringCast(src Expression) Expression {
	var cast Expression

	if src.typeS().IsBool() {
		cast = createCastExpression(CastTypeBoolToString, src)
	} else if src.typeS().IsInt() {
		cast = createCastExpression(CastTypeIntToString, src)
	} else if src.typeS().IsFloat() {
		cast = createCastExpression(CastTypeFloatToString, src)
	} else {
		panic("TODO")
	}

	return cast
}

func castMismatchError(pos Position, src, dest *TypeSpecifier) {
	srcName := src.GetTypeName()
	destName := src.GetTypeName()
	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
