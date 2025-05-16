package auth

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/browser"
	"github.com/Phillezi/kthcloud-cli/pkg/session"
	"github.com/Phillezi/kthcloud-cli/pkg/web"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

func (c *Client) Login() (sess *session.Session, err error) {
	keycloakURL, err := c.constructKeycloakURL()
	if err != nil {
		// log err
		return nil, err
	}

	url, err := url.Parse(c.redirectHost)
	if err != nil {
		// log err
		return nil, err
	}
	serverPort := func() string {
		if port := url.Port(); port != "" {
			return port
		}
		if url.Scheme == "https" {
			return "443"
		}
		return "80"
	}()

	serverAddr := ":" + serverPort

	serverCtx, cancelServerCtx := context.WithCancel(c.ctx)
	defer cancelServerCtx()
	sessionChannel := make(chan *session.Session)
	server := web.New(web.ServerOpts{
		Address:         &serverAddr,
		KeycloakURL:     &keycloakURL,
		RedirectHost:    &c.redirectHost,
		RedirectPath:    &c.redirectPath,
		SessionChannel:  &sessionChannel,
		FetchOAuthToken: c.fetchOAuthToken,
	}).WithContext(serverCtx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(); err != nil {
			// log err here
			logrus.Errorln(err)
		}
	}()
	defer wg.Wait()

	err = browser.Open(c.redirectHost)
	if err != nil {
		// log err

		logrus.Errorln(err)
		return nil, err
	}

	fmt.Println("Waiting for login...")

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()
	defer s.Stop()

	select {
	case sess = <-sessionChannel:
		if sess != nil {
			logrus.Debug("recv token")
			c.session = sess
			c.client.SetAuthToken(c.session.Token.AccessToken)
			if err := c.session.Save(c.sessionPath); err != nil {
				logrus.Errorln("failed to save session", err)
			}
		} else {
			return nil, fmt.Errorf("server was closed without sending token")
		}
	case <-c.ctx.Done():
		return nil, fmt.Errorf("interrupted")
	}

	return sess, nil
}
