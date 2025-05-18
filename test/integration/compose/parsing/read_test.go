package parsing

import (
	"os"
	"testing"

	"github.com/Phillezi/kthcloud-cli/internal/load"
)

const (
	basePath = "../data/"
)

func TestComposeParsing(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		expectErr    bool
		serviceCount int
	}{
		{"Empty File", basePath + "empty.yaml", true, 0},
		{"Single Service", basePath + "single-service.yaml", false, 1},
		{"Multiple Services", basePath + "multiple-services.yaml", false, 3},
		{"Invalid YAML", basePath + "invalid.yaml", true, 0},
		{"Missing File", basePath + "missing.yaml", true, 0},
		{"No Services", basePath + "no-services.yaml", false, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := os.Stat(tc.filePath); err != nil && !tc.expectErr {
				t.Fatalf("error reading file: %s", tc.filePath)
			}

			comp, err := load.InternalGetCompose(load.LoadOpts{File: tc.filePath})
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error status: %v", err)
			}

			if comp == nil {
				return
			}

			if len(comp.Source.Services) != tc.serviceCount {
				t.Errorf("expected %d services, got %d", tc.serviceCount, len(comp.Source.Services))
			}
		})
	}
}
