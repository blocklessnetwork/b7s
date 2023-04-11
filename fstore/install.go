package fstore

import (
	"errors"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Installed checks if the function with the given CID is installed.
func (h *FStore) Installed(cid string) (bool, error) {

	fn, err := h.getFunction(cid)
	if err != nil && errors.Is(err, blockless.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not get function from store: %w", err)
	}

	haveArchive, haveFiles, err := h.checkFunctionFiles(*fn)
	if err != nil {
		return false, fmt.Errorf("could not verify function cache: %w", err)
	}

	// If we don't have all files found, treat it as not installed.
	if !haveArchive || !haveFiles {
		return false, nil
	}

	// We have the function in the database and all files - we're good.
	return true, nil
}
