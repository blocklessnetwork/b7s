package fstore

import (
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func (h *FStore) getFunction(cid string) (*blockless.FunctionRecord, error) {

	fn, err := h.store.RetrieveFunction(cid)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve function record: %w", err)
	}

	// Update the "last retrieved" timestamp.
	fn.LastRetrieved = time.Now().UTC()
	err = h.store.SaveFunction(cid, fn)
	if err != nil {
		h.log.Warn().Err(err).Str("cid", cid).Msg("could not update function record timestamp")
	}

	return &fn, nil
}

func (h *FStore) saveFunction(fn blockless.FunctionRecord) error {

	// Clean paths - make them relative to the current working directory.
	fn.Archive = h.cleanPath(fn.Archive)
	fn.Files = h.cleanPath(fn.Files)

	fn.UpdatedAt = time.Now().UTC()
	return h.store.SaveFunction(fn.CID, fn)
}
