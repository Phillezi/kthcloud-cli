package resources

import (
	"go-deploy/dto/v2/body"
	"time"
)

type CachedResource[T any] struct {
	Data      T
	CachedAt  time.Time
	ExpiresIn time.Duration
}

type Resources struct {
	User        *CachedResource[*body.UserRead]
	Deployments *CachedResource[[]body.DeploymentRead]
	Vms         *CachedResource[[]body.VmRead]
}

func (r *CachedResource[T]) IsExpired() bool {
	if r == nil {
		return true
	}
	return time.Since(r.CachedAt) > r.ExpiresIn
}

func (r *Resources) DropUserCache() {
	if r.User != nil {
		r.User = nil
	}
}

func (r *Resources) DropDeploymentsCache() {
	if r.Deployments != nil {
		r.Deployments = nil
	}
}

func (r *Resources) DropVmsCache() {
	if r.Vms != nil {
		r.Vms = nil
	}
}
