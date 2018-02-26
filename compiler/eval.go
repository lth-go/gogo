package compiler

import (
	"strconv"

	"../vm"
)

func evalMathExpression(currentBlock *Block, binaryExpr *BinaryExpression) Expression {
	switch leftExpr := binaryExpr.left.(type) {

	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {

		case *IntExpression:
			newExpr := evalMathExpressionInt(binaryExpr, leftExpr.intValue, rightExpr.intValue)
			return newExpr

		case *DoubleExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, float64(leftExpr.intValue), rightExpr.doubleValue)
			return newExpr
		}

	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {

		case *IntExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, leftExpr.doubleValue, float64(rightExpr.intValue))
			return newExpr

		case *DoubleExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, leftExpr.doubleValue, rightExpr.doubleValue)
			return newExpr
		}

	case *StringExpression:
		if binaryExpr.operator == AddOperator {
			newExpr := chainString(binaryExpr)
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

	newExpr := &IntExpression{intValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.IntType})

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
	newExpr := &DoubleExpression{doubleValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.DoubleType})

	return newExpr
}

func chainString(binaryExpr *BinaryExpression) Expression {

	rightStr := expressionToString(binaryExpr.right)
	if rightStr == "" {
		return binaryExpr
	}

	leftStringExpr := binaryExpr.left.(*StringExpression)

	newStr := leftStringExpr.stringValue + rightStr

	newExpr := &StringExpression{stringValue: newStr}
	newExpr.setType(&TypeSpecifier{basicType: vm.StringType})

	return newExpr
}

func expressionToString(expr Expression) string {
	var newStr string

	switch e := expr.(type) {
	case *BooleanExpression:
		if e.booleanValue {
			newStr = "true"
		} else {
			newStr = "false"
		}
	case *IntExpression:
		newStr = strconv.Itoa(e.intValue)
	case *DoubleExpression:
		newStr = strconv.FormatFloat(e.doubleValue, 'f', -1, 64)
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
			newExpr := evalCompareExpressionBoolean(binaryExpr, leftExpr.booleanValue, rightExpr.booleanValue)
			return newExpr
		}

	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionInt(binaryExpr, leftExpr.intValue, rightExpr.intValue)
			return newExpr
		case *DoubleExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, float64(leftExpr.intValue), rightExpr.doubleValue)
			return newExpr
		}

	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.doubleValue, float64(rightExpr.intValue))
			return newExpr
		case *DoubleExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.doubleValue, rightExpr.doubleValue)
			return newExpr
		}

	case *StringExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *StringExpression:
			newExpr := evalCompareExpressionString(binaryExpr, leftExpr.stringValue, rightExpr.stringValue)
			return newExpr
		}
	case *NullExpression:
		switch binaryExpr.right.(type) {
		case *NullExpression:
			newExpr := &BooleanExpression{booleanValue: true}
			newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})
			return newExpr
		}
	}

	return binaryExpr
}

func evalCompareExpressionBoolean(binaryExpr *BinaryExpression, left, right bool) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})

	return newExpr
}

func evalCompareExpressionInt(binaryExpr *BinaryExpression, left, right int) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	case GtOperator:
		value = (left > right)
	case GeOperator:
		value = (left >= right)
	case LtOperator:
		value = (left < right)
	case LeOperator:
		value = (left <= right)
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})
	return newExpr
}

func evalCompareExpressionDouble(binaryExpr *BinaryExpression, left, right float64) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	case GtOperator:
		value = (left > right)
	case GeOperator:
		value = (left >= right)
	case LtOperator:
		value = (left < right)
	case LeOperator:
		value = (left <= right)
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})
	return newExpr
}

func evalCompareExpressionString(binaryExpr *BinaryExpression, left, right string) Expression {
	var value bool

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	case GtOperator:
		value = (left > right)
	case GeOperator:
		value = (left >= right)
	case LtOperator:
		value = (left < right)
	case LeOperator:
		value = (left <= right)
	default:
		compileError(binaryExpr.Position(), COMPARE_TYPE_MISMATCH_ERR)
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})

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
	newBinaryExpr := castBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.typeS()
	newBinaryExprRightType := newBinaryExpr.right.typeS()

	if isInt(newBinaryExprLeftType) && isInt(newBinaryExprRightType) {
		newBinaryExpr.setType(&TypeSpecifier{basicType: vm.IntType})

	} else if isDouble(newBinaryExprLeftType) && isDouble(newBinaryExprRightType) {
		newBinaryExpr.setType(&TypeSpecifier{basicType: vm.DoubleType})

	} else if expr.operator == AddOperator {
		if (isString(newBinaryExprLeftType) && isString(newBinaryExprRightType)) ||
			(isString(newBinaryExprLeftType) && isNull(newBinaryExpr.left)) {
			newBinaryExpr.setType(&TypeSpecifier{basicType: vm.StringType})
		}
	} else {
		compileError(expr.Position(), MATH_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(newBinaryExprLeftType.basicType), int(newBinaryExprRightType.basicType))
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

	newBinaryExpr := castBinaryExpression(expr)

	newBinaryExprLeftType := newBinaryExpr.left.typeS()
	newBinaryExprRightType := newBinaryExpr.right.typeS()

	// TODO 字符串是否能跟null比较
	if !(compareType(newBinaryExprLeftType, newBinaryExprRightType) ||
		(isObject(newBinaryExprLeftType) && isNull(newBinaryExpr.right) ||
			(isNull(newBinaryExpr.left) && isObject(newBinaryExprRightType)))) {
		compileError(expr.Position(), COMPARE_TYPE_MISMATCH_ERR, getTypeName(newBinaryExprLeftType), getTypeName(newBinaryExprRightType))
	}

	newBinaryExpr.setType(&TypeSpecifier{basicType: vm.BooleanType})

	return newBinaryExpr
}

func fixLogicalBinaryExpression(expr *BinaryExpression, currentBlock *Block) Expression {
	expr.left = expr.left.fix(currentBlock)
	expr.right = expr.right.fix(currentBlock)

	if isBoolean(expr.left.typeS()) && isBoolean(expr.right.typeS()) {
		expr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
		expr.typeS().fix(0)
		return expr
	}

	compileError(expr.Position(), LOGICAL_TYPE_MISMATCH_ERR, "Left: %d, Right: %d\n", int(expr.left.typeS().basicType), int(expr.right.typeS().basicType))
	return nil
}
