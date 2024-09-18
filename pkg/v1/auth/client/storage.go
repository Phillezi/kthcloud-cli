package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
	"github.com/sirupsen/logrus"
)

func (c *Client) StorageAuth() (bool, error) {
	user, err := c.User()
	if err != nil {
		return false, err
	}
	// , kooky.NameHasPrefix("_oauth2_proxy_")
	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("cloud.cbh.kth.se"))
	httpCookies := make([]*http.Cookie, 0)
	auth := ""

	storageURL, err := url.Parse(*user.StorageURL)
	if err != nil {
		return false, err
	}

	for _, cookie := range cookies {
		if cookie.Name == "auth" && cookie.Domain == storageURL.Host {
			auth = cookie.Value
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

	logrus.Infoln("header: ", resp.Header(), " resp: ", resp.String(), " code: ", resp.StatusCode())

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
	url := fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath)
	logrus.Infoln("URL: ", url, " Body: ", req.Body)

	log.Println("Request Headers:")
	for key, value := range req.Header {
		log.Printf("%s: %s\n", key, value)
	}

	resp, err := req.Post(url)
	if err != nil {
		logrus.Errorln("Create initial file")
		return false, err
	}

	if resp.IsError() {
		logrus.Errorln("Create initial file")
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	// 2. Get upload offset
	resp, err = req.Head(fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath))
	if err != nil {
		logrus.Errorln("Get upload offset")
		return false, err
	}

	if resp.IsError() {
		logrus.Errorln("Get upload offset")
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	// Parse the Upload-Offset from the response header
	uploadOffset := resp.Header().Get("Upload-Offset")
	if uploadOffset == "" {
		logrus.Errorln("Missing Upload-Offset header")
		return false, fmt.Errorf("missing Upload-Offset header")
	}

	offsetValue, err := strconv.ParseInt(uploadOffset, 10, 64)
	if err != nil {
		logrus.Errorln("Invalid Upload-Offset value")
		return false, err
	}

	// 3. Upload contents
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", strconv.FormatInt(offsetValue, 10))
	req.Body = content
	resp, err = req.Patch(fmt.Sprintf("%s/api/tus/%s?override=false", *user.StorageURL, filePath))
	if err != nil {
		logrus.Errorln("Upload contents")
		return false, err
	}

	if resp.IsError() {
		logrus.Errorln("Upload contents")
		return false, fmt.Errorf("request error: %s", resp.Status())
	}

	return true, nil
}
