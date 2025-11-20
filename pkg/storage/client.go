package storage

import (
	"io"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
	IsDir   bool
}

type Client interface {
	MkdirAll(path string) error
	WriteFile(path string, r io.Reader) error
	ReadFile(path string) (io.ReadCloser, error)
	DeleteFile(path string) error
	ListDir(path string) ([]FileInfo, error)
	RemoveAll(path string) error
}
