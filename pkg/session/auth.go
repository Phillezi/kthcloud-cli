package session

import (
	"context"
	"net/http"
)

type Auth interface {
	// Middleware for outgoing requests that add auth info
	AuthMiddleware(ctx context.Context, req *http.Request) error
}
