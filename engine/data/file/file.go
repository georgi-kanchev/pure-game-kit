package file

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

func PathOfExecutable() string {
	var execPath, err = os.Executable()
	if err != nil {
		return ""
	}
	return execPath
}
func Exists(filePath string) bool {
	var info, err = os.Stat(filePath)
	return err == nil && !info.IsDir()
}
func Extension(filePath string) string {
	if !Exists(filePath) {
		return ""
	}
	return filepath.Ext(filePath)
}
func ByteSize(filePath string) int64 {
	var info, err = os.Stat(filePath)
	if err != nil {
		return 0
	}

	return info.Size()
}
func TimeOfLastEdit(filePath string) (year, month, day, minute int) {
	if !Exists(filePath) {
		return 0, 0, 0, 0
	}

	var info, err = os.Stat(filePath)
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

func LoadText(filePath string) string {
	var data, err = os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(data)
}
func LoadBytes(filePath string) []byte {
	var data, err = os.ReadFile(filePath)
	if err != nil {
		return []byte{}
	}
	return data
}

func SaveText(filePath, content string) bool {
	var err = os.WriteFile(filePath, []byte(content), 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}
func SaveTextAppend(filePath string, content string) bool {
	if !Exists(filePath) {
		return false
	}

	var file, err = os.OpenFile("example.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	var _, err2 = file.WriteString(content)
	return err2 == nil
}
func SaveBytes(filePath string, content []byte) bool {
	var err = os.WriteFile(filePath, content, 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}

func Delete(filePath string) bool {
	if !Exists(filePath) {
		return false
	}
	var err = os.Remove(filePath)
	return err == nil
}
func Rename(filePath, newName string) bool {
	var newFilePath = filepath.Join(filepath.Dir(filePath), newName)

	if !Exists(filePath) || Exists(newFilePath) {
		return false
	}

	var err = os.Rename(filePath, newFilePath)
	return err == nil
}
func Move(filePath, toFolderPath string) bool {
	var info, err = os.Stat(toFolderPath)
	var folderExists = err == nil && info.IsDir()
	if !Exists(filePath) || !folderExists {
		return false
	}
	return Rename(filePath, path.Join(toFolderPath, filePath))
}
func Copy(filePath, toFolderPath string) bool {
	var srcFile, err = os.Open(filePath)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	var destFile, err2 = os.Create(path.Join(toFolderPath, filePath))
	if err2 != nil {
		return false
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err == nil
}
