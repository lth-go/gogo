package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func Hash(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IntToBool(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

// Get2ByteInt
// TODO: max is 1 << 15
func Get2ByteInt(b []byte) int {
	return int(int16(binary.BigEndian.Uint16(b)))
}

func Set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

func JsonBytes(v interface{}) []byte {
	buf, _ := json.Marshal(v)

	return buf
}
