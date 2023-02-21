package request

import (
	"github.com/blocklessnetworking/b7s/models/execute"
)

// Execute describes the payload for the REST API request for function execution.
type Execute execute.Request

// InstallFunction describes the payload for the REST API request for function install.
type InstallFunction struct {
	CID string `query:"cid"`
	URI string `query:"uri"`
}
