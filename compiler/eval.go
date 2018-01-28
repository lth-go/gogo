package compiler

import "strconv"

func evalMathExpression(currentBlock *Block, expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch leftExpr := binaryExpr.left.(type) {

	case *IntExpression:
		switch rightExpr := binaryExpr.right.(type) {

		case *IntExpression:
			newExpr := evalMathExpressionInt(binaryExpr, left.intValue, right.intValue)
			return newExpr

		case *DoubleExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, float64(left.intValue), right.doubleValue)
			return newExpr
		}

	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {

		case *IntExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, left.doubleValue, float64(right.intValue))
			return newExpr

		case *DoubleExpression:
			newExpr := evalMathExpressionDouble(binaryExpr, left.doubleValue, right.doubleValue)
			return newExpr
		}

	case *StringExpression:
		if binaryExpr.operator == AddOperator {
			newExpr := chainString(expr)
			return newExpr
		}
	}

	return expr
}

func evalMathExpressionInt(expr Expression, left, right int) Expression {
	var value int

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case AddOperator:
		value = left + right
	case SubOperator:
		value = left - right
	case MulOperator:
		value = left * right
	case DivOperator:
		value = left / right
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &IntExpression{intValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.IntType}
	return newExpr
}
func evalMathExpressionDouble(expr Expression, left, right float64) Expression {
	var value float64

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case AddOperator:
		value = left + right
	case SubOperator:
		value = left - right
	case MulOperator:
		value = left * right
	case DivOperator:
		value = left / right
	default:
		compileError(binaryExpr.Position(), 0, "")
	}
	newExpr := &DoubleExpression{doubleValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.DoubleType}
	return newExpr
}

func chainString(expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	rightStr := expressionToString(binaryExpr.right)
	if rightStr != "" {
		return expr
	}

	leftStringExpr, ok := binaryExpr.left.(*StringExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	newStr := leftStringExpr.stringValue + rightStr

	newExpr := &StringExpression{stringValue: newStr}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.StringType}

	return newExpr
}

func expressionToString(expr Expression) string {
	var newStr string

	switch e := expr.(type) {
	case *BooleanExpression:
		if e.booleanValue == true {
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

func evalCompareExpression(expr Expression) Expression {
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

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
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.intValue, rightExpr.doubleValue)
			return newExpr
		}
	case *DoubleExpression:
		switch rightExpr := binaryExpr.right.(type) {
		case *IntExpression:
			newExpr := evalCompareExpressionDouble(binaryExpr, leftExpr.doubleValue, rightExpr.intValue)
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
	}
	return expr
}

func evalCompareExpressionBoolean(expr Expression, left, right bool) Expression {
	var value bool

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

	switch binaryExpr.operator {
	case EqOperator:
		value = (left == right)
	case NeOperator:
		value = (left != right)
	default:
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
	return newExpr
}

func evalCompareExpressionInt(expr Expression, left, right int) Expression {
	var value bool

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

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
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
	return newExpr
}

func evalCompareExpressionDouble(expr Expression, left, right float64) Expression {
	var value bool

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

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
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}
	return newExpr
}

func evalCompareExpressionString(expr Expression, left, right string) Expression {
	var value bool
	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(binaryExpr.Position(), 0, "")
	}

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
		compileError(binaryExpr.Position(), 0, "")
	}

	newExpr := &BooleanExpression{booleanValue: value}
	newExpr.typeSpecifier = &TypeSpecifier{basicType: vm.BooleanType}

	return newExpr
}
