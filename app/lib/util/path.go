package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StandardPath(rpath string) string {
	return path.Join("/", path.Clean(filepath.ToSlash(rpath)))
}

func SimplePath(rpath string) string {
	return strings.TrimPrefix(rpath, "/")
}

func IsDirExist(path string) (bool, error) {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return dirInfo.Mode().IsDir(), nil
}

func SplitPath(rpath string) (string, string) {
	dirPath := filepath.Dir(rpath)
	fileName := filepath.Base(rpath)
	dirPath = filepath.ToSlash(dirPath)
	return dirPath, fileName
}

func GetParentDir(rpath string) string {
	return filepath.ToSlash(filepath.Dir(rpath))
}
