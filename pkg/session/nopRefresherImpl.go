package session

type NopRefresherImpl struct {
}

func (NopRefresherImpl) Refresh(refreshToken string) (*Session, error) {
	return nil, ErrNoRefresher
}
