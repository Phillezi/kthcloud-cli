package storage

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/client"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/models/compose"
	"github.com/sirupsen/logrus"
)

type FileType int

const (
	Nonexistent FileType = iota
	File
	Dir
)

func (ft FileType) String() string {
	return [...]string{"Nonexistent", "File", "Dir"}[ft]
}

func CreateVolumes(c *client.Client, composeInstance *compose.Compose) (string, error) {
	projectDir := composeInstance.Hash()
	user, err := c.User()
	if err != nil {
		return "", err
	}
	storageURL := user.StorageURL
	if storageURL == nil {
		return "", errors.New("user does not have a storageURL")
	}

	if !c.HasValidSession() {
		// TODO: Bring back API key support
		/*if session.ApiKey != nil {
			return "", errors.New("volume creation requires being logged in, api keys do not work for this, please log in")
		} else {
			return "", errors.New("volume creation requires being logged in, please log in")
		}*/
		return "", errors.New("volume creation requires being logged in, please log in")
	}

	isAuth, err := c.StorageAuth()
	if err != nil {
		return "", err
	}
	if !isAuth {
		return "", errors.New("Not authenticated on storage Url" + *user.StorageURL)
	} else {
		logrus.Info("yayyy")
	}

	created, err := c.StorageCreateDir(projectDir)
	if err != nil {
		logrus.Info("dir")
		return "", err
	}
	if !created {
		return "", errors.New("did not create the dir: " + projectDir)
	}

	volumes, err := checkLocalPaths(composeInstance)
	if err != nil {
		return "", err
	}

	for filePath, fileType := range volumes {
		err = handlePath(filePath, fileType, c, projectDir)
		if err != nil {
			return "", err
		}
	}

	return projectDir, nil
}

func checkLocalPaths(composeInstance *compose.Compose) (map[string]FileType, error) {
	pathsStatus := make(map[string]FileType)

	for _, service := range composeInstance.Services {
		for _, volume := range service.Volumes {
			// Split the volume by ':'
			paths := strings.Split(volume, ":")
			if len(paths) == 0 {
				continue
			}

			// Check the first part of the volume (local path)
			localPath := paths[0]

			// Determine if it's a file, directory, or nonexistent
			fileInfo, err := os.Stat(localPath)
			if os.IsNotExist(err) {
				pathsStatus[localPath] = Nonexistent
			} else if err != nil {
				return nil, err
			} else if fileInfo.IsDir() {
				pathsStatus[localPath] = Dir
			} else {
				pathsStatus[localPath] = File
			}
		}
	}

	return pathsStatus, nil
}

func handlePath(filePath string, fileType FileType, c *client.Client, projectDir string) error {
	serverPath := path.Join(projectDir, filePath)
	if fileType == Dir || fileType == Nonexistent {
		created, err := c.StorageCreateDir(serverPath)
		if err != nil {
			return err
		}
		if !created {
			return errors.New("did not create the dir: " + serverPath)
		}
		if fileType == Dir {
			children, err := os.ReadDir(filePath)
			if err != nil {
				return err
			}

			for _, child := range children {
				child.Name()
				childFT := File
				if child.IsDir() {
					childFT = Dir
				}
				childPath := path.Join(filePath, child.Name())
				err = handlePath(childPath, childFT, c, projectDir)
				if err != nil {
					return err
				}
			}

		}
	} else if fileType == File {
		fileName := path.Base(filePath)

		logrus.Info("fileName: " + fileName)
		logrus.Info("filePath: " + filePath)
		logrus.Info("serverPath: " + serverPath)

		filecontent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		created, err := c.StorageCreateFile(serverPath, filecontent)
		if err != nil {
			return err
		}
		if !created {
			return errors.New("did not create the file: " + serverPath)
		}
	}
	return nil
}
