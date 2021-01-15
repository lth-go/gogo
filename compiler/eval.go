package compiler

import (
	"strconv"

	"github.com/lth-go/gogo/vm"
)

func evalMathExpression(currentBlock *Block, binaryExpr *BinaryExpression) Expression {
	switch leftExpr := binaryExpr.left.(type) {
	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionInt(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr

		case *DoubleExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, float64(leftExpr.Value), rightExpr.Value)
			return newExpr
		}
	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, leftExpr.Value, float64(rightExpr.Value))
			return newExpr
		case *DoubleExpression:
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
	newExpr.setType(NewType(vm.BasicTypeInt))

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
	newExpr := &DoubleExpression{Value: value}
	newExpr.setType(NewType(vm.BasicTypeFloat))

	return newExpr
}

func chainBinaryExpressionString(binaryExpr *BinaryExpression) Expression {
	rightStr := expressionToString(binaryExpr.right)
	if rightStr == "" {
		return binaryExpr
	}

	leftStringExpr := binaryExpr.left.(*StringExpression)

	newStr := leftStringExpr.stringValue + rightStr

	newExpr := &StringExpression{stringValue: newStr}
	newExpr.setType(NewType(vm.BasicTypeString))

	return newExpr
}

func expressionToString(expr Expression) string {
	var newStr string

	switch e := expr.(type) {
	case *BooleanExpression:
		if e.Value {
			newStr = "true"
		} else {
			newStr = "false"
		}
	case *IntExpression:
		newStr = strconv.Itoa(e.Value)
	case *DoubleExpression:
		newStr = strconv.FormatFloat(e.Value, 'f', -1, 64)
	case *StringExpression:
		newStr = e.stringValue
	default:
		newStr = ""
	}

	return newStr
}

func evalCompareExpression(binaryExpr *BinaryExpression) Expression {

	switch leftExpr := binaryExpr.left.(type) {

	case *BooleanExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *BooleanExpression:
			newExpr := evalCompareExpressionBoolean(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}

	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionInt(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		case *DoubleExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, float64(leftExpr.Value), rightExpr.Value)
			return newExpr
		}

	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.Value, float64(rightExpr.Value))
			return newExpr
		case *DoubleExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.Value, rightExpr.Value)
			return newExpr
		}

	case *StringExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *StringExpression:
			newExpr := evalCompareExpressionString(binaryExpr, leftExpr.stringValue, rightExpr.stringValue)
			return newExpr
		}
	case *NilExpression:
		switch binaryExpr.right.(type) {
		case *NilExpression:
			newExpr := &BooleanExpression{Value: true}
			newExpr.setType(NewType(vm.BasicTypeBool))
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

	newExpr := &BooleanExpression{Value: value}
	newExpr.setType(NewType(vm.BasicTypeBool))

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

	newExpr := &BooleanExpression{Value: value}
	newExpr.setType(NewType(vm.BasicTypeBool))
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

	newExpr := &BooleanExpression{Value: value}
	newExpr.setType(NewType(vm.BasicTypeBool))
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

	newExpr := &BooleanExpression{Value: value}
	newExpr.setType(NewType(vm.BasicTypeBool))

	return newExpr
}

func fixMathBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	// 能否合并计算
	newExpr := evalMathExpression(currentBlock, expr)
	switch newExpr.(type) {
	case *IntExpression, *DoubleExpression, *StringExpression:
		return newExpr
	}

	// 类型转换
	newBinaryExpr := CastBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.typeS()
	newBinaryExprRightType := newBinaryExpr.right.typeS()

	if newBinaryExprLeftType.IsInt() && newBinaryExprRightType.IsInt() {
		newBinaryExpr.setType(NewType(vm.BasicTypeInt))

	} else if newBinaryExprLeftType.IsFloat() && newBinaryExprRightType.IsFloat() {
		newBinaryExpr.setType(NewType(vm.BasicTypeFloat))

	} else if expr.operator == AddOperator {
		if (newBinaryExprLeftType.IsString() && newBinaryExprRightType.IsString()) ||
			(newBinaryExprLeftType.IsString() && isNull(newBinaryExpr.left)) {
			newBinaryExpr.setType(NewType(vm.BasicTypeString))
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

func fixCompareBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	newExpr := evalCompareExpression(expr)
	switch newExpr.(type) {
	case *BooleanExpression:
		return newExpr
	}

	newBinaryExpr := CastBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.typeS()
	newBinaryExprRightType := newBinaryExpr.right.typeS()

	// TODO 字符串是否能跟null比较
	if !(compareType(newBinaryExprLeftType, newBinaryExprRightType) ||
		(newBinaryExprLeftType.IsObject() &&
			isNull(newBinaryExpr.right) ||
			(isNull(newBinaryExpr.left) && newBinaryExprRightType.IsObject()))) {
		compileError(expr.Position(), COMPARE_TYPE_MISMATCH_ERR, newBinaryExprLeftType.GetTypeName(), newBinaryExprRightType.GetTypeName())
	}

	newBinaryExpr.setType(NewType(vm.BasicTypeBool))

	return newBinaryExpr
}

func fixLogicalBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	if expr.left.typeS().IsBool() && expr.right.typeS().IsBool() {
		expr.typeSpecifier = NewType(vm.BasicTypeBool)
		expr.typeS().fix()
		return expr
	}

	compileError(
		expr.Position(),
		LOGICAL_TYPE_MISMATCH_ERR,
		"Left: %d, Right: %d\n", expr.left.typeS().GetBasicType(), expr.right.typeS().GetBasicType(),
	)
	return nil
}
