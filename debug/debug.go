package debug

import (
	"fmt"
)

func Printf(format string, a ...interface{}) {
	fmt.Println("=========")
	fmt.Printf(format, a)
	fmt.Println("\n=========")
}
