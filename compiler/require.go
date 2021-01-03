package compiler

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	requireSuffix = ".4g"
)

type Require struct {
	PosImpl
	packageName string
}

// 获取导入文件的相对路径
func (r *Require) getRelativePath() string {
	return r.packageName + requireSuffix
}

func (r *Require) getFullPath() string {
	// TODO 暂时写死, 方便测试
	searchBasePath := os.Getenv("REQUIRE_SEARCH_PATH")
	if searchBasePath == "" {
		searchBasePath = "."
	}
	// searchBasePath := "/home/lth/toy/gogogogo/test"

	relativePath := r.getRelativePath()

	fullPath := filepath.Join(searchBasePath, relativePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		compileError(r.Position(), REQUIRE_FILE_NOT_FOUND_ERR, fullPath)
	}

	return fullPath
}

func (r *Require) getPackageNameList() []string {
	return strings.Split(r.packageName, "/")
}

func createImport(packageName string) *Require {
	return &Require{
		packageName: packageName,
	}
}
