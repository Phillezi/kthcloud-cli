package storage

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/convert"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
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

func CreateVolumes(c *deploy.Client, compose *convert.Wrap) (string, error) {
	projectDir := convert.HashServices(compose.Source.Services)
	user, err := c.User()
	if err != nil {
		return "", err
	}
	storageURL := user.StorageURL
	if storageURL == nil {
		return "", errors.New("user does not have a storageURL")
	}

	isAuth, err := c.Storage().Auth()
	if err != nil {
		return "", err
	}
	if !isAuth {
		return "", errors.New("Not authenticated on storage Url" + *user.StorageURL)
	}

	created, err := c.Storage().CreateDir(projectDir)
	if err != nil {
		return "", err
	}
	if !created {
		return "", errors.New("did not create the dir: " + projectDir)
	}

	volumes, err := checkLocalPaths(compose)
	if err != nil {
		return "", err
	}

	for filePath, fileType := range volumes {
		err = handlePath(filePath, fileType, c, projectDir, compose.Source.WorkingDir)
		if err != nil {
			return "", err
		}
	}

	return projectDir, nil
}

func checkLocalPaths(compose *convert.Wrap) (map[string]FileType, error) {
	pathsStatus := make(map[string]FileType)

	for _, service := range compose.Source.Services {
		for _, volume := range service.Volumes {

			// Determine if it's a file, directory, or nonexistent
			fileInfo, err := os.Stat(volume.Source)
			if os.IsNotExist(err) {
				pathsStatus[volume.Source] = Nonexistent
			} else if err != nil {
				return nil, err
			} else if fileInfo.IsDir() {
				pathsStatus[volume.Source] = Dir
			} else {
				pathsStatus[volume.Source] = File
			}
		}
	}

	return pathsStatus, nil
}

func handlePath(filePath string, fileType FileType, c *deploy.Client, projectDir, composeCWD string) error {
	serverPath := path.Join(projectDir, strings.TrimPrefix(filePath, composeCWD))
	if fileType == Dir || fileType == Nonexistent {
		created, err := c.Storage().CreateDir(serverPath)
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
				err = handlePath(childPath, childFT, c, projectDir, composeCWD)
				if err != nil {
					return err
				}
			}

		}
	} else if fileType == File {
		filecontent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		created, err := c.Storage().UploadFile(serverPath, filecontent)
		if err != nil {
			return err
		}
		if !created {
			return errors.New("did not create the file: " + serverPath)
		}
	}
	return nil
}
