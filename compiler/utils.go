package compiler

import (
	"fmt"
	"strings"
)

func printWithIdent(a string, ident int) {
	fmt.Print(strings.Repeat(" ", ident))
	fmt.Println(a)
}
