package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	// Names used for created files.
	privKeyName  = "priv.bin"
	pubKeyName   = "pub.bin"
	identityName = "identity"
)

const (
	// Permissions used for created files.
	privKeyPermissions = 0600
	pubKeyPermissions  = 0644
)

func main() {

	var (
		flagOutputDir string
	)

	pflag.StringVarP(&flagOutputDir, "output", "o", ".", "directory where keys should be stored")

	pflag.Parse()

	// Create output directory, if it doesn't exist.
	err := os.MkdirAll(flagOutputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("could not create output directory: %s", err)
	}

	// Generate key pair.
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	if err != nil {
		log.Fatalf("could not generate key pair: %s", err)
	}

	// Encode keys and extract peer ID from key.
	privPayload, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		log.Fatalf("could not marshal private key: %s", err)
	}

	pubPayload, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		log.Fatalf("could not marshal public key: %s", err)
	}

	identity, err := peer.IDFromPublicKey(pub)
	if err != nil {
		log.Fatalf("failed to get peer identity from public key: %s", err)
	}

	// Write keys and identity to files.

	pubKeyFile := filepath.Join(flagOutputDir, pubKeyName)
	err = os.WriteFile(pubKeyFile, pubPayload, pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write private key to file: %s", err)
	}

	idFile := filepath.Join(flagOutputDir, identityName)
	err = os.WriteFile(idFile, []byte(identity), pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write private key to file: %s", err)
	}

	privKeyFile := filepath.Join(flagOutputDir, privKeyName)
	err = os.WriteFile(privKeyFile, privPayload, privKeyPermissions)
	if err != nil {
		log.Fatalf("could not write private key to file: %s", err)
	}

	fmt.Printf("generated private key: %s\n", privKeyFile)
	fmt.Printf("generated public key: %s\n", pubKeyFile)
	fmt.Printf("generated identity file: %s\n", idFile)
}
