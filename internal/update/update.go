package update

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func Update(skipPrompt bool) (bool, error) {
	execPath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %v", err)
	}

	latestRelease, err := GetLatestRelease()
	if err != nil {
		return false, fmt.Errorf("failed to get latest release: %v", err)
	}
	isNewer, err := latestRelease.IsNewer()
	if err != nil {
		return false, err
	}
	if !isNewer {
		log.Infoln("Installation is already latest")
		return false, nil
	}

	if !skipPrompt {
		resp, err := PromptYesNo("Newer version found, do you want to update?")
		if err != nil {
			return false, err
		}
		if !resp {
			log.Info("ok, stopping update...")
			return false, nil
		}
	}

	newBinaryPath := filepath.Join(os.TempDir(), "kthcloud-cli-new")
	downloadURL, err := FindBinaryForCurrentPlatform(latestRelease)
	if err != nil {
		return false, err
	}

	err = DownloadBinary(downloadURL, newBinaryPath)
	if err != nil {
		return false, err
	}

	err = CopyFile(newBinaryPath, execPath)
	if err != nil {
		return false, fmt.Errorf("error replacing binary: %v", err)
	}

	err = os.Remove(newBinaryPath)
	if err != nil {
		return false, fmt.Errorf("error removing temporary file: %v", err)
	}

	return true, nil
}
