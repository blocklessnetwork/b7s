package node

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// getFunctionManifest retrieves the function manifest for the function with the given ID.
func (n *Node) getFunctionManifest(id string) (*blockless.FunctionManifest, error) {

	// Try to get function manifest from the store.
	var manifest blockless.FunctionManifest
	err := n.store.GetRecord(id, &manifest)
	if err != nil {
		// TODO: Check - error not found.
		return nil, fmt.Errorf("could not retrieve function manifest: %w", err)
	}

	return &manifest, nil
}
