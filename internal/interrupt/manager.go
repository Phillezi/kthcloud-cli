package interrupt

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	instance *Manager
	once     sync.Once
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	shutdowns []func()
	mu        sync.Mutex
}

func GetInstance() *Manager {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		instance = &Manager{
			ctx:    ctx,
			cancel: cancel,
		}
		instance.listenForSignals()
	})
	return instance
}

func (im *Manager) listenForSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		zap.L().Info("Received shutdown signal", zap.String("signal", sig.String()))
		im.Shutdown()
	}()
}

func (im *Manager) Context() context.Context {
	return im.ctx
}

func (im *Manager) AddShutdownHook(hook func()) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.shutdowns = append(im.shutdowns, hook)
}

func (im *Manager) Shutdown() {
	im.cancel()
	im.mu.Lock()
	defer im.mu.Unlock()

	zap.L().Info("Shutting down, executing hooks...")
	for _, hook := range im.shutdowns {
		hook()
	}
	zap.L().Info("Shutdown complete.")
}

func (im *Manager) Wait(timeout time.Duration) {
	select {
	case <-im.ctx.Done():
		zap.L().Info("Context canceled, exiting.")
	case <-time.After(timeout):
		zap.L().Warn("Timeout reached, forcing exit.")
	}
}
