package compiler

import (
	"os"
	"path/filepath"
)

const (
	requireSuffix = ".4g"
)

type Require struct {
	PosImpl

	packageNameList []string // remove
	packageName     string
}

// 获取导入文件的相对路径
func (r *Require) getRelativePath() string {
	path := filepath.Join(r.packageNameList...)
	path = path + requireSuffix
	return path
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

func createImport(packageName string) *Require {
	return &Require{
		packageName:     packageName,
		packageNameList: []string{packageName},
	}
}
