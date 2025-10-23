package renderer

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func renderYAML(w io.Writer, obj any) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, string(data))
	return err
}
