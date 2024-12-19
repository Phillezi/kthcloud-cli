package cicd

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(basePath, fileName, content string) error {
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

func ReadFile(basePath, fileName string) (string, error) {
	filePath := filepath.Join(basePath, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(content), nil
}

func FileExists(basePath, fileName string) bool {
	filePath := filepath.Join(basePath, fileName)

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
