package compiler

import (
	"encoding/binary"
	"fmt"
	"strings"

	"../vm"
)

func printWithIndent(a string, indent int) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(a)
}

func isNull(expr Expression) bool {
	_, ok := expr.(*NullExpression)
	return ok
}

// TODO 作为TypeSpecifier的方法
func isVoid(t *TypeSpecifier) bool    { return t.basicType == vm.VoidType }
func isBoolean(t *TypeSpecifier) bool { return t.basicType == vm.BooleanType }
func isInt(t *TypeSpecifier) bool     { return t.basicType == vm.IntType }
func isDouble(t *TypeSpecifier) bool  { return t.basicType == vm.DoubleType }
func isString(t *TypeSpecifier) bool  { return t.basicType == vm.StringType }
func isClass(t *TypeSpecifier) bool   { return t.basicType == vm.ClassType }
func isModule(t *TypeSpecifier) bool  { return t.basicType == vm.ModuleType }
func isObject(t *TypeSpecifier) bool  { return isString(t) || isArray(t) }
func isArray(t *TypeSpecifier) bool {
	if t.deriveList == nil || len(t.deriveList) == 0 {
		return false
	}
	firstElem := t.deriveList[0]
	_, ok := firstElem.(*ArrayDerive)
	return ok
}

func getTypeName(typ *TypeSpecifier) string {
	typeName := getBasicTypeName(typ.basicType)

	for _, derive := range typ.deriveList {
		switch derive.(type) {
		case *FunctionDerive:
			panic("TODO:derive_tag, func")
		case *ArrayDerive:
			typeName = typeName + "[]"
		default:
			print("=====\n", typ.Position().Line)
			panic("TODO:derive_tag")
		}
	}

	return typeName
}

func getBasicTypeName(typ vm.BasicType) string {
	switch typ {
	case vm.BooleanType:
		return "boolean"
	case vm.IntType:
		return "int"
	case vm.DoubleType:
		return "double"
	case vm.StringType:
		return "string"
	case vm.NullType:
		return "null"
	default:
		panic(fmt.Sprintf("bad case. type..%d\n", typ))
	}
}

func getOpcodeTypeOffset(typ *TypeSpecifier) byte {

	if typ.deriveList != nil && len(typ.deriveList) != 0 {
		if !typ.isArrayDerive() {
			panic("TODO")
		}
		return 2
	}
	switch typ.basicType {
	case vm.VoidType:
		panic("basic type is void")
	case vm.BooleanType, vm.IntType:
		return byte(0)
	case vm.DoubleType:
		return byte(1)
	case vm.StringType, vm.ClassType:
		return byte(2)
	case vm.NullType, vm.BaseType:
		fallthrough
	default:
		panic("basic type")
	}
	return byte(0)
}

func get2ByteInt(b []byte) int {
	return int(binary.BigEndian.Uint16(b))
}
func set2ByteInt(b []byte, value int) {
	binary.BigEndian.PutUint16(b, uint16(value))
}

//
// compare
//
func compareType(typ1 *TypeSpecifier, typ2 *TypeSpecifier) bool {
	if typ1.basicType != typ2.basicType {
		return false
	}

	typ1Len := len(typ1.deriveList)
	typ2Len := len(typ2.deriveList)
	if typ1Len != typ2Len {
		return false
	}

	for i := 0; i < typ1Len; i++ {
		derive1 := typ1.deriveList[i]
		derive2 := typ2.deriveList[i]
		switch d1 := derive1.(type) {
		case *ArrayDerive:
			switch derive2.(type) {
			case *ArrayDerive:
				// pass
			default:
				return false
			}
		case *FunctionDerive:
			switch d2 := derive2.(type) {
			case *FunctionDerive:
				if !compareParameter(d1.parameterList, d2.parameterList) {
					return false
				}
			default:
				return false
			}
		default:
			panic("TODO")
		}
	}
	return true
}

func compareParameter(paramList1, paramList2 []*Parameter) bool {
	length1 := len(paramList1)
	length2 := len(paramList2)
	if length1 != length2 {
		return false
	}

	for i := length1; i < length1; i++ {
		param1 := paramList1[i]
		param2 := paramList2[i]
		if param1.name != param2.name {
			return false
		}
		if !compareType(param1.typeSpecifier, param2.typeSpecifier) {
			return false
		}
	}
	return true
}

//
// search
//
func searchDeclaration(name string, currentBlock *Block) *Declaration {

	// 从局部作用域查找
	for block := currentBlock; block != nil; block = block.outerBlock {
		for _, declaration := range block.declarationList {
			if declaration.name == name {
				return declaration
			}
		}
	}

	// 从全局作用域查找
	compiler := getCurrentCompiler()
	for _, declaration := range compiler.declarationList {
		if declaration.name == name {
			return declaration
		}
	}

	return nil
}

func searchFunction(name string) *FunctionDefinition {
	compiler := getCurrentCompiler()

	// 当前compiler查找
	for _, pos := range compiler.funcList {
		if pos.name == name {
			return pos
		}
	}

	// 导入的compiler查找
	for _, required := range compiler.requiredList {
		for _, fd := range required.funcList {
			if fd.name == name && fd.classDefinition == nil {
				return fd
			}
		}
	}

	return nil
}

func searchModule(name string) *Module {
	compiler := getCurrentCompiler()

	for _, requiredCompiler := range compiler.requiredList {
		// 暂无处理重名
		lastName := requiredCompiler.packageNameList[len(requiredCompiler.packageNameList)-1]
		if name == lastName {
			return &Module{
				compiler: requiredCompiler,
				typ: &TypeSpecifier{basicType: vm.ModuleType},
			}
		}

	}
	return nil
}

// 根据名字在当前compiler, 及required里搜索类定义
func searchClass(identifier string) *ClassDefinition {

	compiler := getCurrentCompiler()

	for _, cd := range compiler.classDefinitionList {
		if cd.name == identifier {
			return cd
		}
	}

	for _, requiredCompiler := range compiler.requiredList {
		for _, cd := range requiredCompiler.classDefinitionList {
			if cd.name == identifier {
				return cd
			}
		}
	}

	return nil
}

func searchClassAndAdd(pos Position, name string, classIndexP *int) *ClassDefinition {

	cd := searchClass(name)

	if cd == nil {
		compileError(pos, CLASS_NOT_FOUND_ERR, name)
	}

	*classIndexP = cd.addToCurrentCompiler()

	return cd
}

func searchCompiler(list []*Compiler, packageName []string) *Compiler {
	for _, c := range list {
		if comparePackageName(c.packageNameList, packageName) {
			return c
		}
	}
	return nil
}

func comparePackageName(packageNameList1, packageNameList2 []string) bool {
	if packageNameList1 == nil {
		if packageNameList2 == nil {
			return true
		}
		return false
	}

	length1 := len(packageNameList1)
	length2 := len(packageNameList2)

	if length1 != length2 {
		return false
	}

	for i := 0; i < length1; i++ {
		if packageNameList1[i] != packageNameList2[i] {
			return false
		}
	}

	return true
}

// TODO
func createMethodFunctionName(className, methodName string) string {
	return fmt.Sprintf("%s#%s", className, methodName)
}
