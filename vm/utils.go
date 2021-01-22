package vm

import (
	"encoding/binary"
	"fmt"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}

func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

func createMethodFunctionName(className, methodName string) string {
	ret := fmt.Sprintf("%s#%s", className, methodName)
	return ret
}

func checkNullPointer(obj Object) {
	_, ok := obj.(*ObjectNil)
	if ok {
		vmError(NULL_POINTER_ERR)
	}
}

func GetObjectByType(typ *Type) Object {
	var value Object

	if typ.IsSliceType() {
		value = NilObject
		return value
	}

	switch typ.BasicType {
	case BasicTypeVoid, BasicTypeBool, BasicTypeInt:
		value = NewObjectInt(0)
	case BasicTypeFloat:
		value = NewObjectFloat(0.0)
	case BasicTypeString:
		value = NewObjectString("")
	case BasicTypeNil:
		fallthrough
	default:
		panic("TODO")
	}

	return value
}
