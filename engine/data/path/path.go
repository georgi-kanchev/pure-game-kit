package path

import (
	"os"
	"path/filepath"
	"strings"
)

func New(elements ...string) string {
	return filepath.Join(elements...)
}

func Separator() string {
	return string(filepath.Separator)
}

func Executable() string {
	var execPath, err = os.Executable()
	if err != nil {
		return ""
	}
	return execPath
}

func IsDirectory(path string) bool {
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return true
	}
	return Extension(path) == ""
}
func IsFile(path string) bool {
	return Extension(path) != ""
}

func LastElement(path string) string {
	return filepath.Base(path)
}
func Folder(path string) string {
	return filepath.Dir(path)
}
func Extension(path string) string {
	return filepath.Ext(path)
}
func RemoveExtension(path string) string {
	var ext = Extension(path)
	if ext == "" {
		return path
	}
	return path[:len(path)-len(ext)]
}
