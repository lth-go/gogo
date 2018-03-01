package compiler

import (
	"os"
	"path/filepath"
)

const (
	HEAD_SOURCE int = iota
	SOURCE
)

const (
	REQUIRE_SUFFIX string = ".gogogo"
)

type Require struct {
	PosImpl

	packageNameList []string
}

func createRequire(packageNameList []string) *Require {
	return &Require{
		packageNameList: packageNameList,
	}
}

func createRequireList(packageNameList []string) []*Require {
	req := createRequire(packageNameList)

	return []*Require{req}
}

func chainRequireList(requireList1, requireList2 []*Require) []*Require {
	return append(requireList1, requireList2...)
}

// 获取导入文件的相对路径
func (r *Require) getRelativePath() string {
	path := filepath.Join(r.packageNameList...)
	path = path + REQUIRE_SUFFIX
	return path
}

func (r *Require) getFullPath() string {
	searchBasePath := os.Getenv("REQUIRE_SEARCH_PATH")
	if searchBasePath == "" {
		searchBasePath = "."
	}

	relativePath := r.getRelativePath()

	fullPath := filepath.Join(searchBasePath, relativePath)
	_, err := os.Stat(fullPath)
	if err != nil {
		compileError(r.Position(), REQUIRE_FILE_NOT_FOUND_ERR, fullPath)
	}

	return fullPath
}

func createPackageName(lit string) []string {
	return []string{lit}
}

func chainPackageName(list []string, lit string) []string {
	return append(list, lit)
}
