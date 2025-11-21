package logs

import (
	"context"
	"time"

	"github.com/kthcloud/cli/pkg/session"
	"go.uber.org/zap"
)

type Option func(c *ClientImpl)

func WithContext(ctx context.Context) Option {
	return func(c *ClientImpl) {
		c.ctx = ctx
	}
}

func WithSession(session session.Auth) Option {
	return func(c *ClientImpl) {
		c.session = session
	}
}

func WithAPIURL(apiURL string) Option {
	return func(c *ClientImpl) {
		c.apiURL = apiURL
	}
}

func WithInitialBatchDelay(initialBatchDelay time.Duration) Option {
	return func(c *ClientImpl) {
		c.initialBatchDelay = initialBatchDelay
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(c *ClientImpl) {
		c.logger = logger
	}
}
