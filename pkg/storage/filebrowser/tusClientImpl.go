package filebrowser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// UploadProgress represents progress information during an upload.
type UploadProgress struct {
	Uploaded int64
	Total    int64
}

// Middleware is a function that can modify an outgoing request before it’s sent.
type Middleware func(ctx context.Context, req *http.Request) error

// FileBrowserClient implements resumable uploads with middleware and logging.
type FileBrowserClient struct {
	client         *http.Client
	filebrowserURL string
	token          string
	middleware     []Middleware
	logger         *zap.Logger
	maxParallelism int
	maxRetries     int
	baseBackoff    time.Duration
	chunkSize      int64
}

// Option is a functional option for configuring the client.
type Option func(*FileBrowserClient)

func WithLogger(l *zap.Logger) Option {
	return func(c *FileBrowserClient) { c.logger = l }
}
func WithMiddleware(m ...Middleware) Option {
	return func(c *FileBrowserClient) { c.middleware = append(c.middleware, m...) }
}
func WithChunkSize(size int64) Option {
	return func(c *FileBrowserClient) { c.chunkSize = size }
}
func WithRetries(max int, baseBackoff time.Duration) Option {
	return func(c *FileBrowserClient) {
		c.maxRetries = max
		c.baseBackoff = baseBackoff
	}
}

// NewFileBrowserClient returns a configured FileBrowserClient.
func NewFileBrowserClient(url, token string, opts ...Option) *FileBrowserClient {
	c := &FileBrowserClient{
		client:         &http.Client{Timeout: 0},
		filebrowserURL: url,
		token:          token,
		maxParallelism: 3,
		maxRetries:     5,
		baseBackoff:    500 * time.Millisecond,
		chunkSize:      1 << 20, // 1 MB
		logger:         zap.NewNop(),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WriteFile uploads a file via TUS with resumable chunks.
func (c *FileBrowserClient) WriteFile(path string, r io.Reader) error {
	ctx := context.Background()
	return c.UploadFileWithProgress(ctx, path, r, nil)
}

func (c *FileBrowserClient) UploadFileWithProgress(
	ctx context.Context,
	path string,
	r io.Reader,
	progressCh chan<- UploadProgress,
) error {
	defer func() {
		if progressCh != nil {
			close(progressCh)
		}
	}()

	// Step 1: Check server offset
	serverOffset, err := c.headFile(ctx, path)
	if err != nil {
		return fmt.Errorf("head failed: %w", err)
	}

	var offset int64

	if serverOffset == 0 {
		// File does not exist on server -> create new
		if _, err := c.postFile(ctx, path); err != nil {
			return fmt.Errorf("failed to create new upload: %w", err)
		}
		offset = 0
	} else {
		offset = int64(serverOffset)
		// Compare local file with server offset to detect overwrite
		localStart := make([]byte, offset)
		n, err := io.ReadFull(r, localStart)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("failed to read initial bytes: %w", err)
		}

		if int64(n) != offset {
			// Local file shorter or modified -> overwrite from 0
			c.logger.Info("local file shorter than server, overwriting from offset 0")
			offset = 0
			if seeker, ok := r.(io.Seeker); ok {
				seeker.Seek(0, io.SeekStart)
			} else {
				r = io.MultiReader(bytes.NewReader(localStart[:n]), r)
			}
		} else {
			c.logger.Info("resuming upload", zap.Int64("offset", offset))
		}
	}

	// Step 2: Upload remaining chunks
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		chunk := make([]byte, c.chunkSize)
		n, err := io.ReadFull(r, chunk)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n == 0 {
				break
			}
			chunk = chunk[:n]
		} else if err != nil {
			return fmt.Errorf("read failed: %w", err)
		}

		if err := c.uploadChunkWithRetry(ctx, path, chunk, &offset, progressCh); err != nil {
			return err
		}
	}

	return nil
}

// uploadChunkWithRetry uploads one chunk with exponential backoff retries.
func (c *FileBrowserClient) uploadChunkWithRetry(
	ctx context.Context,
	path string,
	chunk []byte,
	offset *int64,
	progressCh chan<- UploadProgress,
) error {
	var attempt int
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		ok, err := c.patchFile(ctx, path, chunk, int(*offset))
		if err == nil && ok {
			*offset += int64(len(chunk))
			if progressCh != nil {
				progressCh <- UploadProgress{Uploaded: *offset}
			}
			return nil
		}

		if errors.Is(err, ErrChunkTooBig) {
			return fmt.Errorf("chunk too big: %w", err)
		}

		attempt++
		if attempt > c.maxRetries {
			return fmt.Errorf("upload failed after %d retries: %w", c.maxRetries, err)
		}

		sleep := c.baseBackoff * time.Duration(math.Pow(2, float64(attempt)))
		sleep += time.Duration(rand.Int63n(int64(sleep / 2))) // jitter

		c.logger.Warn("chunk upload failed, retrying",
			zap.Int("attempt", attempt),
			zap.Duration("sleep", sleep),
			zap.Error(err))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}
}

func (c *FileBrowserClient) postFile(ctx context.Context, filePath string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath),
		nil,
	)
	if err != nil {
		return false, err
	}
	c.applyMiddleware(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("unexpected response: %s", resp.Status)
	}
	return true, nil
}

func (c *FileBrowserClient) headFile(ctx context.Context, filePath string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD",
		fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath),
		nil,
	)
	if err != nil {
		return 0, err
	}
	c.applyMiddleware(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// no existing upload — this is not an error
		return 0, nil
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	offsetStr := resp.Header.Get("Upload-Offset")
	if offsetStr == "" {
		return 0, ErrMissingUploadOffsetHeader
	}

	return strconv.Atoi(offsetStr)
}

func (c *FileBrowserClient) patchFile(ctx context.Context, filePath string, chunk []byte, uploadOffset int) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("%s/api/tus/%s?override=false", c.filebrowserURL, filePath),
		bytes.NewReader(chunk),
	)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", fmt.Sprintf("%d", uploadOffset))

	c.applyMiddleware(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return true, nil
	case http.StatusRequestEntityTooLarge:
		return false, ErrChunkTooBig
	default:
		return false, fmt.Errorf("unexpected status: %s", resp.Status)
	}
}

func (c *FileBrowserClient) Auth() (bool, error) {
	if c.token == "" {
		requestBody := strings.NewReader(`{"username": "", "password": "", "recaptcha": ""}`)

		req, err := http.NewRequest("POST", c.filebrowserURL+"/api/login", requestBody)
		if err != nil {
			return false, err
		}

		c.applyExternalMiddleware(context.TODO(), req)

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
			return false, err
		}
		c.token = string(bodyBytes)
	}
	return true, nil
}

// applyMiddleware applies all middleware (including default token injection).
func (c *FileBrowserClient) applyExternalMiddleware(ctx context.Context, req *http.Request) {
	for _, m := range c.middleware {
		if err := m(ctx, req); err != nil {
			c.logger.Warn("middleware error", zap.Error(err))
		}
	}
}

func (c *FileBrowserClient) applyMiddleware(ctx context.Context, req *http.Request) {
	c.Auth()
	req.Header.Set("X-Auth", c.token)
	c.applyExternalMiddleware(ctx, req)
}

func (c *FileBrowserClient) URL() string {
	return c.filebrowserURL
}
