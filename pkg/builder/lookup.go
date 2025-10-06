package builder

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

func GetCICDDeploymentID(contextPath string, onDeplNotCICDConfigured func(baseDir string)) (string, error) {
	if contextPath == "" {
		contextPath = "."
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var fullpath string
	if !filepath.IsAbs(contextPath) {
		fullpath = path.Join(wd, contextPath)
		if fullpath == "" {
			// will probably never happen
			return "", errors.New("fullcontextpath is empty")
		}
	} else {
		fullpath = contextPath
	}
	medadataDir := fullpath + "/.kthcloud"
	callbackCalls := 0 // could be bool but might want retries?

	for {
		deplfileExists, err := util.FileExists(medadataDir + "/DEPLOYMENT")
		if err != nil {
			return "", err
		}

		if deplfileExists {
			break
		}

		if callbackCalls >= 1 {
			return "", errors.New("max callback calls reached but CICD config still doesn't exist")
		}

		if onDeplNotCICDConfigured == nil {
			return "", errors.New("callback function is not defined")
		}

		onDeplNotCICDConfigured(fullpath)
		callbackCalls++
	}
	content, err := os.ReadFile(medadataDir + "/DEPLOYMENT")
	if err != nil {
		return "", err
	}
	return string(content), nil
}
