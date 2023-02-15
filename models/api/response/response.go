package response

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the REST API response for function execution.
type Execute execute.Result

// InstallFunction describes the REST API response for the function install.
type InstallFunction struct {
	Code string `json:"code"`
}
