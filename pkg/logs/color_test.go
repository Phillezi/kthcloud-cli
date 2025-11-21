package logs_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/kthcloud/cli/pkg/logs"
)

func TestUUIDToANSI(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"UUID1", "550e8400-e29b-41d4-a716-446655440000"},
		{"UUID2", "123e4567-e89b-12d3-a456-426614174000"},
		{"UUID3", "ffffffff-ffff-ffff-ffff-ffffffffffff"},
	}

	ansiPattern := regexp.MustCompile(`^\x1b\[38;2;(\d{1,3});(\d{1,3});(\d{1,3});m$`)

	colors := make(map[string]string)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := logs.UUIDColor(tt.uuid)

			// Verify ANSI format
			matches := ansiPattern.FindStringSubmatch(color)
			if matches == nil {
				t.Fatalf("invalid ANSI format: %q", color)
			}

			// RGB must be in allowed range 80–215
			for i := 1; i <= 3; i++ {
				v := atoi(t, matches[i])
				if v < 80 || v > 215 {
					t.Fatalf("RGB component out of expected range: %d", v)
				}
			}

			// Deterministic: same UUID always returns same ANSI
			color2 := logs.UUIDColor(tt.uuid)
			if color != color2 {
				t.Fatalf("nondeterministic result: %q vs %q", color, color2)
			}

			// Collect for uniqueness check
			colors[tt.uuid] = color
		})
	}

	// Different UUIDs should not map to the same color
	if colors[tests[0].uuid] == colors[tests[1].uuid] ||
		colors[tests[1].uuid] == colors[tests[2].uuid] ||
		colors[tests[0].uuid] == colors[tests[2].uuid] {
		t.Errorf("different UUIDs produced identical colors: %+v", colors)
	}
}

func atoi(t *testing.T, s string) int {
	t.Helper()
	v, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("invalid integer %q: %v", s, err)
	}
	return v
}
