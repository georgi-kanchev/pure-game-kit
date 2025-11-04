package path

import (
	"path/filepath"
	"pure-game-kit/utility/text"
)

func New(elements ...string) string {
	return filepath.Join(elements...)
}

func Divider() string {
	return string(filepath.Separator)
}

func IsDirectory(path string) bool {
	if text.EndsWith(path, string(filepath.Separator)) {
		return true
	}
	return Extension(path) == ""
}
func IsFile(path string) bool {
	return Extension(path) != ""
}

func LastPart(path string) string {
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
