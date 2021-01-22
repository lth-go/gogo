package compiler

import (
	"os"
	"path/filepath"
)

const (
	importSuffix = ".gogo"
)

type Import struct {
	PosBase
	packageName string
}

// 获取导入文件的相对路径
func (i *Import) getRelativePath() string {
	return i.packageName + importSuffix
}

func (i *Import) GetPath() string {
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

func CreateImportList(importSpec *Import) []*Import {
	return []*Import{importSpec}
}

func CreateImport(packageName string) *Import {
	return &Import{
		packageName: packageName,
	}
}
