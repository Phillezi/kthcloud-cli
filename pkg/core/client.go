package core

import (
	"context"
	"net/http"

	"github.com/kthcloud/cli/pkg/core/middleware"
	"golang.org/x/net/context/ctxhttp"
)

type ClientImpl struct {
	ctx context.Context

	http *http.Client

	middlewares []middleware.Middleware
}

func New(ctx context.Context, middlewares ...middleware.Middleware) *ClientImpl {
	c := ClientImpl{
		http:        &http.Client{},
		middlewares: middlewares,
	}
	return &c
}

func (c *ClientImpl) Do(req *http.Request) (*http.Response, error) {
	if len(c.middlewares) > 0 {
		for _, mw := range c.middlewares {
			mw(req)
		}
	}
	return ctxhttp.Do(c.ctx, c.http, req)
}
