package deploy

import "errors"

var (
	ErrNilResponse     = errors.New("response is nil")
	ErrInvalidResponse = errors.New("invalid response object")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrServerError     = errors.New("internal server error")
	ErrUnexpected      = errors.New("unexpected response status")
)
