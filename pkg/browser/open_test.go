package browser_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kthcloud/cli/pkg/browser"
)

func TestOpenDryRun(t *testing.T) {
	var buf bytes.Buffer
	url := "https://example.com"

	err := browser.Open(url, browser.WithDryRun(&buf))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, url) {
		t.Errorf("expected output to contain URL, got: %q", output)
	}

	if !strings.Contains(output, "[browser] would execute:") {
		t.Errorf("expected output to show dry run, got: %q", output)
	}
}
