package filebrowser

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/browser"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/browserutils/kooky"
	"github.com/sirupsen/logrus"
)

func (c *Client) loadCookies() (bool, error) {
	storageURL, err := url.Parse(c.filebrowserURL)
	if err != nil {
		return false, err
	}

	kcURL, err := url.Parse(c.keycloakBaseURL)
	if err != nil {
		return false, err
	}

	commonDomain, err := util.GetCommonDomain(storageURL.Host, kcURL.Host)
	if err != nil {
		return false, err
	}

	try := 0

	for {
		allCookies := kooky.AllCookies(kooky.Valid, kooky.DomainHasSuffix(commonDomain))

		iamCookies := make([]*http.Cookie, 0)
		storageCookies := make([]*http.Cookie, 0)

		for _, cookie := range allCookies {
			cookie.Value = strings.Trim(cookie.Value, "\"")
			if cookie.Domain == storageURL.Host {
				storageCookies = append(storageCookies, &cookie.Cookie)
			}
			if cookie.Domain == kcURL.Host {
				iamCookies = append(iamCookies, &cookie.Cookie)
			}
		}
		c.client.Jar.SetCookies(storageURL, storageCookies)
		c.client.Jar.SetCookies(kcURL, iamCookies)

		if len(iamCookies) == 0 || len(storageCookies) == 0 {
			logrus.Warn("no cookies from keycloak, try to log in")
			if try >= 3 {
				return false, errors.New("could not get required cookies")
			}

			browser.Open(c.filebrowserURL)
			logrus.Info("retrying in 5s...")
			time.Sleep(5 * time.Second)
			try++
			continue
		} else {
			logrus.Debugln("cookies found")
			logrus.Debugln(iamCookies)
			logrus.Debugln(storageCookies)
		}
		break
	}

	return true, nil
}
