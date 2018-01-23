package compiler

//
// cast func
//

// 声明类型转换
func createAssignCast(src Expression, dest *TypeSpecifier) Expression {
	var castExpr Expression
	if src.typeS().deriveList != nil || dest.deriveList != nil {
		compileError(src.Position(), 0, "")
	}

	if src.typeS().basicType == dest.basicType {
		return src
	}

	if src.typeS().basicType == vm.IntType && dest.basicType == DoubleType {
		castExpr = allocCastExpression(IntToDoubleCast, src)
		return castExpr

	} else if src.typeS().basicType == DoubleType && dest.basicType == vm.IntType {
		castExpr = alloc_cast_expression(DoubleToIntCast, src)
		return castExpr
	}

	compileError(src.Position(), 0, "")
	return nil
}

func castBinaryExpression(expr Expression) Expression {

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		compileError(expr.Position(), 0, "")
	}

	if isString(binaryExpr.left.typeS()) && isBoolean(binaryExpr.right.typeS()) {
		newExpr := allocCastExpression(BooleanToStringCast, binaryExpr.right)
		return newExpr

	} else if isString(binaryExpr.left.typeS()) && isDouble(binaryExpr.right.typeS()) {
		newExpr := allocCastExpression(NumberToStringCast, binaryExpr.right)
		return newExpr

	}
	return expr
}
