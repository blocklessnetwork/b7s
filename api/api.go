package api

import (
	"github.com/rs/zerolog"
)

// API provides REST API functionality for the Blockless head node.
type API struct {
	log  zerolog.Logger
	node Node
}

// New creates a new instance of a Blockless head node REST API. Access to node data is provided by the provided `node`.
func New(log zerolog.Logger, node Node) *API {

	api := API{
		log:  log,
		node: node,
	}

	return &api
}
