package compiler

import (
	"fmt"
	"strings"
)

func printWithIndent(a string, indent int) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(a)
}

//
// compare
//
func compareType(typ1 *Type, typ2 *Type) bool {
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
		if param1.Name != param2.Name {
			return false
		}
		if !compareType(param1.Type, param2.Type) {
			return false
		}
	}
	return true
}
