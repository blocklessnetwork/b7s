package overseer

import (
	"github.com/Maelkum/limits/limits"
)

type Limiter interface {
	CreateGroup(name string, opts ...limits.LimitOption) (uintptr, error)
	GetGroupHandle(id string) (uintptr, error)
	AssignProcessToGroup(pid uint64, groupID string) error
	DeleteGroup(id string) error
}
