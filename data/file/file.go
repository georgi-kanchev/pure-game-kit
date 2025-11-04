package file

import (
	"os"
	"pure-game-kit/debug"
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
	if IsExisting(path) {
		debug.LogError("Failed to find file: \"", path, "\"")
		return []byte{}
	}

	var data, err = os.ReadFile(path)
	if err != nil {
		debug.LogError("Failed to load file: \"", path, "\"\n", err)
		return []byte{}
	}
	return data
}
func LoadText(path string) string {
	return string(LoadBytes(path))
}

func SaveBytes(path string, content []byte) bool {
	var err = os.WriteFile(path, content, 0644)
	if err != nil {
		debug.LogError("Failed to save file: \"", path, "\"\n", err)
	}
	return err == nil
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
	var _, err2 = file.WriteString(content)
	defer file.Close()

	if err != nil {
		debug.LogError("Failed to open file: \"", path, "\"\n", err)
		return false
	}
	if err2 != nil {
		debug.LogError("Failed to append file: \"", path, "\"\n", err2)
		return false
	}

	return true
}

/*
func Delete(path string) bool {
	if !IsExisting(path) {
		debug.LogError("Failed to find file: \"", path, "\"")
		return false
	}
	var err = os.Remove(path)
	if err != nil {
		debug.LogError("Failed to delete file: \"", path, "\"\n", err)
	}

	return err == nil
}
func Rename(path, newName string) bool {
	if !IsExisting(path) {
		debug.LogError("Failed to find file: \"", path, "\"")
		return false
	}

	var newpath = ph.New(ph.Folder(path), newName)
	var err = os.Rename(path, newpath)
	if err != nil {
		debug.LogError("Failed to rename file: \"", path, "\"\n", err)
	}
	return err == nil
}
func Move(path, toFolderPath string) bool {
	if !IsExisting(path) {
		debug.LogError("Failed to find file: \"", path, "\"")
		return false
	}

	var info, err = os.Stat(toFolderPath)
	if err != nil || !info.IsDir() {
		debug.LogError("Failed to find target folder: \"", path, "\"\n", err)
		return false
	}

	return Rename(path, ph.New(toFolderPath, path))
}
func Duplicate(path, toFolderPath string) bool {
	if !IsExisting(path) {
		debug.LogError("Failed to find file: \"", path, "\"")
		return false
	}

	var info, err = os.Stat(toFolderPath)
	if err != nil || !info.IsDir() {
		debug.LogError("Failed to find target folder: \"", toFolderPath, "\"\n", err)
		return false
	}

	var from, err2 = os.Open(path)
	if err2 != nil {
		debug.LogError("Failed to open file: \"", path, "\"\n", err2)
		return false
	}
	defer from.Close()

	var targetPath = ph.New(toFolderPath, path)
	var to, err3 = os.Create(targetPath)
	if err3 != nil {
		debug.LogError("Failed to create file: \"", targetPath, "\"\n", err3)
		return false
	}
	defer to.Close()

	var _, err4 = io.Copy(to, from)
	if err4 != nil {
		debug.LogError("Failed to copy file: \"", path, "\" -> \"", targetPath, "\"\n", err4)
		return false
	}

	return true
}
*/
