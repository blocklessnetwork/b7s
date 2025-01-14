package fstore

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/blessnetwork/b7s/models/blockless"
	"github.com/blessnetwork/b7s/telemetry/b7ssemconv"
)

// Install will download and install function identified by the manifest/CID.
func (f *FStore) Install(ctx context.Context, address string, cid string) (retErr error) {

	defer f.metrics.MeasureSince(functionsInstallTimeMetric, time.Now())
	f.metrics.IncrCounter(functionsInstalledMetric, 1)
	defer func() {
		switch retErr {
		case nil:
			f.metrics.IncrCounter(functionsInstalledOkMetric, 1)
		default:
			f.metrics.IncrCounter(functionsInstalledErrMetric, 1)

		}
	}()

	ctx, span := f.tracer.Start(ctx, spanInstall, trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(b7ssemconv.FunctionCID.String(cid)))
	defer span.End()

	f.log.Debug().
		Str("cid", cid).
		Str("address", address).
		Msg("installing function")

	// Retrieve function manifest from the given address.
	var manifest blockless.FunctionManifest
	err := f.getJSON(address, &manifest)
	if err != nil {
		return fmt.Errorf("could not retrieve manifest: %w", err)
	}

	// If the runtime URL is specified, use it to fill in the deployment info.
	if manifest.Runtime.URL != "" {
		err = updateDeploymentInfo(&manifest, address)
		if err != nil {
			return fmt.Errorf("could not update deployment info: %w", err)
		}
	}

	// Download the function identified by the manifest.
	functionPath, err := f.download(ctx, cid, manifest)
	if err != nil {
		return fmt.Errorf("could not download function: %w", err)
	}

	out := filepath.Join(f.workdir, cid)

	// Unpack the .tar.gz archive.
	// TODO: Would be good to know the content of the .tar.gz archive.
	// We're unpacking the archive here and storing the path to the .tar.gz in the DB.
	err = f.unpackArchive(functionPath, out)
	if err != nil {
		return fmt.Errorf("could not unpack gzip archive (file: %s): %w", functionPath, err)
	}

	manifest.Deployment.File = functionPath

	// Store the function record.
	fn := blockless.FunctionRecord{
		CID:      cid,
		URL:      address,
		Manifest: manifest,
		Archive:  functionPath,
		Files:    out,
	}
	err = f.saveFunction(ctx, fn)
	if err != nil {
		f.log.Error().
			Err(err).
			Str("cid", cid).
			Msg("could not save function record")
	}

	f.log.Debug().
		Str("cid", cid).
		Str("address", address).
		Msg("installed function")

	return nil
}

// Installed checks if the function with the given CID is installed.
func (f *FStore) IsInstalled(cid string) (bool, error) {

	fn, err := f.getFunction(context.Background(), cid)
	if err != nil && errors.Is(err, blockless.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not get function from store: %w", err)
	}

	haveArchive, haveFiles, err := f.checkFunctionFiles(fn)
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
