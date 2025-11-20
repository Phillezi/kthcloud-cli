package logs

import (
	"context"
	"sync"

	"github.com/kthcloud/cli/pkg/deploy"
)

type SSEClient struct {
	ctx context.Context

	Connections []*SSEConnection
	middleware  []deploy.RequestEditorFn
}

func New(conns []*SSEConnection, ctx context.Context, mw ...deploy.RequestEditorFn) *SSEClient {
	return &SSEClient{
		Connections: conns,
		ctx:         ctx,
		middleware:  mw,
	}
}

func (m *SSEClient) Start() {
	var wg sync.WaitGroup

	for _, conn := range m.Connections {
		wg.Go(func() {
			conn.OpenConnection(m.ctx, m.middleware...)
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-m.ctx.Done():
	case <-done:
	}
}
