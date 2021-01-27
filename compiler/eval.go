package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func FixMathBinaryExpression(expr *BinaryExpression) Expression {
	expr.left = expr.left.Fix()
	expr.right = expr.right.Fix()

	// 能否合并计算
	newExpr := evalMathExpression(expr)
	switch newExpr.(type) {
	case *IntExpression, *FloatExpression, *StringExpression:
		newExpr = newExpr.Fix()
		return newExpr
	}

	// 类型转换
	newBinaryExpr := CastBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.GetType()
	newBinaryExprRightType := newBinaryExpr.right.GetType()

	if newBinaryExprLeftType.IsInt() && newBinaryExprRightType.IsInt() {
		newBinaryExpr.SetType(NewType(vm.BasicTypeInt))
	} else if newBinaryExprLeftType.IsFloat() && newBinaryExprRightType.IsFloat() {
		newBinaryExpr.SetType(NewType(vm.BasicTypeFloat))
	} else if expr.operator == AddOperator && newBinaryExprLeftType.IsString() && newBinaryExprRightType.IsString() {
		newBinaryExpr.SetType(NewType(vm.BasicTypeString))
	} else {
		compileError(
			expr.Position(),
			MATH_TYPE_MISMATCH_ERR,
			"Left: %d, Right: %d\n",
			newBinaryExprLeftType.GetBasicType(),
			newBinaryExprRightType.GetBasicType(),
		)
	}

	return newBinaryExpr
}

func FixCompareBinaryExpression(expr *BinaryExpression) Expression {
	expr.left = expr.left.Fix()
	expr.right = expr.right.Fix()

	newExpr := evalCompareExpression(expr)
	switch newExpr.(type) {
	case *BoolExpression:
		return newExpr
	}

	newBinaryExpr := CastBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.GetType()
	newBinaryExprRightType := newBinaryExpr.right.GetType()

	ok := func() bool {
		if newBinaryExprLeftType.Equal(newBinaryExprRightType) {
			return true
		}

		if newBinaryExprLeftType.IsComposite() && isNilExpression(newBinaryExpr.right) {
			return true
		}

		if isNilExpression(newBinaryExpr.left) && newBinaryExprRightType.IsComposite() {
			return true
		}

		return false
	}()
	if !ok {
		compileError(expr.Position(),
			COMPARE_TYPE_MISMATCH_ERR,
			newBinaryExprLeftType.GetTypeName(),
			newBinaryExprRightType.GetTypeName(),
		)
	}

	newBinaryExpr.SetType(NewType(vm.BasicTypeBool))

	return newBinaryExpr
}

func FixLogicalBinaryExpression(expr *BinaryExpression) Expression {
	expr.left = expr.left.Fix()
	expr.right = expr.right.Fix()

	if expr.left.GetType().IsBool() && expr.right.GetType().IsBool() {
		expr.Type = NewType(vm.BasicTypeBool)
		expr.GetType().Fix()
		return expr
	}

	compileError(
		expr.Position(),
		LOGICAL_TYPE_MISMATCH_ERR,
		"Left: %d, Right: %d\n",
		expr.left.GetType().GetBasicType(),
		expr.right.GetType().GetBasicType(),
	)
	return nil
}

func evalMathExpression(binaryExpr *BinaryExpression) Expression {
	switch leftExpr := binaryExpr.left.(type) {
	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionInt(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		case *FloatExpression:
			newExpr := evalMathExpressionInt(binaryExpr, leftExpr.Value, int(rightExpr.Value))
			return newExpr
		}
	case *FloatExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionFloat(binaryExpr, leftExpr.Value, float64(rightExpr.Value))
			return newExpr
		case *FloatExpression:
			newExpr := evalMathExpressionFloat(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}
	case *StringExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *StringExpression:
			if binaryExpr.operator == AddOperator {
				newExpr := &StringExpression{Value: leftExpr.Value + rightExpr.Value}
				newExpr.Fix()
				return newExpr
			}
		}
	}

	return binaryExpr
}

func evalMathExpressionInt(binaryExpr *BinaryExpression, left, right int) Expression {
	var value int

	switch binaryExpr.operator {
	case AddOperator:
		value = left + right
	case SubOperator:
		value = left - right
	case MulOperator:
		value = left * right
	case DivOperator:
		if right == 0 {
			compileError(binaryExpr.Position(), DIVISION_BY_ZERO_IN_COMPILE_ERR)
		}
		value = left / right
	default:
		compileError(binaryExpr.Position(), MATH_TYPE_MISMATCH_ERR)
	}

	newExpr := CreateIntExpression(binaryExpr.Position(), value)
	newExpr.SetType(NewType(vm.BasicTypeInt))

	return newExpr
}

