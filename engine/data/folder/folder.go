package folder

import (
	"io"
	"os"
	"pure-kit/engine/data/path"
)

func Exists(folderPath string) bool {
	var info, err = os.Stat(folderPath)
	return err == nil && info.IsDir()
}
func IsEmpty(folderPath string) bool {
	return len(Content(folderPath)) == 0
}
func ByteSize(folderPath string) int64 {
	var totalSize int64

	var files = Files(folderPath)
	for _, file := range files {
		var info, err = os.Stat(path.New(folderPath, file))
		if err == nil {
			totalSize += info.Size()
		}
	}

	var folders = Folders(folderPath)
	for _, sub := range folders {
		totalSize += ByteSize(path.New(folderPath, sub))
	}

	return totalSize
}
func TimeOfLastEdit(folderPath string) (year, month, day, minute int) {
	if !Exists(folderPath) {
		return 0, 0, 0, 0
	}

	var info, err = os.Stat(folderPath)
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

func Create(folderPath string) bool {
	return os.MkdirAll(folderPath, 0755) == nil // 0755 is the file permission: rwxr-xr-x
}
func Delete(folderPath string) bool {
	if !Exists(folderPath) {
		return false
	}
	return os.RemoveAll(folderPath) == nil
}
func DeleteContents(folderPath string) bool {
	if !Delete(folderPath) {
		return false
	}
	if !Create(folderPath) {
		return false
	}
	return true
}
func Rename(folderPath, newName string) bool {
	var newPath = path.New(path.Folder(folderPath), newName)

	if !Exists(folderPath) || Exists(newPath) {
		return false
	}

	return os.Rename(folderPath, newPath) == nil
}
func MoveContents(fromFolderPath, toFolderPath string) bool {
	if !Exists(fromFolderPath) || !Exists(toFolderPath) {
		return false
	}

	var files = Files(fromFolderPath)
	for _, file := range files {
		var srcPath = path.New(fromFolderPath, file)
		var destPath = path.New(toFolderPath, file)
		var err = os.MkdirAll(path.Folder(destPath), 0755)
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

	var folders = Folders(fromFolderPath)
	for _, sub := range folders {
		var srcSub = path.New(fromFolderPath, sub)
		var destSub = path.New(toFolderPath, sub)

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
func CopyContents(fromFolderPath, toFolderPath string) bool {
	if !Exists(fromFolderPath) || !Exists(toFolderPath) {
		return false
	}

	var files = Files(fromFolderPath)
	for _, file := range files {
		var srcPath = path.New(fromFolderPath, file)
		var destPath = path.New(toFolderPath, file)
		var srcFile, err = os.Open(srcPath)
		if err != nil {
			return false
		}

		defer srcFile.Close()

		var err2 = os.MkdirAll(path.Folder(destPath), 0755)
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

	var folders = Folders(fromFolderPath)
	for _, sub := range folders {
		var srcSub = path.New(fromFolderPath, sub)
		var destSub = path.New(toFolderPath, sub)
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

func Content(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	var entries, err = os.ReadDir(folderPath)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}
func Files(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	var entries, err = os.ReadDir(folderPath)
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
func Folders(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	var entries, err = os.ReadDir(folderPath)
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
