package api

// API provides REST API functionality for the Blockless head node.
type API struct {
	node Node
}

// New creates a new instance of a Blockless head node REST API. Access to node data is provided by the provided `node`.
func New(node Node) *API {

	api := API{
		node: node,
	}

	return &api
}
