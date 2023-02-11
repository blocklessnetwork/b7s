package function

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Get retrieves the function manifest from the given address.
func (h *Handler) Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error) {

	h.log.Debug().
		Str("cid", cid).
		Str("address", address).
		Bool("use_cached", useCached).
		Msg("getting manifest")

	var cachedManifest blockless.FunctionManifest
	err := h.store.GetRecord(cid, &cachedManifest)
	if err != nil {
		// TODO: err not found is not an error.
		return nil, fmt.Errorf("could not get function manifest from store: %w", err)
	}

	// Return cached version if so requested.
	if useCached {
		return &cachedManifest, nil
	}

	// Retrieve function manifest from the given address.
	var manifest blockless.FunctionManifest
	err = h.getJSON(address, &manifest)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve manifest: %w", err)
	}

	// If the runtime URL is specified,
	if cachedManifest.Runtime.URL != "" {
		err = updateDeploymentInfo(&cachedManifest, address)
		if err != nil {
			return nil, fmt.Errorf("could not update deployment info: %w", err)
		}
	}

	// Download the function identified by the manifest.
	path, err := h.download(manifest)
	if err != nil {
		return nil, fmt.Errorf("could not download function: %w", err)
	}

	// Unpack the .tar.gz archive.
	// TODO: Would be good to know the content of the .tar.gz archive.
	// We're unpacking the archive to the file, and storing the path to the .tar.gz in the DB.
	err = h.unpackArchive(path, h.workdir)
	if err != nil {
		return nil, fmt.Errorf("could not unpack gzip archive (file: %s): %w", path, err)
	}

	manifest.Deployment.File = path
	manifest.Cached = true

	// Store the retrieved manifest.
	err = h.store.SetRecord(cid, manifest)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("cid", cid).
			Msg("could not store manifest")
	}

	return &manifest, nil
}
