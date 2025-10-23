package renderer

import (
	"encoding/json"
	"io"
)

func renderJSON(w io.Writer, obj any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}
