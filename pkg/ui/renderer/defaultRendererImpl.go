package renderer

import "os"

type DefaultRendererImpl struct {
}

func (DefaultRendererImpl) Render(obj any, options ...RenderOptions) error {
	cfg := &RenderConfig{
		Output: Output_Table,
		W:      os.Stdout,
	}
	for _, opt := range options {
		opt(cfg)
	}

	switch cfg.Output {
	case Output_JSON:
		return renderJSON(cfg.W, obj)
	case Output_YAML:
		return renderYAML(cfg.W, obj)
	default:
		return renderTable(cfg.W, obj)
	}
}
