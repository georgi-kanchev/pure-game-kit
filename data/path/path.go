package path

import (
	"path/filepath"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
)

func New(elements ...string) string {
	return filepath.Join(elements...)
}

func Divider() string {
	return string(filepath.Separator)
}

func Executable() string {
	return internal.ExecutablePath()
}

func IsDirectory(path string) bool {
	path = internal.MakeAbsolutePath(path)
	if text.EndsWith(path, string(filepath.Separator)) {
		return true
	}
	return Extension(path) == ""
}
func IsFile(path string) bool {
	path = internal.MakeAbsolutePath(path)
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