func evalMathExpressionFloat(binaryExpr *BinaryExpression, left, right float64) Expression {
	var value float64

	switch binaryExpr.operator {
	case AddOperator:
		value = left + right
	case SubOperator:
		value = left - right
	case MulOperator:
		value = left * right
	case DivOperator:
		if right == 0.0 {
			compileError(binaryExpr.Position(), DIVISION_BY_ZERO_IN_COMPILE_ERR)
		}
		value = left / right
	default:
		compileError(binaryExpr.Position(), MATH_TYPE_MISMATCH_ERR)
	}
	newExpr := CreateFloatExpression(binaryExpr.Position(), value)
	newExpr.SetType(NewType(vm.BasicTypeFloat))

	return newExpr
}

func evalCompareExpression(binaryExpr *BinaryExpression) Expression {
	switch leftExpr := binaryExpr.left.(type) {

	case *BoolExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *BoolExpression:
			newExpr := evalCompareExpressionBoolean(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}
	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionInt(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		case *FloatExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, float64(leftExpr.Value), rightExpr.Value)
			return newExpr
		}
	case *FloatExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.Value, float64(rightExpr.Value))
			return newExpr
		case *FloatExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}
	case *StringExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *StringExpression:
			newExpr := evalCompareExpressionString(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}
	case *NilExpression:
		switch binaryExpr.right.(type) {
		case *NilExpression:
			newExpr := &BoolExpression{Value: true}
			newExpr.SetType(NewType(vm.BasicTypeBool))
			return newExpr
		}
	}

	return binaryExpr
}

func evalCompareExpressionBoolean(binaryExpr *BinaryExpression, left, right bool) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = left == right
	case NeOperator:
		value = left != right
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BoolExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeBool))

	return newExpr
}

func evalCompareExpressionInt(binaryExpr *BinaryExpression, left, right int) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = left == right
	case NeOperator:
		value = left != right
	case GtOperator:
		value = left > right
	case GeOperator:
		value = left >= right
	case LtOperator:
		value = left < right
	case LeOperator:
		value = left <= right
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BoolExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeBool))
	return newExpr
}

func evalCompareExpressionDouble(binaryExpr *BinaryExpression, left, right float64) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = left == right
	case NeOperator:
		value = left != right
	case GtOperator:
		value = left > right
	case GeOperator:
		value = left >= right
	case LtOperator:
		value = left < right
	case LeOperator:
		value = left <= right
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BoolExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeBool))
	return newExpr
}

func evalCompareExpressionString(binaryExpr *BinaryExpression, left, right string) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = left == right
	case NeOperator:
		value = left != right
	case GtOperator:
		value = left > right
	case GeOperator:
		value = left >= right
	case LtOperator:
		value = left < right
	case LeOperator:
		value = left <= right
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BoolExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeBool))

	return newExpr
}

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
			newExpr := CreateFloatExpression(expr.Position(), float64(expr.Value))
			newExpr.Fix()
			return newExpr
		}
	}

	if destType.IsInt() {
		expr, ok := src.(*FloatExpression)
		if ok {
			newExpr := CreateIntExpression(expr.Position(), int(expr.Value))
			newExpr.Fix()
			return newExpr
		}
	}

	castMismatchError(src.Position(), srcTye, destType)
	return nil
}

func CastBinaryExpression(binaryExpr *BinaryExpression) *BinaryExpression {
	leftType := binaryExpr.left.GetType()
	rightType := binaryExpr.right.GetType()

	if leftType.IsInt() && rightType.IsFloat() {
		right, ok := binaryExpr.right.(*FloatExpression)
		if !ok {
			compileError(binaryExpr.Position(), CAST_MISMATCH_ERR, leftType.GetBasicType(), rightType.GetBasicType())
		}

		binaryExpr.right = CreateIntExpression(right.Position(), int(right.Value))
		binaryExpr.right.Fix()
	} else if leftType.IsFloat() && rightType.IsInt() {
		right, ok := binaryExpr.right.(*IntExpression)
		if !ok {
			compileError(binaryExpr.Position(), CAST_MISMATCH_ERR, leftType.GetBasicType(), rightType.GetBasicType())
		}

		binaryExpr.right = CreateFloatExpression(right.Position(), float64(right.Value))
		binaryExpr.right.Fix()
	}

	return binaryExpr
}

func castMismatchError(pos Position, src, dest *Type) {
	srcName := src.GetTypeName()
	destName := dest.GetTypeName()
	compileError(pos, CAST_MISMATCH_ERR, srcName, destName)
}
