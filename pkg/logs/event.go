package logs

import (
	"fmt"
	"time"
)

type Event struct {
	ID    string
	Event string
	Data  string
}

type LogMessage struct {
	Source    string    `json:"source"`
	Prefix    string    `json:"prefix"`
	Line      string    `json:"line"`
	CreatedAt time.Time `json:"createdAt"`
}

type LogLine struct {
	LogMessage `json:",inline"`
	Color      string `json:"color"`
	Name       string `json:"name"`
}

func (l *LogLine) String() string {
	return fmt.Sprintf("%-26s%s%-20s%s", l.CreatedAt.Local().Format(time.RFC3339), l.Color, l.Name+"\033[0m:", l.Line)
}
