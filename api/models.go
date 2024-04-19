package api

import (
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/node/aggregate"
)

// ExecuteRequest describes the payload for the REST API request for function execution.
type ExecuteRequest struct {
	execute.Request
	Topic string `json:"topic,omitempty"`
}

// ExecuteResponse describes the REST API response for function execution.
type ExecuteResponse struct {
	Code      codes.Code        `json:"code,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Message   string            `json:"message,omitempty"`
	Results   aggregate.Results `json:"results,omitempty"`
	Cluster   execute.Cluster   `json:"cluster,omitempty"`
}

// InstallFunctionRequest describes the payload for the REST API request for function install.
type InstallFunctionRequest struct {
	CID      string `json:"cid"`
	URI      string `json:"uri"`
	Subgroup string `json:"subgroup"`
}

// InstallFunctionResponse describes the REST API response for the function install.
type InstallFunctionResponse struct {
	Code string `json:"code"`
}

// ExecutionResultRequest describes the payload for the REST API request for execution result.
type ExecutionResultRequest struct {
	ID string `json:"id"`
}
