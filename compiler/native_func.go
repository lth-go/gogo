package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func (c *Compiler) AddNativeFunctions() {
	c.AddNativeFunctionPrint()
	c.AddNativeFunctionItoa()
}

func (c *Compiler) AddNativeFunctionPrint() {
	paramsType := []*Parameter{{Type: NewType(vm.BasicTypeString), Name: "str"}}
	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, nil),
		Name:            "print",
		PackageName:     "_sys",
		ParameterList:   paramsType,
		Block:           nil,
		DeclarationList: nil,
	}

	c.funcList = append(c.funcList, fd)
}

func (c *Compiler) AddNativeFunctionItoa() {
	paramsType := []*Parameter{{Type: NewType(vm.BasicTypeInt), Name: "int"}}
	resultsType := []*Parameter{{Type: NewType(vm.BasicTypeString)}}
	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, resultsType),
		Name:            "itoa",
		PackageName:     "_sys",
		ParameterList:   paramsType,
		Block:           nil,
		DeclarationList: nil,
	}

	c.funcList = append(c.funcList, fd)
}
