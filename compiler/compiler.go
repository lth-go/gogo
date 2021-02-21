package compiler

import (
	"log"
)

type Compiler struct {
	lexer           *Lexer                // 词法解析器
	path            string                // 源文件路径
	packageName     string                // 包名
	importList      []*Import             // 依赖的包
	funcList        []*FunctionDefinition // 函数列表
	declarationList []*Declaration        // 声明列表
	currentBlock    *Block                // 当前块
}

func (c *Compiler) GetPackageName() string {
	return c.packageName
}

func (c *Compiler) SetPackageName(packageName string) {
	c.packageName = packageName
}

func (c *Compiler) Parse() {
	if yyParse(c.lexer) != 0 {
		log.Fatalf("\nFileName: %s%s", c.path, c.lexer.e)
	}
}

func NewCompiler(path string) *Compiler {
	c := &Compiler{
		lexer:           NewLexer(path),
		path:            path,
		importList:      []*Import{},
		funcList:        []*FunctionDefinition{},
		declarationList: []*Declaration{},
	}

	return c
}

// 添加声明
func (c *Compiler) AddDeclarationList(decl *Declaration) {
	c.declarationList = append(c.declarationList, decl)
}

func (c *Compiler) SearchDeclaration(name string) *Declaration {
	for _, decl := range c.declarationList {
		if decl.Name == name {
			return decl
		}
	}
	return nil
}
