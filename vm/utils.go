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
			case 's':
				fallthrough
			case 'p':
				i += 2
			default:
				panic("TODO")
			}
		}
		i += 1
	}
	print("=====\n", "code list end\n=====\n")
}
