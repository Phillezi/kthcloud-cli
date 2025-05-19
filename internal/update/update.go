package update

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/minio/selfupdate"
)

func updateInteractively() (bool, error) {

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

	downloadURL, err := FindBinaryForCurrentPlatform(selectedRelease)
	if err != nil {
		return false, err
	}

	binResp, err := http.Get(downloadURL)
	if err != nil {
		return false, err
	}
	defer binResp.Body.Close()
	err = selfupdate.Apply(binResp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			log.Errorln("Failed to rollback from bad update: %v", rerr)
		}
		return false, err
	}

	return true, nil
}

func Update(skipPrompt bool, interactive bool) (bool, error) {
	if interactive {
		return updateInteractively()
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

	downloadURL, err := FindBinaryForCurrentPlatform(latestRelease)
	if err != nil {
		return false, err
	}

	binResp, err := http.Get(downloadURL)
	if err != nil {
		return false, err
	}
	defer binResp.Body.Close()
	err = selfupdate.Apply(binResp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			log.Errorln("Failed to rollback from bad update: %v", rerr)
		}
		return false, err
	}

	return true, nil
}
