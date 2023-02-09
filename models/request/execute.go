package request

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	Type       string             `json:"type,omitempty"`
	From       peer.ID            `json:"from,omitempty"`
	Code       string             `json:"code,omitempty"`
	FunctionID string             `json:"functionId,omitempty"`
	Method     string             `json:"method,omitempty"`
	Parameters []ExecuteParameter `json:"parameters,omitempty"`
	Config     ExecutionConfig    `json:"config,omitempty"`
}

// ExecuteParameter represents an execution parameter, modeled by a key-value pair.
// TODO: All key-value pairs can perhaps be modeled the same?
type ExecuteParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ExecutionConfig represents the configurable options for an execution request.
type ExecutionConfig struct {
	Environment       []ExecuteEnvVars         `json:"env_vars"`
	NodeCount         int                      `json:"number_of_nodes"`
	ResultAggregation ExecuteResultAggregation `json:"result_aggregation"`
	Stdin             *string                  `json:"stdin"`
	Permissions       []string                 `json:"permissions"`
}

// ExecuteEnvVars represents the name and value of the environment variables set for the execution.
type ExecuteEnvVars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExecuteResultAggregation struct {
	Enable     bool               `json:"enable"`
	Type       string             `json:"type"`
	Parameters []ExecuteParameter `json:"parameters"`
}
