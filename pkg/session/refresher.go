package session

type Refresher interface {
	Refresh(refreshToken string) (*Session, error)
}
