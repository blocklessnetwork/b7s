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

func installBinary(url, folder string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	targetPath := filepath.Join(usr.HomeDir, folder)
	os.MkdirAll(targetPath, os.ModePerm)

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

		path := filepath.Join(targetPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				log.Fatal(err)
			}

		case tar.TypeReg:
			outFile, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatal(err)
			}

			outFile.Close()

			if err := os.Chmod(path, os.FileMode(header.Mode)); err != nil {
				log.Fatal(err)
			}

			log.Printf("File %s installed in %s", header.Name, targetPath)
		}
	}
}

func installB7s(baseURL, version string) {
	arch := runtime.GOARCH
	platform := runtime.GOOS
	url := fmt.Sprintf("%s/%s/b7s-%s.%s.tar.gz", baseURL, version, platform, arch)
	installBinary(url, ".b7s/networking")
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

func installRuntime(baseURL, version string) {
	arch := runtime.GOARCH
	platform := runtime.GOOS

	if platform == "darwin" {
		platform = "macos"
	}

	if arch == "amd64" {
		arch = "x86_64"
	}

	if arch == "arm64" {
		arch = "aarch64"
	}

	url := fmt.Sprintf("%s/%s/blockless-runtime.%s-latest.%s.tar.gz", baseURL, version, platform, arch)
	installBinary(url, ".b7s/runtime")
}
