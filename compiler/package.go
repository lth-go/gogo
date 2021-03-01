package compiler

import (
	"log"
)

type Package struct {
	lexer           *Lexer                // 词法解析器
	path            string                // 源文件路径
	packageName     string                // 包名
	importDeclList  []*ImportDecl         // 依赖的包
	funcList        []*FunctionDefinition // 函数列表
	declarationList []*Declaration        // 声明列表
	currentBlock    *Block                // 当前块
}

func (c *Package) GetPackageName() string {
	return c.packageName
}

func (c *Package) SetPackageName(packageName string) {
	c.packageName = packageName
}
func (c *Package) Parse() {
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}
}

func NewPackage(path string) *Package {
	c := &Package{
		lexer:           NewLexer(path),
		path:            path,
		importDeclList:  []*ImportDecl{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
	}

	return c
}

//
//
//

//
// 函数定义
//
func CreateFunctionDefine(pos Position, receiver *Parameter, identifier string, typ *Type, block *Block) {
	c := GetCurrentPackage()

	fd := &FunctionDefinition{
		Type:            typ,
		Name:            identifier,
		PackageName:     c.GetPackageName(),
		Block:           block,
		DeclarationList: nil,
	}

	if block != nil {
		block.parent = &FunctionBlockInfo{Function: fd}
	}

	c.funcList = append(c.funcList, fd)
}

func CreateDeclaration(pos Position, typ *Type, name string, value Expression) *Declaration {
	decl := NewDeclaration(pos, typ, name, value)

	decl.Block = GetCurrentPackage().currentBlock
	if decl.Block != nil {
		decl.IsLocal = true
	}

	return decl
}

func AddDeclList(stmt Statement) {
	decl := stmt.(*Declaration)

	c := GetCurrentPackage()
	decl.PackageName = c.GetPackageName()

	c.declarationList = append(c.declarationList, decl)
}

func SetPackageName(packageName string) {
	c := GetCurrentPackage()
	c.SetPackageName(packageName)
}

func SetImportDeclList(importList []*ImportDecl) {
	c := GetCurrentPackage()
	c.importDeclList = importList
}

func PushCurrentBlock() *Block {
	c := GetCurrentPackage()
	c.currentBlock = &Block{outerBlock: c.currentBlock}

	return c.currentBlock
}

func PopCurrentBlock() *Block {
	c := GetCurrentPackage()
	b := c.currentBlock
	c.currentBlock = c.currentBlock.outerBlock

	return b
}
