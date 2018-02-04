package compiler

import (
	"fmt"
)

func debug(format string, a ...interface{}) {
	fmt.Println("=========")
	fmt.Printf(format, a)
	fmt.Println("\n=========")
}
