package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kthcloud/cli/pkg/browser"
	"github.com/kthcloud/cli/pkg/session"
	"go.uber.org/zap"
)

func (a *App) Login() error {
	if m, ok := a.session.(session.Manager); ok {
		s, err := m.GetSession(a.sessionKey)
		if err != nil && !errors.Is(err, session.ErrNotFound) {
			return errors.Join(err, ErrLogin)
		}

		if s == nil || !s.IsValid() {
			if err := a.login(m); err != nil {
				return errors.Join(err, ErrLogin)
			}
		} else {
			a.l.Debug("re-using valid session")
		}

		return nil
	} else {
		return ErrSessionDoesntSupportLogin
	}
}

func (a *App) login(manager session.Manager) error {
	errCh := make(chan error, 1)

	go func() {
		if err := a.loginServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()
	defer a.loginServer.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	if err := browser.Open(a.loginServer.Url()); err != nil {
		if errors.Is(err, browser.ErrUnsupportedPlatform) {
			fmt.Printf("Open \033]8;;%s\a%s\033]8;;\a in the browser\n", a.loginServer.Url(), a.loginServer.Url())
		} else {
			a.l.Error("Failed to open the url", zap.String("url", a.loginServer.Url()), zap.Error(err))
			return err
		}
	}

	select {
	case <-a.ctx.Done():
	case err := <-errCh:
		a.l.Error("Failed to start login server", zap.Error(err))
		return err
	case tok := <-a.loginServer.Token():
		manager.SaveSession(a.sessionKey, session.FromOAuth2Token(tok))
	}

	return nil
}

func (a *App) Logout() error {
	if m, ok := a.session.(session.Manager); ok {
		return m.DeleteSession(a.sessionKey)
	} else {
		return ErrSessionDoesntSupportLogout
	}
}
