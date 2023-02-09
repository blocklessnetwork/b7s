package node

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

type Execute interface {
	Function(execute.Request) (execute.Response, error)
}
