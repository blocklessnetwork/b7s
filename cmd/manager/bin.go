package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)


func installBlsCLI(baseURL string, version string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	binPath := filepath.Join(usr.HomeDir, ".b7s", "bin")
	os.MkdirAll(binPath, os.ModePerm)

	arch := runtime.GOARCH
	platform := runtime.GOOS

	// maybe change this in ci
	if platform == "darwin" {
		platform = "macOS"
	}

	url := fmt.Sprintf("%s/%s/bls-%s-%s-blockless-cli.tar.gz", baseURL, version, platform, arch)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	archiveData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(archiveData))
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if header.Typeflag == tar.TypeReg {
			target := filepath.Join(binPath, "b7s")
			outFile, err := os.Create(target)
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatal(err)
			}

			if err := os.Chmod(target, 0755); err != nil {
				log.Fatal(err)
			}

			log.Printf("b7s CLI installed in %s", binPath)
			break
		}
	}
}



func removeB7s() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	b7sPath := filepath.Join(usr.HomeDir, ".b7s")
	err = os.RemoveAll(b7sPath)
	if err != nil {
		log.Println("Error removing b7s:", err)
	}

	log.Println("b7s removed.")
}

