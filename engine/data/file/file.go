package file

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"pure-kit/engine/data/path"
)

func Exists(filePath string) bool {
	var info, err = os.Stat(filePath)
	return err == nil && !info.IsDir()
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

func LoadBytes(filePath string) []byte {
	var data, err = os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file '%s': %v\n", filePath, err)
		return []byte{}
	}
	return data
}
func LoadText(filePath string) string {
	return string(LoadBytes(filePath))
}

func SaveBytes(filePath string, content []byte) bool {
	var err = os.WriteFile(filePath, content, 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}
func SaveText(filePath, content string) bool {
	return SaveBytes(filePath, []byte(content))
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

func Delete(filePath string) bool {
	if !Exists(filePath) {
		return false
	}
	var err = os.Remove(filePath)
	return err == nil
}
func Rename(filePath, newName string) bool {
	var newFilePath = path.New(path.Folder(filePath), newName)

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
	return Rename(filePath, path.New(toFolderPath, filePath))
}
func Copy(filePath, toFolderPath string) bool {
	var srcFile, err = os.Open(filePath)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	var destFile, err2 = os.Create(path.New(toFolderPath, filePath))
	if err2 != nil {
		return false
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err == nil
}

func Compress(data []byte) []byte {
	var buf bytes.Buffer
	var gw = gzip.NewWriter(&buf)
	var _, err = gw.Write(data)

	if err != nil {
		return data
	}

	if err := gw.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}
func Decompress(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = gzip.NewReader(buf)

	if err != nil {
		return data
	}
	defer gr.Close()

	var result, err2 = io.ReadAll(gr)
	if err2 != nil {
		return data
	}
	return result

}
