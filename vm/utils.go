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

func initializeValue(typ *TypeSpecifier) Value {
	var value Value

	if typ.isArrayDerive() {
		value = vmNullObjectRef
		return value
	}

	switch typ.BasicType {
	case VoidType, BooleanType, IntType:
		value = &IntValue{intValue: 0}

	case DoubleType:
		value = &DoubleValue{doubleValue: 0.0}

	case StringType:
		value = vmNullObjectRef

	case NullType, BaseType:
		fallthrough
	default:
		panic("TODO")
	}

	return value
}

func createMethodFunctionName(className, methodName string) string {

	ret := fmt.Sprintf("%s#%s", className, methodName)

	return ret
}

func checkNullPointer(obj *ObjectRef) {
	if obj.data == nil {
		vmError(NULL_POINTER_ERR)
	}
}
