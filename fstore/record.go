package fstore

import (
	"context"
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// Get retrieves a function manifest for the given function from storage.
func (h *FStore) Get(ctx context.Context, cid string) (blockless.FunctionRecord, error) {

	fn, err := h.getFunction(ctx, cid)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not get function from store: %w", err)
	}

	return fn, nil
}

func (h *FStore) getFunction(ctx context.Context, cid string) (blockless.FunctionRecord, error) {

	function, err := h.store.RetrieveFunction(ctx, cid)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not retrieve function record: %w", err)
	}

	go func() {
		// Update the "last retrieved" timestamp.
		function.LastRetrieved = time.Now().UTC()
		err = h.store.SaveFunction(context.Background(), function)
		if err != nil {
			h.log.Warn().Err(err).Str("cid", cid).Msg("could not update function record timestamp")
		}
	}()

	return function, nil
}

func (h *FStore) saveFunction(ctx context.Context, fn blockless.FunctionRecord) error {

	// Clean paths - make them relative to the current working directory.
	fn.Archive = h.cleanPath(fn.Archive)
	fn.Files = h.cleanPath(fn.Files)
	fn.Manifest.Deployment.File = h.cleanPath(fn.Manifest.Deployment.File)

	fn.UpdatedAt = time.Now().UTC()
	return h.store.SaveFunction(ctx, fn)
}
