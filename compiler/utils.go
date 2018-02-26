package compiler

import (
	"encoding/binary"
	"fmt"
	"strings"

	"../vm"
)

func printWithIdent(a string, ident int) {
	fmt.Print(strings.Repeat(" ", ident))
	fmt.Println(a)
}

func isInt(t *TypeSpecifier) bool     { return t.basicType == vm.IntType }
func isDouble(t *TypeSpecifier) bool  { return t.basicType == vm.DoubleType }
func isBoolean(t *TypeSpecifier) bool { return t.basicType == vm.BooleanType }
func isString(t *TypeSpecifier) bool  { return t.basicType == vm.StringType }
func isVoid(t *TypeSpecifier) bool    { return t.basicType == vm.VoidType }

func isNull(expr Expression) bool {
	_, ok := expr.(*NullExpression)
	return ok
}

func isArray(t *TypeSpecifier) bool {
	if t.deriveList == nil || len(t.deriveList) == 0 {
		return false
	}
	firstElem := t.deriveList[0]
	_, ok := firstElem.(*ArrayDerive)
	return ok
}

func isObject(t *TypeSpecifier) bool {
	return isString(t) || isArray(t)
}

func getOpcodeTypeOffset(typ *TypeSpecifier) byte {

	if typ.deriveList != nil && len(typ.deriveList) != 0 {
		if !typ.isArrayDerive() {
			panic("TODO")
		}
		return 2
	}
	switch typ.basicType {
	case vm.BooleanType:
		return byte(0)
	case vm.IntType:
		return byte(0)
	case vm.DoubleType:
		return byte(1)
	case vm.StringType:
		return byte(2)
	case vm.NullType:
		fallthrough
	default:
		panic("basic type")
	}
	return byte(0)
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}
func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

func compareType(typ1 *TypeSpecifier, typ2 *TypeSpecifier) bool {
	if typ1.basicType != typ2.basicType {
		return false
	}

	typ1Len := len(typ1.deriveList)
	typ2Len := len(typ2.deriveList)
	if typ1Len != typ2Len {
		return false
	}

	for i := 0; i < typ1Len; i++ {
		derive1 := typ1.deriveList[i]
		derive2 := typ2.deriveList[i]
		switch d1 := derive1.(type) {
		case *ArrayDerive:
			switch derive2.(type) {
			case *ArrayDerive:
				// pass
			default:
				return false
			}
		case *FunctionDerive:
			switch d2 := derive2.(type) {
			case *FunctionDerive:
				if !compareParameter(d1.parameterList, d2.parameterList) {
					return false
				}
			default:
				return false
			}
		default:
			panic("TODO")
		}
	}
	return true
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
