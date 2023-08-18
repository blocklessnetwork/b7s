package node

import (
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/request"
)

// convert the incoming message format to an execution request.
func createExecuteRequest(req request.Execute) execute.Request {

	er := execute.Request{
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
	}

	return er
}
