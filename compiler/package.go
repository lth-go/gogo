package compiler

import (
	"log"
)

type Package struct {
	lexer           *Lexer                // 词法解析器
	path            string                // 源文件路径
	packageName     string                // 包名
	importList      []*Import             // 依赖的包
	funcList        []*FunctionDefinition // 函数列表
	declarationList []*Declaration        // 声明列表
	typeDefList     []*TypeDefDecl        // 类型声明列表
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
		importList:      []*Import{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
		typeDefList:     []*TypeDefDecl{},
	}

	return c
}
