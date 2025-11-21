package app

import "errors"

var (
	ErrLogin = errors.New("error on login")

	ErrUnsupportedSessionType     = errors.New("error session doesnt support")
	ErrSessionDoesntSupportLogin  = errors.Join(ErrUnsupportedSessionType, errors.New("login"))
	ErrSessionDoesntSupportLogout = errors.Join(ErrUnsupportedSessionType, errors.New("logout"))
)
