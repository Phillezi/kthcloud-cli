package storageclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
	"github.com/sirupsen/logrus"
)

type Client struct {
	storageURL  string
	keycloakURL string
	client      *http.Client
	token       string
}

var (
	instance *Client
	once     sync.Once
)

func GetInstance(storageURL, keycloakURL string) *Client {
	once.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Error creating cookie jar: %v", err)
		}
		client := &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
		}
		instance = &Client{
			storageURL:  storageURL,
			keycloakURL: keycloakURL,
			client:      client,
			token:       "",
		}
		_, err = instance.loadCookies()
		if err != nil {
			logrus.Fatal(err)
		}
	})
	return instance
}

func (c *Client) loadCookies() (bool, error) {
	storageURL, err := url.Parse(c.storageURL)
	if err != nil {
		return false, err
	}

	kcURL, err := url.Parse(c.keycloakURL)
	if err != nil {
		return false, err
	}

	commonDomain, err := util.GetCommonDomain(storageURL.Host, kcURL.Host)
	if err != nil {
		return false, err
	}

	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix(commonDomain))
	iamCookies := make([]*http.Cookie, 0)
	storageCookies := make([]*http.Cookie, 0)
	for _, cookie := range cookies {
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
	return true, nil
}

func (c *Client) Auth() (bool, error) {
	resp, err := c.client.Get(c.storageURL + "/oauth2/auth")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	if c.token == "" {
		requestBody := strings.NewReader(`{"username": "", "password": "", "recaptcha": ""}`)

		req, err := http.NewRequest(
			"POST",
			c.storageURL+"/api/login",
			requestBody,
		)
		if err != nil {
			return false, err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("request error: %s", resp.Status)
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return false, err
		}

		c.token = string(bodyBytes)
	}

	return true, nil
}

func (c *Client) CreateDir(filePath string) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	endpointURL := fmt.Sprintf("%s/api/resources/%s/?override=false", c.storageURL, filePath)
	req, err := http.NewRequest(
		"POST",
		endpointURL,
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-Auth", c.token)
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	return true, nil
}

func (c *Client) UploadFile(filePath string, content []byte) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	fileCreated, err := c.postFile(filePath)
	if err != nil {
		return fileCreated, err
	}

	uploadOffset, err := c.headFile(filePath)
	if err != nil {
		return fileCreated, err
	}

	fileUploaded, err := c.patchFile(filePath, content, uploadOffset)
	if err != nil {
		return fileUploaded, err
	}

	return true, nil
}

func (c *Client) postFile(filePath string) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.storageURL, filePath)
	req, err := http.NewRequest(
		"POST",
		endpointURL,
		nil,
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-Auth", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	return true, nil
}

func (c *Client) headFile(filePath string) (int, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return 0, err
		}
	}

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.storageURL, filePath)
	req, err := http.NewRequest(
		"HEAD",
		endpointURL,
		nil,
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("X-Auth", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("request error: %s", resp.Status)
	}

	uploadOffsetStr := resp.Header.Get("Upload-Offset")
	if uploadOffsetStr == "" {
		return 0, errors.New("no upload-offset provided in response")
	}
	uploadOffset, err := strconv.Atoi(uploadOffsetStr)
	if err != nil {
		return 0, err
	}

	return uploadOffset, nil
}

func (c *Client) patchFile(filePath string, content []byte, uploadOffset int) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.storageURL, filePath)
	req, err := http.NewRequest(
		"PATCH",
		endpointURL,
		bytes.NewReader(content),
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-Auth", c.token)
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", fmt.Sprintf("%d", uploadOffset))

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	return true, nil
}
