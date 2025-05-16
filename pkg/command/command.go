package command

import "context"

type Command interface {
	WithContext(ctx context.Context) Command
	Run() error
}
