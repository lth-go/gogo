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

func getOpcodeTypeOffset(basicType vm.BasicType) byte {
	switch basicType {
	case vm.BooleanType:
		return byte(0)
	case vm.IntType:
		return byte(0)
	case vm.DoubleType:
		return byte(1)
	case vm.StringType:
		return byte(2)
	default:
		panic("basic type")
	}
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}
func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}
