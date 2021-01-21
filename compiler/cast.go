package compiler

import (
	"github.com/lth-go/gogo/vm"
)

// 声明类型转换
func CreateAssignCast(src Expression, destType *Type) Expression {
	srcTye := src.GetType()

	if srcTye.Equal(destType) {
		return src
	}

	if destType.IsComposite() && srcTye.IsNil() {
		return src
	}

	if destType.IsFloat() {
		expr, ok := src.(*IntExpression)
		if ok {
			return CreateFloatExpression(expr.Position(), float64(expr.Value))
		}
	}

	if destType.IsInt() {
		expr, ok := src.(*FloatExpression)
		if ok {
			return CreateIntExpression(expr.Position(), int(expr.Value))
		}
	}

	castMismatchError(src.Position(), srcTye, destType)
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

func castMismatchError(pos Position, src, dest *Type) {
	srcName := src.GetTypeName()
	destName := dest.GetTypeName()
	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
