package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
		Name               string `json:"name"`
	} `json:"assets"`
}

func extractTimestamp(release string) (time.Time, error) {
	// Remove the "release-" prefix
	timestampStr := strings.TrimPrefix(release, "release-")

	// Parse the timestamp as a time.Time object
	layout := "20060102150405" // Format of YYYYMMDDHHMMSS
	return time.Parse(layout, timestampStr)
}

func (r *GitHubRelease) IsNewer() (bool, error) {
	currentReleaseName := viper.GetString("release")

	if currentReleaseName == "" {
		return false, fmt.Errorf("no current release found in configuration")
	}

	currentTimestamp, err := extractTimestamp(currentReleaseName)
	if err != nil {
		return false, fmt.Errorf("error parsing current release timestamp: %v", err)
	}

	latestTimestamp, err := extractTimestamp(r.TagName)
	if err != nil {
		return false, fmt.Errorf("error parsing latest release timestamp: %v", err)
	}

	return latestTimestamp.After(currentTimestamp), nil
}

func GetLatestRelease() (*GitHubRelease, error) {
	url := "https://api.github.com/repos/Phillezi/kthcloud-cli/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching release info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release info, status: %s", resp.Status)
	}

	var release GitHubRelease
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, fmt.Errorf("error decoding release info: %v", err)
	}

	return &release, nil
}

func GetReleases() ([]GitHubRelease, error) {
	url := "https://api.github.com/repos/Phillezi/kthcloud-cli/releases"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching release info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release info, status: %s", resp.Status)
	}

	var releases []GitHubRelease
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, fmt.Errorf("error decoding releases: %v", err)
	}

	return releases, nil
}

func SelectReleaseInteractively(releases []GitHubRelease) (*GitHubRelease, error) {
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases available for selection")
	}

	current := viper.GetString("release")
	startIndex := 0

	releaseItems := make([]string, len(releases))
	for i, r := range releases {
		if r.TagName == current {
			startIndex = i
		}
		releaseItems[i] = fmt.Sprintf("%s", r.TagName)
	}

	prompt := promptui.Select{
		Label:     "Select a Release\nCurrent:" + current,
		Items:     releaseItems,
		CursorPos: startIndex,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to select release: %v", err)
	}

	return &releases[index], nil
}
