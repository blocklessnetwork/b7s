package fstore

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Get retrieves a function manifest for the given function from storage.
func (h *FStore) Get(cid string) (*blockless.FunctionManifest, error) {

	fn, err := h.getFunction(cid)
	if err != nil {
		return nil, fmt.Errorf("could not get function from store: %w", err)
	}

	return &fn.Manifest, nil
}
