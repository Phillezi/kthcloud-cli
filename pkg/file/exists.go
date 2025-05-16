package file

import (
	"os"
	"path/filepath"
)

func Exists(basePath, fileName string) bool {
	filePath := filepath.Join(basePath, fileName)

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
