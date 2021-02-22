package compiler

func createNativeFuncParamTypeList(typeList []BasicType) []*Parameter {
	if len(typeList) == 0 {
		return nil
	}

	list := make([]*Parameter, 0)
	for _, basicType := range typeList {
		p := &Parameter{
			Type: NewType(basicType),
		}

		if p.Type.IsArray() {
			p.Type.arrayType = NewArrayType(NewType(BasicTypeInterface))
		}

		list = append(list, p)
	}

	return list
}

func (c *Compiler) AddNativeFunctionList() {
	c.AddNativeFunctionPrintf()
	c.AddNativeFunctionLen()
	c.AddNativeFunctionAppend()
	c.AddNativeFunctionDelete()
}

func (c *Compiler) AddNativeFunc(name string, pType, rType []BasicType, ellipsis bool) {
	paramsType := createNativeFuncParamTypeList(pType)
	resultsType := createNativeFuncParamTypeList(rType)

	if ellipsis {
		paramsType[len(paramsType)-1].Ellipsis = true
	}

	fd := &FunctionDefinition{
		Type:            CreateFuncType(paramsType, resultsType),
		Name:            name,
		PackageName:     "_sys",
		Block:           nil,
		DeclarationList: nil,
	}

	c.FuncList = append(c.FuncList, fd)
}

func (c *Compiler) AddNativeFunctionPrintf() {
	c.AddNativeFunc(
		"printf",
		[]BasicType{BasicTypeString, BasicTypeArray},
		nil,
		true,
	)
}

func (c *Compiler) AddNativeFunctionItoa() {
	c.AddNativeFunc(
		"itoa",
		[]BasicType{BasicTypeInt},
		[]BasicType{BasicTypeString},
		false,
	)
}

func (c *Compiler) AddNativeFunctionLen() {
	c.AddNativeFunc(
		"len",
		[]BasicType{BasicTypeInterface},
		[]BasicType{BasicTypeInt},
		false,
	)
}

func (c *Compiler) AddNativeFunctionAppend() {
	c.AddNativeFunc(
		"append",
		[]BasicType{BasicTypeArray, BasicTypeArray},
		[]BasicType{BasicTypeArray},
		true,
	)
}

func (c *Compiler) AddNativeFunctionDelete() {
	c.AddNativeFunc(
		"delete",
		[]BasicType{BasicTypeMap, BasicTypeInterface},
		nil,
		false,
	)
}
