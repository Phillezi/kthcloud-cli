package session

type Manager interface {
	// SaveSession stores the session securely.
	SaveSession(key string, session *Session) error

	// GetSession retrieves a session by key. If the session is expired,
	// it will attempt to refresh it automatically.
	GetSession(key string) (*Session, error)

	// DeleteSession removes a session.
	DeleteSession(key string) error

	// Clear removes all sessions.
	Clear() error

	Auth
}
