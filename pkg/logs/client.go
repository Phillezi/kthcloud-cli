package logs

import "io"

type Client interface {
	Consume(writer io.Writer) error
	Subscribe() error
}
