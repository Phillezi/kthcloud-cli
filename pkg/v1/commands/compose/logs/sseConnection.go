package logs

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kthcloud/go-deploy/dto/v2/body"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/sirupsen/logrus"
)

type SSEConnection struct {
	Name  string
	URL   string
	Color string
	Token string
	Key   string
}

func (sse *SSEConnection) OpenConnection(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s%s\033[0m: Context cancelled, closing connection...", sse.Color, sse.Name)
			return
		default:
			req, err := http.NewRequest("GET", sse.URL, nil)
			if err != nil {
				log.Printf("%s%s\033[0m: Failed to create request for %s: %v", sse.Color, sse.Name, sse.URL, err)
				return
			}

			if sse.Key != "" {
				req.Header.Set("X-Api", sse.Key)
			} else if sse.Token != "" {
				req.Header.Set("Authorization", "Bearer "+sse.Token)
			}

			req.Header.Set("Accept", "text/event-stream")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("%s%s\033[0m: Failed to connect to %s: %v", sse.Color, sse.Name, sse.URL, err)
				return
			}
			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)

			var event Event
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						log.Printf("%s%s\033[0m: Connection closed by server, reconnecting...", sse.Color, sse.Name)
					} else {
						log.Printf("%s%s\033[0m: Error reading stream: %v", sse.Color, sse.Name, err)
					}
					break
				}

				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, "data:") {
					event.Data = strings.TrimSpace(line[5:])
					log, err := util.ProcessResponse[body.LogMessage](event.Data)
					if err != nil {
						logrus.Errorln(err)
					}
					if !strings.Contains(event.Data, `"source":"keep-alive"`) {
						fmt.Printf("%s%s\033[0m: %s\n", sse.Color, sse.Name, log.Line)
					}
				}

				if line == "" {
					event = Event{}
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func CreateConns(depls []*body.DeploymentRead, apiURL, token, key string) []*SSEConnection {
	conns := make([]*SSEConnection, 0)
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
		conns = append(conns, &SSEConnection{
			Name:  depl.Name,
			URL:   apiURL + "/v2/deployments/" + depl.ID + "/logs-sse",
			Token: token,
			Key:   key,
			Color: colors[i%len(colors)],
		})
	}
	return conns
}
