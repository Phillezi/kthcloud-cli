package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FileStoreImpl struct {
	path string // folder path for storing session files
}

func NewFileStoreImpl(path string) *FileStoreImpl {
	return &FileStoreImpl{
		path: path,
	}
}

func (s *FileStoreImpl) filePath(key string) string {
	return filepath.Join(s.path, key+".session")
}

func (s *FileStoreImpl) Set(key string, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath(key), data, 0600)
}

func (s *FileStoreImpl) Get(key string) (*Session, error) {
	data, err := os.ReadFile(s.filePath(key))
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *FileStoreImpl) Delete(key string) error {
	return os.Remove(s.filePath(key))
}

func (s *FileStoreImpl) Clear() error {
	files, err := os.ReadDir(s.path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".session" {
			_ = os.Remove(filepath.Join(s.path, f.Name()))
		}
	}
	return nil
}
