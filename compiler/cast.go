package compiler

import (
	"github.com/lth-go/gogo/vm"
)

// 声明类型转换
func CreateAssignCast(src Expression, destTye *Type) Expression {
	var castExpr Expression

	srcTye := src.GetType()

	if srcTye.Equal(destTye) {
		return src
	}

	if destTye.IsComposite() && srcTye.IsNil() {
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
	leftType := binaryExpr.left.GetType()
	rightType := binaryExpr.right.GetType()

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
	var typ *Type

	castExpr := &CastExpression{castType: castType, operand: expr}
	castExpr.SetPosition(expr.Position())

	switch castType {
	case CastTypeIntToFloat:
		typ = NewType(vm.BasicTypeFloat)
	case CastTypeFloatToInt:
		typ = NewType(vm.BasicTypeInt)
	case CastTypeBoolToString, CastTypeIntToString, CastTypeFloatToString:
		typ = NewType(vm.BasicTypeString)
	}
	castExpr.SetType(typ)

	return castExpr
}

func createToStringCast(src Expression) Expression {
	var cast Expression

	if src.GetType().IsBool() {
		cast = createCastExpression(CastTypeBoolToString, src)
	} else if src.GetType().IsInt() {
		cast = createCastExpression(CastTypeIntToString, src)
	} else if src.GetType().IsFloat() {
		cast = createCastExpression(CastTypeFloatToString, src)
	} else {
		panic("TODO")
	}

	return cast
}

func castMismatchError(pos Position, src, dest *Type) {
	srcName := src.GetTypeName()
	destName := dest.GetTypeName()
	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
