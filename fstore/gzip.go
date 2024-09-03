package fstore

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (f *FStore) unpackArchive(filename string, destination string) error {

	// Use CWD if not specified.
	if destination == "" {
		destination = "."
	}

	f.log.Debug().
		Str("archive", filename).
		Str("destination", destination).
		Msg("unpacking gzip archive")

	// Create output directory.
	err := os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create destination directory (dir: %s): %w", destination, err)
	}

	// Open gzip archive.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open gzip archive (file: %s): %w", filename, err)
	}
	defer file.Close()

	// Create reader for compressed data.
	reader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {

		// Get the next record from the archive.
		entry, err := tarReader.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("could not read archive: %w", err)
			}

			break
		}

		typ := entry.Typeflag

		f.log.Debug().
			Str("archive", filename).
			Str("entry", entry.Name).
			Str("type", fmt.Sprintf("%d", typ)).
			Msg("processing archive entry")

		switch typ {
		case tar.TypeDir:
			// Entry is a directory - create output dir.
			dir := filepath.Join(destination, entry.Name)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("could not create directory (dir: %s): %w", dir, err)
			}

		case tar.TypeReg:
			// Entry is a file - create file.
			file := filepath.Join(destination, entry.Name)
			of, err := os.Create(file)
			if err != nil {
				return fmt.Errorf("could not create file (file: %s): %w", file, err)
			}

			// Copy file content.
			_, err = io.Copy(of, tarReader)
			of.Close()
			if err != nil {
				return fmt.Errorf("could not write file content (file: %s): %w", file, err)
			}

		default:
			return fmt.Errorf("unexpected entry found (name: %s, type: %d)", entry.Name, typ)
		}
	}

	f.log.Debug().
		Str("archive", filename).
		Str("destination", destination).
		Msg("gzip archive unpacked")

	return nil
}
