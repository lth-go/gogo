package vm

import (
	"fmt"
)

import (
	"encoding/binary"
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

func PrintCode(codeList []byte) {
	print("=====\n", "code list start\n=====\n")
	for i := 0; i < len(codeList); {
		code := codeList[i]
		info := OpcodeInfo[int(code)]
		paramList := []byte(info.Parameter)

		fmt.Println(info.Mnemonic)
		for _, param := range paramList {
			switch param {
			case 'b':
				i += 1
			case 's', 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
		i += 1
	}
	print("=====\n", "code list end\n=====\n")
}

func initializeValue(typ *VmTypeSpecifier) VmValue {
	var value VmValue

	if typ.isArrayDerive() {
		value = vmNullObjectRef
		return value
	}

	switch typ.BasicType {
	case VoidType, BooleanType, IntType:
		value = &VmIntValue{intValue: 0}

	case DoubleType:
		value = &VmDoubleValue{doubleValue: 0.0}

	case StringType, ClassType:
		value = vmNullObjectRef

	case NullType, BaseType:
		fallthrough
	default:
		panic("TODO")
	}

	return value
}

func createMethodFunctionName(class_name, method_name string) string {

	ret := fmt.Sprintf("%s#%s", class_name, method_name)

	return ret
}

func check_null_pointer(obj *VmObjectRef) {
	if obj.data == nil {
		vmError(NULL_POINTER_ERR)
	}
}
