package compiler

// ==============================
// MemberDeclaration
// ==============================
type MemberDeclaration interface{}

//
// MethodMember
//
type MethodMember struct {
	PosImpl
	functionDefinition *FunctionDefinition
	methodIndex        int
}

//
// FieldMember
//
type FieldMember struct {
	PosImpl
	name          string
	typeSpecifier *TypeSpecifier
	fieldIndex    int
}
