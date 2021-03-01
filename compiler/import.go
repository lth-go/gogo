package compiler

import (
	"os"
	"path/filepath"
)

const (
	importSuffix = ".gogo"
)

type ImportDecl struct {
	PosBase
	packageName string
}

// 获取导入文件的相对路径
func (i *ImportDecl) getRelativePath() string {
	return i.packageName + importSuffix
}

func (i *ImportDecl) GetPath() string {
	searchBasePath := os.Getenv("IMPORT_SEARCH_PATH")
	if searchBasePath == "" {
		searchBasePath = "."
	}

	relativePath := i.getRelativePath()

	fullPath := filepath.Join(searchBasePath, relativePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		compileError(i.Position(), REQUIRE_FILE_NOT_FOUND_ERR, fullPath)
	}

	return fullPath
}

func CreateImportDeclList(importSpec *ImportDecl) []*ImportDecl {
	return []*ImportDecl{importSpec}
}

func CreateImportDecl(packageName string) *ImportDecl {
	return &ImportDecl{
		packageName: packageName,
	}
}
