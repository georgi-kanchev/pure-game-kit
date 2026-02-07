/*
Wraps some essential Operating System/Input-Output (OS/IO) folder features with helper functions
to make them more digestible and clarify their API.
*/
package folder

import (
	"os"
	ph "pure-game-kit/data/path"
	"pure-game-kit/debug"
)

func Exists(path string) bool {
	var info, err = os.Stat(path)
	return err == nil && info.IsDir()
}
func IsEmpty(path string) bool {
	return len(Content(path, false)) == 0
}
func ByteSize(path string) int64 {
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
	if !Exists(path) {
		debug.LogError("Failed to find folder: \"", path, "\"")
		return 0, 0, 0, 0
	}

	var info, err = os.Stat(path)
	if err != nil {
		debug.LogError("Failed to read folder: \"", path, "\"\n", err)
		return 0, 0, 0, 0
	}

	var t = info.ModTime()
	year = t.Year()
	month = int(t.Month()) // time.Month is 1-based already
	day = t.Day()          // day of the month
	minute = t.Hour()*60 + t.Minute()
	return
}

func Content(path string, includeFullPaths bool) []string {
	if !Exists(path) {
		return []string{}
	}

	var entries, err = os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, entry := range entries {
		var value = entry.Name()
		if includeFullPaths {
			value = ph.New(path, value)
		}

		names = append(names, value)
	}
	return names
}
func Files(path string) []string {
	if !Exists(path) {
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
	if !Exists(path) {
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

func Create(path string) bool {
	var err = os.MkdirAll(path, 0755)
	if err != nil {
		debug.LogError("Failed to create folders: \"", path, "\"")
	}
	return err == nil
}

/*
func Delete(path string) bool {
	if !IsExisting(path) {
		return false
	}
	return os.RemoveAll(path) == nil
}
func DeleteContents(path string) bool {
	if !Delete(path) {
		return false
	}
	if !Create(path) {
		return false
	}
	return true
}
func Rename(path, newName string) bool {
	var newPath = ph.New(ph.Folder(path), newName)

	if !IsExisting(path) || IsExisting(newPath) {
		return false
	}

	return os.Rename(path, newPath) == nil
}
func MoveContents(fromPath, toPath string) bool {
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
*/
