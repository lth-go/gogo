package compiler

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	src, err := ioutil.ReadFile("../test.4g")
	if err != nil {
		panic(err)
	}
	compiler, err := ParseSrc(string(src))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, stmt := range compiler.statementList {
		stmt.show(0)
	}
}
