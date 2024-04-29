package fstore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

// Sync will verify that the function identified by `cid` is still found on the local filesystem.
// If the function archive of function files are missing, they will be recreated.
func (h *FStore) Sync(cid string) error {

	h.log.Debug().Str("cid", cid).Msg("checking function installation")

	// Read the function directly from storage - we don't want to update the timestamp
	// since this is a 'maintenance' access.
	var fn blockless.FunctionRecord
	err := h.store.GetRecord(cid, &fn)
	if err != nil {
		return fmt.Errorf("could not get function record: %w", err)
	}

	h.log.Debug().
		Str("cid", cid).
		Str("archive", fn.Archive).
		Str("files", fn.Files).
		Msg("checking function installation")

	haveArchive, haveFiles, err := h.checkFunctionFiles(fn)
	if err != nil {
		return fmt.Errorf("could not verify function cache: %w", err)
	}

	// If both archive and files are there - we're done.
	if haveArchive && haveFiles {
		h.log.Debug().Str("cid", cid).Msg("function files found, done")
		return nil
	}

	h.log.Debug().
		Bool("have_archive", haveArchive).
		Bool("have_files", haveFiles).
		Str("cid", cid).
		Msg("function installation missing files")

	// If we don't have the archive - redownload it.
	if !haveArchive {
		path, err := h.download(cid, fn.Manifest)
		if err != nil {
			return fmt.Errorf("could not download the function archive (cid: %v): %w", cid, err)
		}

		// Update path in case it changed.
		fn.Archive = h.cleanPath(path)
	}

	// If we don't have the files OR if we redownloaded the archive - recreate the files.
	if !haveFiles || !haveArchive {

		archivePath := filepath.Join(h.workdir, fn.Archive)
		files := filepath.Join(h.workdir, fn.CID)

		h.log.Info().
			Str("archive", archivePath).
			Str("fn_archive", fn.Archive).
			Msg("archive path to use")

		err = h.unpackArchive(archivePath, files)
		if err != nil {
			return fmt.Errorf("could not unpack gzip archive (cid: %v, file: %s): %w", cid, fn.Archive, err)
		}

		fn.Files = files
	}

	// Save the updated function record.
	err = h.saveFunction(fn)
	if err != nil {
		return fmt.Errorf("could not save function (cid: %v): %w", cid, err)
	}

	return nil
}

// checkFunctionFiles checks if the files required by the function are found on local storage.
// It returns two booleans indicating presence of the archive file, the unpacked files, and a potential error.
func (h *FStore) checkFunctionFiles(fn blockless.FunctionRecord) (bool, bool, error) {

	// Check if the archive is found.
	archiveFound := true

	apath := filepath.Join(h.workdir, fn.Archive)
	_, err := os.Stat(apath)
	if err != nil && os.IsNotExist(err) {
		archiveFound = false
	} else if err != nil {
		return false, false, fmt.Errorf("could not stat function archive: %w", err)
	}

	// NOTE: We could check that it's a regular file (plus cheksum), but lets not go overboard for now.

	// Check if the files are found.
	filesFound := true

	fpath := filepath.Join(h.workdir, fn.Files)
	_, err = os.Stat(fpath)
	if err != nil && os.IsNotExist(err) {
		filesFound = false
	} else if err != nil {
		return false, false, fmt.Errorf("could not stat function files: %w", err)
	}

	return archiveFound, filesFound, nil
}
