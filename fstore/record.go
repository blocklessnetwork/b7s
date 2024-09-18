package fstore

import (
	"context"
	"fmt"
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// Get retrieves a function manifest for the given function from storage.
func (f *FStore) Get(ctx context.Context, cid string) (blockless.FunctionRecord, error) {

	fn, err := f.getFunction(ctx, cid)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not get function from store: %w", err)
	}

	return fn, nil
}

func (f *FStore) getFunction(ctx context.Context, cid string) (blockless.FunctionRecord, error) {

	function, err := f.store.RetrieveFunction(ctx, cid)
	if err != nil {
		return blockless.FunctionRecord{}, fmt.Errorf("could not retrieve function record: %w", err)
	}

	go func() {
		// Update the "last retrieved" timestamp.
		function.LastRetrieved = time.Now().UTC()
		err = f.store.SaveFunction(context.Background(), function)
		if err != nil {
			f.log.Warn().Err(err).Str("cid", cid).Msg("could not update function record timestamp")
		}
	}()

	return function, nil
}

func (f *FStore) saveFunction(ctx context.Context, fn blockless.FunctionRecord) error {

	// Clean paths - make them relative to the current working directory.
	fn.Archive = f.cleanPath(fn.Archive)
	fn.Files = f.cleanPath(fn.Files)
	fn.Manifest.Deployment.File = f.cleanPath(fn.Manifest.Deployment.File)

	fn.UpdatedAt = time.Now().UTC()
	return f.store.SaveFunction(ctx, fn)
}
