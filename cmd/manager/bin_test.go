package main

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallB7s(t *testing.T) {
	baseURL := "https://github.com/blocklessnetwork/b7s/releases/download"
	version := "v0.0.25"

	// Run the installBlsCLI function
	installB7s(baseURL, version)

	// Check if the b7s binary was installed
	usr, _ := user.Current()
	binPath := filepath.Join(usr.HomeDir, ".b7s", "bin")
	b7sPath := filepath.Join(binPath, "b7s")

	_, err := os.Stat(b7sPath)
	assert.NoError(t, err)
}

func TestRemoveB7s(t *testing.T) {
	// Run the removeB7s function
	removeB7s()

	// Check if the b7s folder was removed
	usr, _ := user.Current()
	b7sPath := filepath.Join(usr.HomeDir, ".b7s")

	_, err := os.Stat(b7sPath)
	assert.True(t, os.IsNotExist(err))
}
