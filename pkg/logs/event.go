package logs

import "time"

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
