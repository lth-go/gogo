package compiler

import (
	"fmt"
	"strings"
)

func printWithIndent(a string, indent int) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(a)
}

func isNull(expr Expression) bool {
	_, ok := expr.(*NullExpression)
	return ok
}

//
// compare
//
func compareType(typ1 *TypeSpecifier, typ2 *TypeSpecifier) bool {
	if typ1.GetBasicType() != typ2.GetBasicType() {
		return false
	}

	t1 := typ1.sliceType
	t2 := typ2.sliceType

	if t1 == nil && t2 == nil {
		return true
	}

	if t1 == nil && t2 != nil {
		return false
	}

	if t1 != nil && t2 == nil {
		return false
	}

	return compareType(t1.ElementType, t2.ElementType)
}

func compareParameter(paramList1, paramList2 []*Parameter) bool {
	length1 := len(paramList1)
	length2 := len(paramList2)
	if length1 != length2 {
		return false
	}

	for i := length1; i < length1; i++ {
		param1 := paramList1[i]
		param2 := paramList2[i]
		if param1.name != param2.name {
			return false
		}
		if !compareType(param1.typeSpecifier, param2.typeSpecifier) {
			return false
		}
	}
	return true
}
