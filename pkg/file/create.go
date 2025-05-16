package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Create(basePath, fileName, content string) error {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create folder: %v", err)
		}
	}

	filePath := filepath.Join(basePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %v", err)
	}

	return nil
}
