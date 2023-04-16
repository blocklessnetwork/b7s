package response

import (
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the REST API response for function execution.
type Execute struct {
	Code      codes.Code                `json:"code,omitempty"`
	RequestID string                    `json:"request_id,omitempty"`
	Results   map[string]execute.Result `json:"results,omitempty"`
	// NOTE: Not sending the usage information for now.
	Usage execute.Usage `json:"-"`
}

// InstallFunction describes the REST API response for the function install.
type InstallFunction struct {
	Code string `json:"code"`
}
