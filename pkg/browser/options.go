package browser

import "io"

type Option func(*options)

type options struct {
	dryRun bool
	out    io.Writer
}

func WithDryRun(w io.Writer) Option {
	return func(o *options) {
		o.dryRun = true
		o.out = w
	}
}

func WithOutput(w io.Writer) Option {
	return func(o *options) {
		o.out = w
	}
}
