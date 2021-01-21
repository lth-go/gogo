package compiler

import (
	"strconv"

	"github.com/lth-go/gogo/vm"
)

func EvalMathExpression(binaryExpr *BinaryExpression) Expression {
	switch leftExpr := binaryExpr.left.(type) {
	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionInt(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr

		case *FloatExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, float64(leftExpr.Value), rightExpr.Value)
			return newExpr
		}
	case *FloatExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, leftExpr.Value, float64(rightExpr.Value))
			return newExpr
		case *FloatExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}
	case *StringExpression:
		if binaryExpr.operator == AddOperator {
			newExpr := chainBinaryExpressionString(binaryExpr)
			return newExpr
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

	newExpr := &IntExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeInt))

	return newExpr
}

func evalMathExpressionDouble(binaryExpr *BinaryExpression, left, right float64) Expression {
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
	newExpr := &FloatExpression{Value: value}
	newExpr.SetType(NewType(vm.BasicTypeFloat))

	return newExpr
}

func chainBinaryExpressionString(binaryExpr *BinaryExpression) Expression {
	rightStr := expressionToString(binaryExpr.right)
	if rightStr == "" {
		return binaryExpr
	}

	leftStringExpr := binaryExpr.left.(*StringExpression)

	newStr := leftStringExpr.Value + rightStr

	newExpr := &StringExpression{Value: newStr}
	newExpr.SetType(NewType(vm.BasicTypeString))

	return newExpr
}

// TODO: 只能字符串相加
func expressionToString(expr Expression) string {
	var newStr string

	switch e := expr.(type) {
	case *BoolExpression:
		if e.Value {
			newStr = "true"
		} else {
			newStr = "false"
		}
	case *IntExpression:
		newStr = strconv.Itoa(e.Value)
	case *FloatExpression:
		newStr = strconv.FormatFloat(e.Value, 'f', -1, 64)
	case *StringExpression:
		newStr = e.Value
	default:
		newStr = ""
	}

	return newStr
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

func FixMathBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	// 能否合并计算
	newExpr := EvalMathExpression(expr)
	switch newExpr.(type) {
	case *IntExpression, *FloatExpression, *StringExpression:
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

	} else if expr.operator == AddOperator {
		if (newBinaryExprLeftType.IsString() && newBinaryExprRightType.IsString()) ||
			(newBinaryExprLeftType.IsString() && isNilExpression(newBinaryExpr.left)) {
			newBinaryExpr.SetType(NewType(vm.BasicTypeString))
		}
	} else {
		compileError(
			expr.Position(),
			MATH_TYPE_MISMATCH_ERR,
			"Left: %d, Right: %d\n", newBinaryExprLeftType.GetBasicType(), newBinaryExprRightType.GetBasicType(),
		)
	}

	return newBinaryExpr
}

func FixCompareBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	newExpr := evalCompareExpression(expr)
	switch newExpr.(type) {
	case *BoolExpression:
		return newExpr
	}

	newBinaryExpr := CastBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.GetType()
	newBinaryExprRightType := newBinaryExpr.right.GetType()

	if !compareType(newBinaryExprLeftType, newBinaryExprRightType) {
		if !(newBinaryExprLeftType.IsComposite() &&
			isNilExpression(newBinaryExpr.right) ||
			(isNilExpression(newBinaryExpr.left) &&
				newBinaryExprRightType.IsComposite())) {

			compileError(expr.Position(), COMPARE_TYPE_MISMATCH_ERR, newBinaryExprLeftType.GetTypeName(), newBinaryExprRightType.GetTypeName())
		}
	}

	newBinaryExpr.SetType(NewType(vm.BasicTypeBool))

	return newBinaryExpr
}

func FixLogicalBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	if expr.left.GetType().IsBool() && expr.right.GetType().IsBool() {
		expr.Type = NewType(vm.BasicTypeBool)
		expr.GetType().Fix()
		return expr
	}

	compileError(
		expr.Position(),
		LOGICAL_TYPE_MISMATCH_ERR,
		"Left: %d, Right: %d\n", expr.left.GetType().GetBasicType(), expr.right.GetType().GetBasicType(),
	)
	return nil
}
