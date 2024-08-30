package fstore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/armon/go-metrics"
	"github.com/cavaliergopher/grab/v3"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func (h *FStore) getJSON(address string, out interface{}) error {

	h.log.Debug().Str("url", address).Msg("retrieving JSON doc")

	res, err := h.http.Get(address)
	if err != nil {
		return fmt.Errorf("could not get resource (url: %s): %w", address, err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return fmt.Errorf("could not unpack data (url: %s): %w", address, err)
	}

	return nil
}

// download will retrieve the function with the given manifest. It returns the full path
// of the file where the function is saved on the local storage or any error that might have
// occurred in the process. The function blocks until the download is complete.
func (h *FStore) download(ctx context.Context, cid string, manifest blockless.FunctionManifest) (string, error) {

	// Determine directory where files should be stored.
	fdir := filepath.Join(h.workdir, cid)

	h.log.Info().
		Str("target_dir", fdir).
		Str("cid", cid).
		Str("function_uri", manifest.Deployment.URI).
		Msg("downloading function")

	// Create destination directory.
	err := os.MkdirAll(fdir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("could not create destination directory (dir: %s): %w", fdir, err)
	}

	// Get expected checksum of the function.
	sum, err := hex.DecodeString(manifest.Deployment.Checksum)
	if err != nil {
		return "", fmt.Errorf("invalid function checksum (sum: %s): %w", manifest.Deployment.Checksum, err)
	}

	// Create a new download request.
	req, err := grab.NewRequest(fdir, manifest.Deployment.URI)
	if err != nil {
		return "", fmt.Errorf("error creating download request: %w", err)
	}
	req.SetChecksum(sha256.New(), sum, true)
	req.NoCreateDirectories = false
	req = req.WithContext(ctx)

	// Execute the download request.
	res := h.downloader.Do(req)

	metrics.IncrCounter(functionsDownloadedSizeMetric, float32(res.HTTPResponse.ContentLength))

	// Wait until the download is complete.
	err = res.Err()
	if err != nil {
		return "", fmt.Errorf("could not download function: %w", err)
	}

	h.log.Info().
		Str("output", res.Filename).
		Str("cid", cid).
		Str("function_uri", manifest.Deployment.URI).
		Msg("downloaded function")

	return res.Filename, nil
}
