package fstore

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Get retrieves the function manifest from the given address. `useCached` indicates whether,
// if the function is found in the store/db, it should be used, or if we should re-download it.
func (h *FStore) Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error) {

	h.log.Debug().
		Str("cid", cid).
		Str("address", address).
		Bool("use_cached", useCached).
		Msg("getting manifest")

	cachedFn, err := h.getFunction(cid)
	// Return cached version if so requested.
	if err == nil && useCached {

		h.log.Debug().
			Str("cid", cid).
			Str("address", address).
			Msg("function manifest was already cached, done")

		return &cachedFn.Manifest, nil
	}
	if err != nil && !errors.Is(err, blockless.ErrNotFound) {
		return nil, fmt.Errorf("could not get function from store: %w", err)
	}

	// Being here means that we either did not find the manifest, or we don't
	// want to use the cached one.

	// Retrieve function manifest from the given address.
	var manifest blockless.FunctionManifest
	err = h.getJSON(address, &manifest)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve manifest: %w", err)
	}

	// If the runtime URL is specified, use it to fill in the deployment info.
	if manifest.Runtime.URL != "" {
		err = updateDeploymentInfo(&manifest, address)
		if err != nil {
			return nil, fmt.Errorf("could not update deployment info: %w", err)
		}
	}

	// Download the function identified by the manifest.
	functionPath, err := h.download(manifest)
	if err != nil {
		return nil, fmt.Errorf("could not download function: %w", err)
	}

	out := filepath.Join(h.workdir, cid)

	// Unpack the .tar.gz archive.
	// TODO: Would be good to know the content of the .tar.gz archive.
	// We're unpacking the archive here and storing the path to the .tar.gz in the DB.
	err = h.unpackArchive(functionPath, out)
	if err != nil {
		return nil, fmt.Errorf("could not unpack gzip archive (file: %s): %w", functionPath, err)
	}

	manifest.Deployment.File = functionPath

	// Store the function record.
	fn := functionRecord{
		CID:      cid,
		URL:      address,
		Manifest: manifest,
		Archive:  functionPath,
		Files:    out,
	}
	err = h.saveFunction(fn)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("cid", cid).
			Msg("could not save function record")
	}

	return &manifest, nil
}