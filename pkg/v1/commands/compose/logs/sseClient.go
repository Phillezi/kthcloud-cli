package logs

import (
	"context"
	"sync"
)

type SSEClient struct {
	Connections []*SSEConnection
	ctx         context.Context
	cancel      context.CancelFunc
}

func New(conns []*SSEConnection, ctx context.Context) *SSEClient {
	ctx, cancel := context.WithCancel(ctx)
	return &SSEClient{
		Connections: conns,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (m *SSEClient) Start() {
	var wg sync.WaitGroup

	for _, conn := range m.Connections {
		wg.Add(1)
		go conn.OpenConnection(&wg, m.ctx)
	}

	wg.Wait()
}

func (m *SSEClient) Stop() {
	m.cancel()
	<-m.ctx.Done()
}
