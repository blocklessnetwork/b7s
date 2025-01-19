package api

import (
	"github.com/rs/zerolog"
)

// API provides REST API functionality for the Bless head node.
type API struct {
	Log  zerolog.Logger
	Node Node
}

// New creates a new instance of a Bless head node REST API. Access to node data is provided by the provided `node`.
func New(log zerolog.Logger, node Node) *API {

	api := API{
		Log:  log,
		Node: node,
	}

	return &api
}
