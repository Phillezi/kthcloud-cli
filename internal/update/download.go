package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
)

type ProgressReader struct {
	Reader       io.Reader
	TotalSize    int
	CurrentBytes int
	Spinner      *spinner.Spinner
}

func DownloadBinary(url, filename string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Downloading " + filename
	s.Start()

	resp, err := http.Get(url)
	if err != nil {
		s.Stop()
		return fmt.Errorf("error downloading binary: %v", err)
	}
	defer resp.Body.Close()

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		s.Stop()
		return fmt.Errorf("unable to get content length: %v", err)
	}

	out, err := os.Create(filename)
	if err != nil {
		s.Stop()
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	progressReader := &ProgressReader{
		Reader:       resp.Body,
		TotalSize:    contentLength,
		Spinner:      s,
		CurrentBytes: 0,
	}

	_, err = io.Copy(out, progressReader)
	if err != nil {
		s.Stop()
		return fmt.Errorf("error saving binary: %v", err)
	}

	s.Stop()
	fmt.Println("Download complete!")

	return nil
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.CurrentBytes += n

	percent := float64(pr.CurrentBytes) / float64(pr.TotalSize) * 100
	pr.Spinner.Suffix = fmt.Sprintf(" Downloading binary... %.2f%%", percent)

	return n, err
}

func FindBinaryForCurrentPlatform(release *GitHubRelease) (string, error) {
	platform := runtime.GOOS
	arch := runtime.GOARCH
	var binaryName string

	if platform == "windows" {
		binaryName = fmt.Sprintf("kthcloud_%s_%s.exe", arch, platform)
	} else {
		binaryName = fmt.Sprintf("kthcloud_%s_%s", arch, platform)
	}

	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no binary found for platform %s-%s", platform, arch)
}
