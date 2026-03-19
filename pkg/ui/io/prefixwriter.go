package io

import (
	"bytes"
	"io"
)

type PrefixWriter struct {
	Prefix string
	Writer io.Writer
	buf    bytes.Buffer
}

func (p *PrefixWriter) Write(data []byte) (int, error) {
	for _, b := range data {
		if b == '\n' {
			if _, err := p.Writer.Write([]byte(p.Prefix + p.buf.String() + "\n")); err != nil {
				return 0, err
			}
			p.buf.Reset()
		} else {
			p.buf.WriteByte(b)
		}
	}
	return len(data), nil
}
