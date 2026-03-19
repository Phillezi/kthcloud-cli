package out

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type Format = string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

func Write(v any, format Format, w io.Writer) error {
	switch format {
	case FormatJSON:
		dat, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(dat)); err != nil {
			return err
		}
	case FormatYAML:
		dat, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(dat)); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintln(w, v); err != nil {
			return err
		}
	}
	return nil
}
