package update

import (
	"fmt"
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer sourceFile.Close()

	err = os.Remove(dst)
	if err != nil {
		return fmt.Errorf("error removing destination file: %v", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	if err = os.Chmod(dst, 0755); err != nil {
		return fmt.Errorf("error setting file permissions: %v", err)
	}

	return nil
}
