package session_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kthcloud/cli/pkg/session"
	"golang.org/x/oauth2"
)

func TestNewSessionStoreFactory(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "session_test")
	defer os.RemoveAll(tmpDir)

	store, err := session.NewSessionStore("test-service", tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if store == nil {
		t.Fatalf("expected a store, got nil")
	}

	switch s := store.(type) {
	case *session.FileStoreImpl:
		t.Log("Using FileStoreImpl fallback")

		info, err := os.Stat(tmpDir)
		if err != nil {
			t.Fatalf("expected fallback directory to exist, got error: %v", err)
		}
		if !info.IsDir() {
			t.Fatalf("expected fallback path to be a directory")
		}

		testSession := &session.Session{Token: &oauth2.Token{
			AccessToken: "abc123",
		}}
		key := "testsession"

		if err := s.Set(key, testSession); err != nil {
			t.Fatalf("Set() failed: %v", err)
		}

		got, err := s.Get(key)
		if err != nil {
			t.Fatalf("Get() failed: %v", err)
		}
		if got.AccessToken() != testSession.AccessToken() {
			t.Fatalf("expected AccessToken %s, got %s", testSession.AccessToken(), got.AccessToken())
		}

		if err := s.Delete(key); err != nil {
			t.Fatalf("Delete() failed: %v", err)
		}

		_, err = s.Get(key)
		if err == nil {
			t.Fatalf("expected error getting deleted session, got nil")
		}

	case *session.KeyringStoreImpl:
		t.Log("Using KeyringStoreImpl")

		testSession := &session.Session{Token: &oauth2.Token{
			AccessToken: "abc123",
		}}
		key := "testsession"

		if err := s.Set(key, testSession); err != nil {
			t.Fatalf("Set() failed: %v", err)
		}

		got, err := s.Get(key)
		if err != nil {
			t.Fatalf("Get() failed: %v", err)
		}
		if got.AccessToken() != testSession.AccessToken() {
			t.Fatalf("expected AccessToken %s, got %s", testSession.AccessToken(), got.AccessToken())
		}

		if err := s.Delete(key); err != nil {
			t.Fatalf("Delete() failed: %v", err)
		}

		_, err = s.Get(key)
		if err == nil {
			t.Fatalf("expected error getting deleted session, got nil")
		}

	default:
		t.Fatalf("unexpected store type %T", store)
	}
}
