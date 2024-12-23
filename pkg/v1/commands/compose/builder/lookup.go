package builder

import (
	"errors"
	"os"
	"path"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

func getCICDDeploymentID(contextPath string, onDeplNotCICDConfigured func(baseDir string)) (string, error) {
	if contextPath == "" {
		contextPath = "."
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fullpath := path.Join(wd, contextPath)
	if fullpath == "" {
		// will probably never happen
		return "", errors.New("fullcontextpath is empty")
	}
	medadataDir := fullpath + "/.kthcloud"
	var deplfileExists bool
	callbackCalls := 0 // could be bool but might want retries?
	for !deplfileExists {
		if callbackCalls >= 1 {
			return "", errors.New("max callback calls reached but cicd config still doesnt exist")
		}
		deplfileExists, err = util.FileExists(medadataDir + "/DEPLOYMENT")
		if err != nil {
			return "", err
		}
		if !deplfileExists {
			// id file doesnt exist call the callback to handle it
			onDeplNotCICDConfigured(fullpath)
			callbackCalls++
		}
	}
	content, err := os.ReadFile(medadataDir + "/DEPLOYMENT")
	if err != nil {
		return "", err
	}
	return string(content), nil
}
