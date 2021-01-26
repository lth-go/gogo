package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func (c *Compiler) AddNativeFunctionList() {
	c.AddNativeFunctionPrint()
	c.AddNativeFunctionPrintf()
	c.AddNativeFunctionItoa()
}
func (c *Compiler) AddNativeFunc(name string, pType, rType []vm.BasicType) {
	paramsType := TODOCreateParam(pType)
	resultsType := TODOCreateParam(rType)

	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, resultsType),
		Name:            name,
		PackageName:     "_sys",
		ParameterList:   paramsType,
		Block:           nil,
		DeclarationList: nil,
	}

	c.funcList = append(c.funcList, fd)
}

func (c *Compiler) AddNativeFunctionPrint() {
	c.AddNativeFunc(
		"print",
		[]vm.BasicType{vm.BasicTypeString},
		nil,
	)
}

func (c *Compiler) AddNativeFunctionPrintf() {
	c.AddNativeFunc(
		"printf",
		[]vm.BasicType{vm.BasicTypeString, vm.BasicTypeArray},
		nil,
	)
}

func (c *Compiler) AddNativeFunctionItoa() {
	c.AddNativeFunc(
		"itoa",
		[]vm.BasicType{vm.BasicTypeInt},
		[]vm.BasicType{vm.BasicTypeString},
	)
}

func TODOCreateParam(typeList []vm.BasicType) []*Parameter {
	if len(typeList) == 0 {
		return nil
	}

	list := make([]*Parameter, 0)
	for _, basicType := range typeList {
		p := &Parameter{
			Type: NewType(basicType),
		}

		if p.Type.IsArray() {
			p.Type.arrayType = NewArrayType(NewType(vm.BasicTypeInterface))
		}

		list = append(list, p)
	}

	return list
}
