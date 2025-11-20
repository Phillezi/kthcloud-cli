package filebrowser_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kthcloud/cli/pkg/storage/filebrowser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// newTestClient creates a FileBrowserClient with a mock server handler.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...filebrowser.Option) (*filebrowser.FileBrowserClient, *httptest.Server) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	opts = append([]filebrowser.Option{
		filebrowser.WithLogger(zap.NewNop()),
		filebrowser.WithChunkSize(4), // small chunk size for test
	}, opts...)

	client := filebrowser.NewFileBrowserClient(server.URL, "test-token", opts...)
	return client, server
}

func TestUploadFile_Success(t *testing.T) {
	var postCalled, headCalled, patchCalled, totalUploaded int32

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			atomic.AddInt32(&postCalled, 1)
			w.WriteHeader(http.StatusCreated)

		case http.MethodHead:
			atomic.AddInt32(&headCalled, 1)
			w.Header().Set("Upload-Offset", "0")
			w.WriteHeader(http.StatusOK)

		case http.MethodPatch:
			atomic.AddInt32(&patchCalled, 1)
			n, _ := io.Copy(io.Discard, r.Body)
			atomic.AddInt32(&totalUploaded, int32(n))
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})

	ctx := context.Background()
	progressCh := make(chan filebrowser.UploadProgress, 10)
	go func() {
		for range progressCh {
		}
	}()

	content := bytes.NewBufferString("abcdabcd") // 8 bytes -> 2 chunks of 4
	err := client.UploadFileWithProgress(ctx, "/upload/test.txt", content, progressCh)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), postCalled)
	assert.GreaterOrEqual(t, headCalled, int32(1))
	assert.GreaterOrEqual(t, patchCalled, int32(2)) // allow multi-chunk
	assert.Equal(t, int32(8), totalUploaded)
}

func TestUploadFile_MiddlewareApplied(t *testing.T) {
	middlewareCalled := false
	middleware := func(ctx context.Context, req *http.Request) error {
		middlewareCalled = true
		req.Header.Set("X-Middleware", "ok")
		return nil
	}

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ok", r.Header.Get("X-Middleware"))
		w.Header().Set("Upload-Offset", "0")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			w.WriteHeader(http.StatusNoContent)
		}
	}, filebrowser.WithMiddleware(middleware))

	err := client.UploadFileWithProgress(context.Background(), "/file", bytes.NewBufferString("abc"), nil)
	assert.NoError(t, err)
	assert.True(t, middlewareCalled)
}

func TestUploadFile_RetryAndBackoff(t *testing.T) {
	var patchAttempts int32

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.Header().Set("Upload-Offset", "0")
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			attempt := atomic.AddInt32(&patchAttempts, 1)
			if attempt < 3 {
				http.Error(w, "temporary error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	}, filebrowser.WithRetries(3, 10*time.Millisecond))

	err := client.UploadFileWithProgress(context.Background(), "/file", bytes.NewBufferString("abc"), nil)
	assert.NoError(t, err)
	assert.Equal(t, int32(3), patchAttempts)
}

func TestUploadFile_ContextCancel(t *testing.T) {
	blocking := make(chan struct{})

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.Header().Set("Upload-Offset", "0")
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			<-blocking // block until canceled
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := client.UploadFileWithProgress(ctx, "/blocked", bytes.NewBufferString("abc"), nil)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	close(blocking)
}

func TestUploadFile_ServerError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})

	err := client.UploadFileWithProgress(context.Background(), "/bad", bytes.NewBufferString("abc"), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response")
}

func TestHeadFile_MissingUploadOffset(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.WriteHeader(http.StatusOK)
		}
	})

	err := client.UploadFileWithProgress(context.Background(), "/missing", bytes.NewBufferString("abc"), nil)
	assert.True(t, errors.Is(err, filebrowser.ErrMissingUploadOffsetHeader))
}

func TestPatchFile_ChunkTooBig(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.Header().Set("Upload-Offset", "0")
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			http.Error(w, "too big", http.StatusRequestEntityTooLarge)
		}
	})
	err := client.UploadFileWithProgress(context.Background(), "/toobig", bytes.NewBufferString("abc"), nil)
	assert.ErrorContains(t, err, "chunk too big")
}

func TestMiddleware_ErrorIgnored(t *testing.T) {
	middleware := func(ctx context.Context, req *http.Request) error {
		return errors.New("middleware failure")
	}
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Upload-Offset", "0")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodHead:
			w.WriteHeader(http.StatusOK)
		case http.MethodPatch:
			w.WriteHeader(http.StatusNoContent)
		}
	}, filebrowser.WithMiddleware(middleware))

	err := client.UploadFileWithProgress(context.Background(), "/middleware", bytes.NewBufferString("abc"), nil)
	assert.NoError(t, err)
}
