package file

import (
	"os"
)

func Exists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func LoadText(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(data)
}
func Load(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}
	}
	return data
}

func SaveText(filePath string, content string) bool {
	err := os.WriteFile(filePath, []byte(content), 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}
func Save(filePath string, content []byte) bool {
	err := os.WriteFile(filePath, content, 0644) // 0644 is the file permission: rw-r--r--
	return err == nil
}

func Delete(filePath string) bool {
	if !Exists(filePath) {
		return false
	}
	err := os.Remove(filePath)
	return err == nil
}
