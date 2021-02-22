package compiler

import (
	"github.com/lth-go/gogo/vm"
)

func (cm *Compiler) GetVmVariableList() []vm.Object {
	variableList := make([]vm.Object, 0)

	for _, decl := range cm.DeclarationList {
		variableList = append(variableList, GetVmVariable(decl.Value))
	}

	return variableList
}

func (cm *Compiler) GetVmFunctionList() []*vm.GoGoFunction {
	vmFuncList := make([]*vm.GoGoFunction, 0)

	for _, fd := range cm.FuncList {
		// TODO: 过滤掉_sys, 由虚拟机自己添加
		if fd.PackageName == "_sys" {
			continue
		}

		variableList := make([]vm.Object, 0)

		for _, variable := range fd.DeclarationList {
			variableList = append(variableList, GetVmVariable(variable.Value))
		}

		vmFuncList = append(vmFuncList, &vm.GoGoFunction{
			ParamCount:   len(fd.GetType().funcType.Params),
			ResultCount:  len(fd.GetType().funcType.Results),
			VariableList: variableList,
			CodeList:     fd.CodeList,
		})
	}

	return vmFuncList
}

func GetVmVariable(valueIFS Expression) vm.Object {
	if valueIFS == nil {
		return nil
	}

	switch value := valueIFS.(type) {
	case *BoolExpression:
		v := 0
		if value.Value {
			v = 1
		}
		return vm.NewObjectInt(v)
	case *IntExpression:
		return vm.NewObjectInt(value.Value)
	case *FloatExpression:
		return vm.NewObjectFloat(value.Value)
	case *StringExpression:
		return vm.NewObjectString(value.Value)
	case *InterfaceExpression:
		return vm.NewObjectInterface(GetVmVariable(value.Data))
	case *NilExpression:
		return vm.NilObject
	case *ArrayExpression:
		arrayValue := vm.NewObjectArray(len(value.List))
		for i, subValue := range value.List {
			arrayValue.List[i] = GetVmVariable(subValue)
		}
		return arrayValue
	case *MapExpression:
		mapValue := vm.NewObjectMap()
		length := len(value.KeyList)
		for i := 0; i < length; i++ {
			mapValue.Set(GetVmVariable(value.KeyList[i]), GetVmVariable(value.ValueList[i]))
		}
		return mapValue
	case *StructExpression:
		structValue := vm.NewObjectStruct(len(value.FieldList))
		for i, subValue := range value.FieldList {
			structValue.FieldList[i] = GetVmVariable(subValue)
		}
		return structValue
	}

	return nil
}
