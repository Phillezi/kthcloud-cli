package app

import (
	"context"

	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
)

type App struct {
	ctx context.Context

	session session.Client
	deploy  deploy.Client
}

func New(ctx context.Context, opts ...Option) *App {
	a := App{ctx: ctx}
	for _, opt := range opts {
		opt(&a)
	}
	return &a
}
