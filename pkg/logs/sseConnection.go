package logs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kthcloud/cli/pkg/deploy"
	"go.uber.org/zap"
)

type SSEConnection struct {
	Name  string
	URL   string
	Color string
}

func (sse *SSEConnection) OpenConnection(ctx context.Context, mw ...deploy.RequestEditorFn) error {
	var backoff time.Duration = 1 * time.Second
	var retry int
	const maxRetries int = 5
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s%s\033[0m: Context cancelled, closing connection...", sse.Color, sse.Name)
			return nil
		default:
			req, err := http.NewRequest("GET", sse.URL, nil)
			if err != nil {
				log.Printf("%s%s\033[0m: Failed to create request for %s: %v", sse.Color, sse.Name, sse.URL, err)
				return err
			}

			for _, m := range mw {
				if err := m(ctx, req); err != nil {
					return err
				}
			}

			req.Header.Set("Accept", "text/event-stream")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("%s%s\033[0m: Failed to connect to %s: %v", sse.Color, sse.Name, sse.URL, err)
				return err
			}
			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)

			var event Event
			for ctx.Err() == nil {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						if retry > maxRetries {
							return err
						}
						log.Printf("%s%s\033[0m: Connection closed by server [%d/%d], reconnecting in %3v...", sse.Color, sse.Name, retry, maxRetries, backoff)
						retry++
						select {
						case <-ctx.Done():
						case <-time.After(backoff):
							backoff = min(backoff*2, 10*time.Second)
						}
					} else {
						log.Printf("%s%s\033[0m: Error reading stream: %v", sse.Color, sse.Name, err)
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
						fmt.Printf("%s%-20s%s\n", sse.Color, sse.Name+"\033[0m:", log.Line)
					}
				}

				if line == "" {
					event = Event{}
				}
			}
		}
	}
}

func CreateConns(depls []*deploy.BodyDeploymentRead, apiURL string) []*SSEConnection {
	conns := make([]*SSEConnection, 0, len(depls))
	colors := []string{
		"\033[31m", // Red
		"\033[32m", // Green
		"\033[33m", // Yellow
		"\033[34m", // Blue
		"\033[35m", // Magenta
		"\033[36m", // Cyan
		"\033[37m", // White
	}
	for i, depl := range depls {
		if depl == nil || depl.Name == nil || depl.Id == nil {
			zap.L().Warn("Skipping logs for nil deployment", zap.Any("deployment", depl))
			continue
		}
		conns = append(conns, &SSEConnection{
			Name:  *depl.Name,
			URL:   apiURL + "/v2/deployments/" + *depl.Id + "/logs-sse",
			Color: colors[i%len(colors)],
		})
	}
	return conns
}

func CreateConn(name, id, apiURL string) *SSEConnection {
	return &SSEConnection{
		Name:  name,
		URL:   apiURL + "/v2/deployments/" + id + "/logs-sse",
		Color: "\033[34m",
	}
}
