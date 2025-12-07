// Wraps some essential OS/IO path functionalities to make them more digestible and clarify their API.
//
// None of those functionalities rely on existing files or folders, they simply operate on a string.
package path

import (
	"path/filepath"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
)

func New(elements ...string) string {
	return internal.Path(filepath.Join(elements...))
}

func IsDirectory(path string) bool {
	path = internal.Path(path)
	if text.EndsWith(path, "/") {
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
	return internal.Path(filepath.Dir(path))
}
func Extension(path string) string {
	return internal.Path(filepath.Ext(path))
}
func RemoveExtension(path string) string {
	var ext = Extension(path)
	if ext == "" {
		return path
	}
	return internal.Path(path[:len(path)-len(ext)])
}
