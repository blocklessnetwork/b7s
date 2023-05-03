package main

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestInstallBlsCLI(t *testing.T) {
	// Test case: Linux x64
	baseURL := "https://github.com/blocklessnetwork/cli/releases/download"
	version := "latest"
	installBlsCLI(baseURL, version)

	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(usr.HomeDir, ".b7s", "bin", "b7s")

	// Check if the b7s CLI binary exists
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatalf("b7s CLI not installed in %s", binPath)
	}
}
