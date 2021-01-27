package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func (c *Compiler) AddNativeFunctionList() {
	c.AddNativeFunctionPrintf()
	c.AddNativeFunctionItoa()
	c.AddNativeFunctionLen()
	c.AddNativeFunctionAppend()
}

func (c *Compiler) AddNativeFunc(name string, pType, rType []vm.BasicType, ellipsis bool) {
	paramsType := createNativeFuncParamTypeList(pType)
	resultsType := createNativeFuncParamTypeList(rType)

	if ellipsis {
		paramsType[len(paramsType)-1].Ellipsis = true
	}

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

func (c *Compiler) AddNativeFunctionPrintf() {
	c.AddNativeFunc(
		"printf",
		[]vm.BasicType{vm.BasicTypeString, vm.BasicTypeArray},
		nil,
		true,
	)
}

func (c *Compiler) AddNativeFunctionItoa() {
	c.AddNativeFunc(
		"itoa",
		[]vm.BasicType{vm.BasicTypeInt},
		[]vm.BasicType{vm.BasicTypeString},
		false,
	)
}

func (c *Compiler) AddNativeFunctionLen() {
	c.AddNativeFunc(
		"len",
		[]vm.BasicType{vm.BasicTypeInterface},
		[]vm.BasicType{vm.BasicTypeInt},
		false,
	)
}

func (c *Compiler) AddNativeFunctionAppend() {
	c.AddNativeFunc(
		"append",
		[]vm.BasicType{vm.BasicTypeArray, vm.BasicTypeArray},
		[]vm.BasicType{vm.BasicTypeArray},
		true,
	)
}

func createNativeFuncParamTypeList(typeList []vm.BasicType) []*Parameter {
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
