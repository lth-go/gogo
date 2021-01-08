package compiler

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	importSuffix = ".gogo"
)

type ImportSpec struct {
	PosImpl
	packageName string
}

// 获取导入文件的相对路径
func (i *ImportSpec) getRelativePath() string {
	return i.packageName + importSuffix
}

func (i *ImportSpec) getFullPath() string {
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

func (i *ImportSpec) getPackageNameList() []string {
	return strings.Split(i.packageName, "/")
}

func createImportSpecList(importSpec *ImportSpec) []*ImportSpec {
	return []*ImportSpec{importSpec}
}

func createImportSpec(packageName string) *ImportSpec {
	return &ImportSpec{
		packageName: packageName,
	}
}
