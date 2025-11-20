package filebrowser

import "errors"

var (
	ErrChunkTooBig               = errors.New("chunk too big")
	ErrMissingUploadOffsetHeader = errors.New("missing Upload-Offset header")
)
