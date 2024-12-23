package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the data to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func EnsureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
		log.Printf("File %s created.", path)
	}
	return nil
}

func FileExists(path string) (bool, error) {
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
