package update

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func updateInteractively() (bool, error) {
	execPath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %v", err)
	}

	releases, err := GetReleases()
	if err != nil {
		return false, fmt.Errorf("failed to get latest release: %v", err)
	}

	selectedRelease, err := SelectReleaseInteractively(releases)
	if err != nil {
		return false, err
	}

	resp, err := PromptYesNo("Are you sure?")
	if err != nil {
		return false, err
	}
	if !resp {
		log.Info("ok, stopping version change...")
		return false, nil
	}

	newBinaryPath := filepath.Join(os.TempDir(), "kthcloud-cli-new")
	downloadURL, err := FindBinaryForCurrentPlatform(selectedRelease)
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

func Update(skipPrompt bool, interactive bool) (bool, error) {
	if interactive {
		return updateInteractively()
	}
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
		log.Infoln("No newer versions found")
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
