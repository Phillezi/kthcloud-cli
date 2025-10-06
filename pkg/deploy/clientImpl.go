package deploy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/kthcloud/cli/pkg/core"
)

type ClientImpl struct {
	ctx context.Context

	client *core.ClientImpl
}

func New(ctx context.Context, opts ...Option) *ClientImpl {
	c := ClientImpl{ctx: ctx}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *ClientImpl) _() {
	c.client.Do(&http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "https",
			Host:   "cloud.cbh.kth.se:443",
			Path:   "/",
		},
		Header: http.Header{
			"Content-Type":  []string{"application/json"},
			"Authorization": []string{""},
		},
	})
}
