package util

import (
	"fmt"
	"time"
)

func TimeAgo(t *time.Time) string {
	if t == nil {
		return ""
	}
	duration := time.Since(*t)

	switch {
	case duration.Minutes() < 1:
		return "just now"
	case duration.Hours() < 1:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}
