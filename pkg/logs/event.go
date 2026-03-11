package logs

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
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

// func (l *LogLine) String() string {
// 	return fmt.Sprintf("%-26s%s%-20s%s", l.CreatedAt.Local().Format(time.RFC3339), l.Color, l.Name+"\033[0m:", l.Line)
// }

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visibleLen counts only printable characters (no ANSI sequences)
func visibleLen(s string) int {
	clean := ansiRegexp.ReplaceAllString(s, "")
	return utf8.RuneCountInString(clean)
}

// wrapText wraps text with the specified indentation.
func wrapText(text string, width int, indent string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var b strings.Builder
	lineLen := 0

	for i, w := range words {
		wLen := utf8.RuneCountInString(w)

		// Start new line
		if i == 0 {
			b.WriteString(w)
			lineLen = wLen
			continue
		}

		if lineLen+1+wLen > width {
			b.WriteString("\n")
			b.WriteString(indent)
			b.WriteString(w)
			lineLen = len(indent) + wLen
		} else {
			b.WriteString(" ")
			b.WriteString(w)
			lineLen += 1 + wLen
		}
	}

	return b.String()
}

func getTerminalWidth() int {
	fd := int(os.Stdout.Fd())

	if term.IsTerminal(fd) {
		w, _, err := term.GetSize(fd)
		if err == nil && w > 0 {
			return w
		}
	}

	// Safe fallback
	return 100
}

func (l *LogLine) String() string {
	// Timestamp + color + name + ":" before log text
	prefix := fmt.Sprintf(
		"%-26s%s%-20s ",
		l.CreatedAt.Local().Format(time.RFC3339),
		l.Color,
		l.Name+"\033[0m:",
	)

	// Visible indent alignment for wrapped lines
	indent := strings.Repeat(" ", visibleLen(prefix))

	termWidth := getTerminalWidth()

	wrapped := wrapText(l.Line, termWidth-visibleLen(prefix), indent)

	return prefix + wrapped
}
