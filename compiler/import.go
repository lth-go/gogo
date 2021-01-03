package compiler

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	requireSuffix = ".4g"
)

type ImportSpec struct {
	PosImpl
	packageName string
}

// 获取导入文件的相对路径
func (i *ImportSpec) getRelativePath() string {
	return i.packageName + requireSuffix
}

func (i *ImportSpec) getFullPath() string {
	// TODO 暂时写死, 方便测试
	searchBasePath := os.Getenv("REQUIRE_SEARCH_PATH")
	if searchBasePath == "" {
		searchBasePath = "."
	}
	// searchBasePath := "/home/lth/toy/gogogogo/test"

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
