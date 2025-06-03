package folder

import "os"

func Exists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	return err == nil && info.IsDir()
}

func Create(folderPath string) bool {
	err := os.MkdirAll(folderPath, 0755) // 0755 is the file permission: rwxr-xr-x
	return err == nil
}

func Delete(folderPath string) bool {
	if !Exists(folderPath) {
		return false
	}
	err := os.RemoveAll(folderPath)
	return err == nil
}
