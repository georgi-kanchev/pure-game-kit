package file

import (
	"io"
	"os"
	ph "pure-game-kit/data/path"
)

func IsExisting(path string) bool {
	var info, err = os.Stat(path)
	return err == nil && !info.IsDir()
}
func ByteSize(path string) int64 {
	var info, err = os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
}
func TimeOfLastEdit(path string) (year, month, day, minute int) {
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

func LoadBytes(path string) []byte {
	var data, err = os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return data
}
func LoadText(path string) string {
	return string(LoadBytes(path))
}

func SaveBytes(path string, content []byte) bool {
	return os.WriteFile(path, content, 0644) == nil // 0644 is the file permission: rw-r--r--
}
func SaveText(path, content string) bool {
	return SaveBytes(path, []byte(content))
}
func SaveTextAppend(path string, content string) bool {
	if !IsExisting(path) {
		SaveText(path, content)
		return true
	}

	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	var _, err2 = file.WriteString(content)
	return err2 == nil
}

func Delete(path string) bool {
	if !IsExisting(path) {
		return false
	}
	var err = os.Remove(path)
	return err == nil
}
func Rename(path, newName string) bool {
	var newpath = ph.New(ph.Folder(path), newName)

	if !IsExisting(path) || IsExisting(newpath) {
		return false
	}

	var err = os.Rename(path, newpath)
	return err == nil
}
func Move(path, toFolderPath string) bool {
	var info, err = os.Stat(toFolderPath)
	var folderExists = err == nil && info.IsDir()
	if !IsExisting(path) || !folderExists {
		return false
	}
	return Rename(path, ph.New(toFolderPath, path))
}
func Copy(path, toFolderPath string) bool {
	var srcFile, err = os.Open(path)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	var destFile, err2 = os.Create(ph.New(toFolderPath, path))
	if err2 != nil {
		return false
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err == nil
}
