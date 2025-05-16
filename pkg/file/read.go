package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Read(basePath, fileName string) (string, error) {
	filePath := filepath.Join(basePath, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(content), nil
}
