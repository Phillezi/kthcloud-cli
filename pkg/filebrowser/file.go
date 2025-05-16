package filebrowser

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

	chunkSize := viper.GetInt("file-upload-chunk-size")
	if chunkSize <= 0 {
		// nginx default size
		// 1mb = 1024kb = 1024^2 bytes
		chunkSize = 1048576
	}
	uploadedbytes := 0

	for uploadedbytes < len(content) {
		uploadOffset, err := c.headFile(filePath)
		if err != nil {
			return fileCreated, err
		}

		_, err = c.patchFile(filePath, content, uploadOffset, chunkSize)
		if err != nil {
			if err.Error() == "chunk too big" {
				logrus.Errorf("chunksize is to big, value: %d reducing it to: %d\n", chunkSize, chunkSize/2)
				chunkSize /= 2
				if chunkSize <= 0 {
					return false, fmt.Errorf("chunk size must be at least 1 byte, but is: %d bytes", chunkSize)
				}
			} else {
				return false, err
			}
		} else {
			uploadedbytes += chunkSize
		}
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

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath)
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

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath)
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

func (c *Client) patchFile(filePath string, content []byte, uploadOffset, chunkSize int) (bool, error) {
	if c.token == "" {
		_, err := c.Auth()
		if err != nil {
			return false, err
		}
	}

	totalBytes := len(content)
	chunkEnd := uploadOffset + chunkSize
	if chunkEnd > totalBytes {
		chunkEnd = totalBytes
	}

	chunk := content[uploadOffset:chunkEnd]

	endpointURL := fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath)
	req, err := http.NewRequest(
		"PATCH",
		endpointURL,
		bytes.NewReader(chunk),
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

	if resp.StatusCode == http.StatusRequestEntityTooLarge {
		return false, errors.New("chunk too big")
	}

	if resp.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("request error: %s", resp.Status)
	}

	return true, nil
}
