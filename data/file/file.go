/*
Wraps some essential Operating System/Input-Output (OS/IO) file features with helper functions
to make them more digestible and clarify their API.
*/
package file

import (
	"os"
	"pure-game-kit/debug"
	"pure-game-kit/utility/text"
)

func Exists(path string) bool {
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
	if !Exists(path) {
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
	if !Exists(path) {
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
	return text.Remove(string(LoadBytes(path)), "\r") // FUCK windows pt1
}

func SaveBytes(path string, content []byte) bool {
	var err = os.WriteFile(path, content, 0644)
	if err != nil {
		debug.LogError("Failed to save file: \"", path, "\"\n", err)
	}
	return err == nil
}
func SaveText(path, content string) bool {
	return SaveBytes(path, []byte(text.Remove(content, "\r"))) // FUCK windows pt2
}
func SaveTextAppend(path string, content string) bool {
	if !Exists(path) {
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
