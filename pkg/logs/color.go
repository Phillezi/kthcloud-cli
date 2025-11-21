package logs

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	supportsTruecolorCached bool
	supports256Cached       bool
	initOnce                sync.Once
)

// init runs once per process and caches terminal capabilities.
func init() {
	initOnce.Do(func() {
		supportsTruecolorCached = detectTruecolor()
		supports256Cached = detect256()
	})
}

// UUIDColor returns the best ANSI color code for the UUID
func UUIDColor(uuid string) string {
	sum := sha1.Sum([]byte(uuid))

	r := int(sum[0])
	g := int(sum[1])
	b := int(sum[2])

	// Adjust for readability on dark terminals
	const min = 80
	const max = 215
	r = min + (r % (max - min))
	g = min + (g % (max - min))
	b = min + (b % (max - min))

	// Truecolor
	if supportsTruecolorCached {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	// 256-color fallback
	if supports256Cached {
		idx := 16 + (int(sum[0])%6)*36 + (int(sum[1])%6)*6 + (int(sum[2]) % 6)
		return fmt.Sprintf("\x1b[38;5;%dm", idx)
	}

	// Basic 8-color fallback
	base := 30 + (int(sum[0]) % 8)
	return fmt.Sprintf("\x1b[%dm", base)
}

const ResetANSI = "\x1b[0m"

func detectTruecolor() bool {
	ct := strings.ToLower(os.Getenv("COLORTERM"))
	if strings.Contains(ct, "truecolor") || strings.Contains(ct, "24bit") {
		return true
	}

	term := strings.ToLower(os.Getenv("TERM"))

	// Common truecolor terminals
	truecolorTerms := []string{
		"xterm-truecolor",
		"xterm-direct",
		"tmux-truecolor",
		"screen-truecolor",
		"alacritty",
		"wezterm",
		"iterm",
		"iterm2",
		"kitty",
		"foot",
		"ghostty",
	}

	for _, t := range truecolorTerms {
		if strings.Contains(term, t) {
			return true
		}
	}

	return false
}

func detect256() bool {
	term := strings.ToLower(os.Getenv("TERM"))
	return strings.Contains(term, "256")
}
