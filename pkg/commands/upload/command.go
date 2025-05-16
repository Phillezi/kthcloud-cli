package upload

import (
	"context"

	"github.com/Phillezi/kthcloud-cli/pkg/defaults"
	"github.com/Phillezi/kthcloud-cli/pkg/deploy"
	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Command struct {
	ctx    context.Context
	client *deploy.Client

	srcPath  string
	destPath string

	storageURL      string
	keycloakBaseURL string
}

func New(opts ...CommandOpts) *Command {
	var opt CommandOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	c := &Command{
		ctx:    util.PtrOr(opt.Context, context.Background()),
		client: opt.Client,

		srcPath:  util.PtrOr(opt.SrcPath),
		destPath: util.PtrOr(opt.DestPath),

		keycloakBaseURL: util.PtrOr(opt.KeycloakBaseURL, defaults.DefaultKeycloakBaseURL),
	}

	c.storageURL = util.PtrOr(opt.StorageURL, func() string {
		if c.client == nil {
			return ""
		}
		if user, err := c.client.User(); err != nil && user != nil {
			return *user.StorageURL
		}
		return ""
	}())

	return c
}

func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}
