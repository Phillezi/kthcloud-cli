package renderer

import (
	"io"
	"strings"
)

type Output int

const (
	Output_Table Output = iota
	Output_JSON
	Output_YAML
)

func OutputFromString(s string) (o Output) {
	switch strings.ToLower(s) {
	case "json":
		o = Output_JSON
	case "yaml":
		o = Output_YAML
	}
	return
}

type RenderConfig struct {
	Output Output
	W      io.Writer
}

type RenderOptions func(rc *RenderConfig)

type Renderer interface {
	Render(obj any, options ...RenderOptions) error
}

// Option helpers
func WithOutput(output Output) RenderOptions {
	return func(rc *RenderConfig) { rc.Output = output }
}

func WithWriter(w io.Writer) RenderOptions {
	return func(rc *RenderConfig) { rc.W = w }
}

// New creates a new generic renderer
func New() Renderer {
	return &DefaultRendererImpl{}
}
