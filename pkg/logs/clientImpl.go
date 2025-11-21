package logs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/kthcloud/cli/internal/defaults"
	"github.com/kthcloud/cli/pkg/deploy"
	"github.com/kthcloud/cli/pkg/session"
	"go.uber.org/zap"
)

type ClientImpl struct {
	ctx context.Context

	in chan LogLine

	initialBatchDelay time.Duration
	apiURL            string

	session session.Auth

	logger *zap.Logger
}

func New(opts ...Option) *ClientImpl {
	c := &ClientImpl{
		in:                make(chan LogLine, 128),
		initialBatchDelay: 500 * time.Millisecond,
		apiURL:            defaults.DefaultDeployAPIBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if c.logger == nil {
		c.logger = zap.NewNop()
	}

	if c.session == nil {
		c.session = session.APITokenSession("")
		c.logger.Warn("No auth provided for logs")
	}

	return c
}

func (c *ClientImpl) Consume(w io.Writer) error {
	// Wait for first message
	firstMsg, ok := <-c.in
	if !ok {
		return io.EOF
	}

	buffer := []LogLine{firstMsg}
	done := make(chan struct{})

	// End batching after initialBatchDelay
	go func() {
		select {
		case <-time.After(c.initialBatchDelay):
		case <-c.ctx.Done():
		}
		close(done)
	}()

	batching := true

	for batching {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case msg, ok := <-c.in:
			if !ok {
				// channel closed during batching, flush sorted buffer
				sort.Slice(buffer, func(i, j int) bool {
					return buffer[i].CreatedAt.Before(buffer[j].CreatedAt)
				})
				for _, m := range buffer {
					fmt.Fprintln(w, m.String())
				}
				return nil
			}
			buffer = append(buffer, msg)

		case <-done:
			// End batching phase
			batching = false
		}
	}

	go func() {
		<-c.ctx.Done()
		close(c.in)
	}()

	// End of batching phase, flush the sorted buffer
	sort.Slice(buffer, func(i, j int) bool {
		return buffer[i].CreatedAt.Before(buffer[j].CreatedAt)
	})
	for _, m := range buffer {
		fmt.Fprintln(w, m.String())
	}

	// Done flushing batch, stream the msgs
	for msg := range c.in {
		fmt.Fprintln(w, msg.String())
	}

	return nil
}

func (c *ClientImpl) Subscribe(depls ...*deploy.BodyDeploymentRead) error {
	for _, depl := range depls {
		if depl == nil || depl.Name == nil || depl.Id == nil {
			c.logger.Warn("Skipping logs for nil deployment", zap.Any("deployment", depl))
			continue
		}
		go c.openConnection(&SSEConnection{
			Name:  *depl.Name,
			URL:   c.apiURL + "/v2/deployments/" + *depl.Id + "/logs-sse",
			Color: UUIDColor(*depl.Id),
		})
	}

	return nil
}

func (c *ClientImpl) openConnection(sse *SSEConnection) error {
	var backoff time.Duration = 1 * time.Second
	var retry int
	const maxRetries int = 5
	for {
		select {
		case <-c.ctx.Done():
			c.logger.Sugar().Infof("%s%s\033[0m: Context cancelled, closing connection...", sse.Color, sse.Name)
			return nil
		default:
			req, err := http.NewRequest("GET", sse.URL, nil)
			if err != nil {
				c.logger.Sugar().Errorf("%s%s\033[0m: Failed to create request for %s: %v", sse.Color, sse.Name, sse.URL, err)
				return err
			}

			if err := c.session.AuthMiddleware(c.ctx, req); err != nil {
				c.logger.Sugar().Errorf("%s%s\033[0m: Failed to add auth middleware for %s: %v", sse.Color, sse.Name, sse.URL, err)
				return err
			}

			req.Header.Set("Accept", "text/event-stream")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				c.logger.Sugar().Errorf("%s%s\033[0m: Failed to connect to %s: %v", sse.Color, sse.Name, sse.URL, err)
				return err
			}
			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)

			var event Event
			for c.ctx.Err() == nil {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						if retry > maxRetries {
							return err
						}
						c.logger.Sugar().Warnf("%s%s\033[0m: Connection closed by server [%d/%d], reconnecting in %3v...", sse.Color, sse.Name, retry, maxRetries, backoff)
						retry++
						select {
						case <-c.ctx.Done():
						case <-time.After(backoff):
							backoff = min(backoff*2, 10*time.Second)
						}
					} else {
						c.logger.Sugar().Errorf("%s%s\033[0m: Error reading stream: %v", sse.Color, sse.Name, err)
						return err
					}
					break
				}

				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, "data:") {
					event.Data = strings.TrimSpace(line[5:])

					var log LogMessage
					if err := json.Unmarshal([]byte(event.Data), &log); err != nil {
						return fmt.Errorf("failed to unmarshal response: %w", err)
					}

					if !strings.Contains(event.Data, `"source":"keep-alive"`) {
						func() {
							defer func() {
								if r := recover(); r != nil {
								}
							}()
							c.in <- LogLine{
								LogMessage: log,
								Color:      sse.Color,
								Name:       sse.Name,
							}
						}()
					}
				}

				if line == "" {
					event = Event{}
				}
			}
		}
	}
}
