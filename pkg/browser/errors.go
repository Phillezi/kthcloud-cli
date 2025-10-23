package browser

import "errors"

var (
	// ErrUnsupportedPlatform is returned when the OS is not supported
	ErrUnsupportedPlatform = errors.New("unsupported platform for opening browser")
)
