package session

import "errors"

var (
	ErrTokenExpired = errors.New("token expired and no refresher available")

	ErrNoRefresher = errors.New("no refresher available")

	ErrNotImplemented = errors.New("not implemented")

	ErrNotFound = errors.New("session not found")

	ErrLoginRequired = errors.New("login is required")

	ErrMiddlewareOnNilReq = errors.New("auth middleware called on nil http.Request")
)
