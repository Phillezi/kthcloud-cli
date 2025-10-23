package session

import (
	"os"
)

type SessionStore interface {
	Set(key string, session *Session) error
	Get(key string) (*Session, error)
	Delete(key string) error
	Clear() error
}

// NewSessionStoreFactory returns a SessionStore depending on the OS and availability.
// service: identifier for the keyring (or app name)
// fallbackDir: path to store fallback encrypted files
// fallbackKey: 16/24/32-byte AES key for file encryption
func NewSessionStore(service, fallbackDir string) (SessionStore, error) {
	// Ensure fallback directory exists
	if err := os.MkdirAll(fallbackDir, 0700); err != nil {
		return nil, err
	}

	if KeyringAvailable() {
		return NewKeyringStoreImpl(service), nil
	}

	return NewFileStoreImpl(fallbackDir), nil
}
