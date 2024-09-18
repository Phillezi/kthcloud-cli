package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

func (c *Client) StorageAuth() (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}

	storageURL, err := url.Parse(*user.StorageURL)
	if err != nil {
		return false, err
	}

	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("cloud.cbh.kth.se"))
	httpCookies := make([]*http.Cookie, 0)
	auth := ""

	for _, cookie := range cookies {
		if cookie.Name == "auth" && cookie.Domain == storageURL.Host {
			auth = cookie.Value
		}
		if strings.Contains(cookie.Value, `"`) {
			continue
		}
		httpCookies = append(httpCookies, &http.Cookie{
			Name:        cookie.Name,
			Value:       cookie.Value,
			Quoted:      cookie.Quoted,
			Domain:      cookie.Domain,
			Path:        cookie.Path,
			Expires:     cookie.Expires,
			RawExpires:  cookie.RawExpires,
			Secure:      cookie.Secure,
			HttpOnly:    cookie.HttpOnly,
			SameSite:    cookie.SameSite,
			MaxAge:      cookie.MaxAge,
			Partitioned: cookie.Partitioned,
			Raw:         cookie.Raw,
			Unparsed:    cookie.Unparsed,
		})
	}
	c.cookies = httpCookies
	c.client.Header.Set("X-Auth", auth)

	c.client.Cookies = httpCookies

	req := c.client.R()

	resp, err := req.Get(fmt.Sprintf("%s/oauth2/auth", *user.StorageURL))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	if auth == "" {
		req.Body = "{\"username\": \"\", \"password\": \"\", \"recaptcha\": \"\"}"
		resp, err := req.Post(fmt.Sprintf("%s/api/login", *user.StorageURL))
		if err != nil {
			return false, err
		}

		if resp.IsError() {
			return false, fmt.Errorf("request error: %s", resp.Status())
		}

		c.client.Header.Set("X-Auth", resp.String())
	}

	return true, nil
}

func (c *Client) StorageCreateDir(filePath string) (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	req := c.client.R()

	resp, err := req.Post(fmt.Sprintf("%s/api/resources/%s/?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}

func (c *Client) StorageCreateFile(filePath string, content []byte) (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	req := c.client.R()

	// 1. Create initial file
	resp, err := req.Post(fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	// 2. Get upload offset
	resp, err = req.Head(fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	// Parse the Upload-Offset from the response header
	uploadOffset := resp.Header().Get("Upload-Offset")
	if uploadOffset == "" {
		return false, fmt.Errorf("missing Upload-Offset header")
	}

	offsetValue, err := strconv.ParseInt(uploadOffset, 10, 64)
	if err != nil {
		return false, err
	}

	if len(strings.Trim(resp.String(), " \n")) > 0 {
		return false, errors.New("file not empty")
	}

	// 3. Upload contents
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", strconv.FormatInt(offsetValue, 10))
	req.Body = content
	resp, err = req.Patch(fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath))
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}
