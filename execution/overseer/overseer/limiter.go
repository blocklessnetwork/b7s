package overseer

import (
	"github.com/blocklessnetwork/b7s/execution/limits"
)

type Limiter interface {
	CreateGroup(id string, opts ...limits.LimitOption) error
	GetGroupHandle(id string) (uintptr, error)
	AssignProcessToGroup(pid uint64, groupID string) error
	DeleteGroup(id string) error
}
