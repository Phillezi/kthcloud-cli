package body

import (
	"bytes"
	"encoding/json"
	"io"
)

func Reader(v any) (io.Reader, error) {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	return buf, nil
}
