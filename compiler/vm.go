package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func (c *Compiler) GetVmVariableList() []*vm.Variable {
	variableList := make([]*vm.Variable, 0)

	for _, decl := range c.declarationList {
		newValue := vm.NewVmVariable(decl.PackageName, decl.Name, CopyToVmType(decl.Type))
		newValue.Value = GetVmVariable(decl.Value)
		variableList = append(variableList, newValue)
	}

	return variableList
}

func GetVmVariable(valueIFS Expression) vm.Object {
	if valueIFS == nil {
		return nil
	}

	switch value := valueIFS.(type) {
	// case *BoolExpression:
	//     return value.Value
	case *IntExpression:
		return vm.NewObjectInt(value.Value)
	case *FloatExpression:
		return vm.NewObjectFloat(value.Value)
	case *StringExpression:
		return vm.NewObjectString(value.Value)
	case *ArrayExpression:
		arrayValue := vm.NewObjectArray(len(value.List))
		for i, subValue := range value.List {
			arrayValue.List[i] = GetVmVariable(subValue)
		}
		return arrayValue
	}

	return nil
}

func (c *Compiler) GetVmFunctionList() []*vm.Function {
	vmFuncList := make([]*vm.Function, 0)

	for _, fd := range c.funcList {
		vmFunc := c.GetVmFunction(fd, fd.GetPackageName() == c.GetPackageName())
		vmFuncList = append(vmFuncList, vmFunc)
	}

	return vmFuncList
}

func (c *Compiler) GetVmFunction(src *FunctionDefinition, inThisExe bool) *vm.Function {
	ob := NewOpCodeBuf()

	dest := &vm.Function{
		PackageName: src.GetPackageName(),
		Name:        src.Name,
		Type:        CopyToVmType(src.GetType().Copy()),
		IsMethod:    false,
	}

	if src.Block != nil && inThisExe {
		generateStatementList(src.Block.statementList, ob)

		dest.IsImplemented = true
		dest.CodeList = ob.fixOpcodeBuf()
		dest.LineNumberList = ob.lineNumberList
		dest.LocalVariableList = copyVmVariableList(src)
	} else {
		dest.IsImplemented = false
		dest.LocalVariableList = nil
	}

	return dest
}
