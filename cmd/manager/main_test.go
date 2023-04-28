package main

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestInstallBlsCLI(t *testing.T) {
	// Test case: Linux x64
	url := "https://github.com/blocklessnetwork/cli/releases/download/0.0.46/bls-linux-x64-blockless-cli.tar.gz"

	installBlsCLI(url)

	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(usr.HomeDir, ".b7s", "bin", "bls")

	// Check if the b7s CLI binary exists
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatalf("b7s CLI not installed in %s", binPath)
	}
}
