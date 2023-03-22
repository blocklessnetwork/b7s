package node

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

type Executor interface {
	Function(string, execute.Request) (execute.Result, error)
}
