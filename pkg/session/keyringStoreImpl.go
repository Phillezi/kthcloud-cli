package session

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/zalando/go-keyring"
)

type KeyringStoreImpl struct {
	service string
	cache   sync.Map // map[string]*Session
}

func NewKeyringStoreImpl(service string) *KeyringStoreImpl {
	return &KeyringStoreImpl{service: service}
}

func (s *KeyringStoreImpl) Set(key string, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Set in keyring
	if err := keyring.Set(s.service, key, string(data)); err != nil {
		return err
	}

	// Cache in memory
	s.cache.Store(key, session)
	return nil
}

func (s *KeyringStoreImpl) Get(key string) (*Session, error) {
	// Check cache first
	if sess, ok := s.cache.Load(key); ok {
		if session, ok := sess.(*Session); ok {
			return session, nil
		}
	}

	// Load from keyring
	data, err := keyring.Get(s.service, key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil, errors.Join(err, ErrNotFound)
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	// Cache it
	s.cache.Store(key, &session)
	return &session, nil
}

func (s *KeyringStoreImpl) Delete(key string) error {
	if err := keyring.Delete(s.service, key); err != nil {
		return err
	}

	// Remove from cache
	s.cache.Delete(key)
	return nil
}

func (s *KeyringStoreImpl) Clear() error {
	if err := keyring.DeleteAll(s.service); err != nil {
		return err
	}

	// Clear cache
	s.cache = sync.Map{}
	return nil
}

func KeyringAvailable() bool {
	testKey := "__keyring_test__"
	testValue := "ok"
	err := keyring.Set("test-service", testKey, testValue)
	if err != nil && errors.Is(err, keyring.ErrUnsupportedPlatform) {
		return false
	}
	_ = keyring.Delete("test-service", testKey)
	return true
}
