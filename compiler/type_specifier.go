package compiler

//
// derive
//

type TypeDerive interface{}

type FunctionDerive struct {
	parameterList []*Parameter
}

type ArrayDerive struct{}

//
// TypeSpecifier
//
type classRef struct {
	identifier      string
	classDefinition *ClassDefinition
	classIndex      int
}

// TypeSpecifier 表达式类型, 包括基本类型和派生类型
type TypeSpecifier struct {
	PosImpl

	// 基本类型
	basicType vm.BasicType

	// 类引用
	classRef classRef

	// 派生类型
	deriveList []TypeDerive
}

func (t *TypeSpecifier) fix() {
	compiler := getCurrentCompiler()

	for _, derive_pos := range t.deriveList {
		derive, ok := derive_pos.(*FunctionDerive)
		if ok {
			for _, parameter := range derive.parameterList {
				parameter.typeSpecifier.fix()
			}
		}
	}

	if typ.basicType == vm.ClassType && typ.classRef.classDefinition == nil {

		cd := searchClass(typ.classRef.identifier)
		if cd == nil {
			compileError(t.Position(), TYPE_NAME_NOT_FOUND_ERR, t.classRef.identifier)
			return nil
		}
		if !compare_package_name(cd.package_name, compiler.package_name) {
			compileError(typ.Position(), PACKAGE_CLASS_ACCESS_ERR, cd.name)
		}
		typ.classRef.classDefinition = cd
		typ.classRef.classIndex = cd.add_class()
		return
	}
}

func createTypeSpecifier(basicType vm.BasicType, pos Position) *TypeSpecifier {
	typ := &TypeSpecifier{basicType: basicType}
	typ.SetPosition(pos)
	return typ
}
func create_class_type_specifier(identifier string, pos Position) *TypeSpecifier {

	typ := &TypeSpecifier{
		basicType: vm.ClassType,
		classRef: classRef{
			identifier: identifier,
		},
	}
	typ.SetPosition(pos)

	return typ
}

func create_array_type_specifier(typ *TypeSpecifier) *TypeSpecifier {
	typ.appendDerive(&ArrayDerive{})
	return typ
}

func (t *TypeSpecifier) appendDerive(derive TypeDerive) {
	if t.deriveList == nil {
		t.deriveList = []TypeDerive{}
	}
	t.deriveList = append(t.deriveList, derive)
}

func (t *TypeSpecifier) isArrayDerive() bool {
	return isArray(t)
}
