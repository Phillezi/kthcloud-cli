package session

import (
	"context"
	"net/http"
)

type APITokenSession string

func (a APITokenSession) AuthMiddleware(ctx context.Context, req *http.Request) error {
	if req == nil {
		return ErrMiddlewareOnNilReq
	}

	req.Header.Set("X-Api-Key", string(a))

	return nil
}
