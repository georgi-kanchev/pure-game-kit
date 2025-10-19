package file

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	ph "pure-game-kit/data/path"
	"pure-game-kit/internal"
)

func IsExisting(path string) bool {
	path = internal.MakeAbsolutePath(path)
	var info, err = os.Stat(path)
	return err == nil && !info.IsDir()
}
func ByteSize(path string) int64 {
	path = internal.MakeAbsolutePath(path)
	var info, err = os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
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

func LoadBytes(path string) []byte {
	path = internal.MakeAbsolutePath(path)
	var data, err = os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file '%s': %v\n", path, err)
		return []byte{}
	}
	return data
}
func LoadText(path string) string {
	path = internal.MakeAbsolutePath(path)
	return string(LoadBytes(path))
}

func SaveBytes(path string, content []byte) bool {
	path = internal.MakeAbsolutePath(path)
	var err = os.WriteFile(path, content, 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}
func SaveText(path, content string) bool {
	return SaveBytes(path, []byte(content))
}
func SaveTextAppend(path string, content string) bool {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
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

func Delete(path string) bool {
	path = internal.MakeAbsolutePath(path)
	if !IsExisting(path) {
		return false
	}
	var err = os.Remove(path)
	return err == nil
}
func Rename(path, newName string) bool {
	path = internal.MakeAbsolutePath(path)
	var newpath = ph.New(ph.Folder(path), newName)

	if !IsExisting(path) || IsExisting(newpath) {
		return false
	}

	var err = os.Rename(path, newpath)
	return err == nil
}
func Move(path, toFolderPath string) bool {
	path = internal.MakeAbsolutePath(path)
	toFolderPath = internal.MakeAbsolutePath(toFolderPath)
	var info, err = os.Stat(toFolderPath)
	var folderExists = err == nil && info.IsDir()
	if !IsExisting(path) || !folderExists {
		return false
	}
	return Rename(path, ph.New(toFolderPath, path))
}
func Copy(path, toFolderPath string) bool {
	path = internal.MakeAbsolutePath(path)
	toFolderPath = internal.MakeAbsolutePath(toFolderPath)
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
