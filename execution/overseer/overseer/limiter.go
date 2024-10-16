package overseer

import (
	"github.com/blocklessnetwork/b7s/execution/limits"
)

type Limiter interface {
	// TODO: Think about this, it's an interface but directly linking to a specific
	// package for limit options..?
	CreateGroup(id string, opts ...limits.LimitOption) error

	// TODO: Too directly tied to the limiter implementation too.
	GetHandle(id string) (uintptr, error)

	DeleteGroup(id string) error
}
