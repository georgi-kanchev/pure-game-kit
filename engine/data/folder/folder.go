package folder

import (
	"io"
	"os"
	"path/filepath"
)

func PathOfExecutable() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	execDir := filepath.Dir(execPath)
	return execDir
}
func Exists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	return err == nil && info.IsDir()
}
func IsEmpty(folderPath string) bool {
	return len(ListContents(folderPath)) == 0
}
func ByteSize(folderPath string) int64 {
	var totalSize int64 = 0

	_ = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable files
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize
}
func TimeOfLastEdit(folderPath string) (year, month, day, minute int) {
	if !Exists(folderPath) {
		return 0, 0, 0, 0
	}

	info, err := os.Stat(folderPath)
	if err != nil {
		return 0, 0, 0, 0
	}

	t := info.ModTime()
	year = t.Year()
	month = int(t.Month()) // time.Month is 1-based already
	day = t.Day()          // day of the month
	minute = t.Hour()*60 + t.Minute()
	return
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
func DeleteContents(folderPath string) bool {
	if !Delete(folderPath) {
		return false
	}
	if !Create(folderPath) {
		return false
	}
	return true
}
func Rename(folderPath string, newName string) bool {
	newPath := filepath.Join(filepath.Dir(folderPath), newName)

	if !Exists(folderPath) || Exists(newPath) {
		return false
	}

	err := os.Rename(folderPath, newPath)
	return err == nil
}
func MoveContents(fromFolderPath string, toFolderPath string) bool {
	if !Exists(fromFolderPath) || !Exists(toFolderPath) {
		return false
	}

	err := filepath.WalkDir(fromFolderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == fromFolderPath {
			return nil // Skip root
		}

		relPath, err := filepath.Rel(fromFolderPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(toFolderPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Ensure parent folder exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Try rename first
		if err := os.Rename(path, targetPath); err == nil {
			return nil
		}

		// Fallback: copy + delete
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}

		return os.Remove(path)
	})

	return err == nil
}
func CopyContents(fromFolderPath, toFolderPath string) bool {
	if !Exists(fromFolderPath) || !Exists(toFolderPath) {
		return false
	}

	err := filepath.WalkDir(fromFolderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(fromFolderPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(toFolderPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
	return err == nil
}

func ListContents(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return []string{}
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}
func ListFiles(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	entries, err := os.ReadDir(folderPath)
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
func ListFolders(folderPath string) []string {
	if !Exists(folderPath) {
		return []string{}
	}

	entries, err := os.ReadDir(folderPath)
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
