package folder

import (
	"io"
	"os"
	ph "pure-kit/engine/data/path"
	"pure-kit/engine/internal"
)

func IsExisting(path string) bool {
	path = internal.MakeAbsolutePath(path)
	var info, err = os.Stat(path)
	return err == nil && info.IsDir()
}
func IsEmpty(path string) bool {
	path = internal.MakeAbsolutePath(path)
	return len(Content(path)) == 0
}
func ByteSize(path string) int64 {
	path = internal.MakeAbsolutePath(path)
	var totalSize int64

	var files = Files(path)
	for _, file := range files {
		var info, err = os.Stat(ph.New(path, file))
		if err == nil {
			totalSize += info.Size()
		}
	}

	var folders = Folders(path)
	for _, sub := range folders {
		totalSize += ByteSize(ph.New(path, sub))
	}

	return totalSize
}
func TimeOfLastEdit(path string) (year, month, day, minute int) {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return 0, 0, 0, 0
	}

	var info, err = os.Stat(path)
	if err != nil {
		return 0, 0, 0, 0
	}

	var t = info.ModTime()
	year = t.Year()
	month = int(t.Month()) // time.Month is 1-based already
	day = t.Day()          // day of the month
	minute = t.Hour()*60 + t.Minute()
	return
}

func Create(path string) bool {
	path = internal.MakeAbsolutePath(path)
	return os.MkdirAll(path, 0755) == nil // 0755 is the file permission: rwxr-xr-x
}
func Delete(path string) bool {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return false
	}
	return os.RemoveAll(path) == nil
}
func DeleteContents(path string) bool {
	path = internal.MakeAbsolutePath(path)
	if !Delete(path) {
		return false
	}
	if !Create(path) {
		return false
	}
	return true
}
func Rename(path, newName string) bool {
	path = internal.MakeAbsolutePath(path)
	var newPath = ph.New(ph.Folder(path), newName)

	if !IsExisting(path) || IsExisting(newPath) {
		return false
	}

	return os.Rename(path, newPath) == nil
}
func MoveContents(fromPath, toPath string) bool {
	fromPath = internal.MakeAbsolutePath(fromPath)
	toPath = internal.MakeAbsolutePath(toPath)
	if !IsExisting(fromPath) || !IsExisting(toPath) {
		return false
	}

	var files = Files(fromPath)
	for _, file := range files {
		var srcPath = ph.New(fromPath, file)
		var destPath = ph.New(toPath, file)
		var err = os.MkdirAll(ph.Folder(destPath), 0755)
		if err != nil {
			return false
		}

		var err2 = os.Rename(srcPath, destPath)
		if err2 == nil {
			continue
		}

		var srcFile, err3 = os.Open(srcPath)
		if err3 != nil {
			return false
		}
		defer srcFile.Close()

		var destFile, err4 = os.Create(destPath)
		if err4 != nil {
			return false
		}
		defer destFile.Close()

		var _, err5 = io.Copy(destFile, srcFile)
		if err5 != nil {
			return false
		}

		var err6 = os.Remove(srcPath)
		if err6 != nil {
			return false
		}
	}

	var folders = Folders(fromPath)
	for _, sub := range folders {
		var srcSub = ph.New(fromPath, sub)
		var destSub = ph.New(toPath, sub)

		var err = os.MkdirAll(destSub, 0755)
		if err != nil {
			return false
		}
		if !MoveContents(srcSub, destSub) {
			return false
		}
	}

	return true
}
func CopyContents(fromPath, toPath string) bool {
	fromPath = internal.MakeAbsolutePath(fromPath)
	toPath = internal.MakeAbsolutePath(toPath)
	if !IsExisting(fromPath) || !IsExisting(toPath) {
		return false
	}

	var files = Files(fromPath)
	for _, file := range files {
		var srcPath = ph.New(fromPath, file)
		var destPath = ph.New(toPath, file)
		var srcFile, err = os.Open(srcPath)
		if err != nil {
			return false
		}

		defer srcFile.Close()

		var err2 = os.MkdirAll(ph.Folder(destPath), 0755)
		if err2 != nil {
			return false
		}

		var destFile, err3 = os.Create(destPath)
		if err3 != nil {
			return false
		}
		defer destFile.Close()

		var _, err4 = io.Copy(destFile, srcFile)
		if err4 != nil {
			return false
		}
	}

	var folders = Folders(fromPath)
	for _, sub := range folders {
		var srcSub = ph.New(fromPath, sub)
		var destSub = ph.New(toPath, sub)
		var err = os.MkdirAll(destSub, 0755)
		if err != nil {
			return false
		}
		if !CopyContents(srcSub, destSub) {
			return false
		}
	}

	return true
}

func Content(path string) []string {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return []string{}
	}

	var entries, err = os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}
func Files(path string) []string {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return []string{}
	}

	var entries, err = os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}
func Folders(path string) []string {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return []string{}
	}

	var entries, err = os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var folders []string
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}
	return folders
}
