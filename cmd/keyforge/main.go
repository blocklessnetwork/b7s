package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

const (
	privKeyName        = "priv.bin"
	pubKeyName         = "pub.bin"
	pubKeyTxtName      = "pubkey.txt"
	identityName       = "identity"
	peerIDFileName     = "peerid.txt"
	privKeyPermissions = 0600
	pubKeyPermissions  = 0644
)

func main() {
	var (
		flagOutputDir string
		flagString    string
		flagFile      string
		flagPublicKey string
		flagMessage   string
		flagSignature string
		flagPeerID    string
	)

	pflag.StringVar(&flagPeerID, "peerid", "", "PeerID for verification")
	pflag.StringVarP(&flagOutputDir, "output", "o", ".", "directory where keys should be stored")
	pflag.StringVarP(&flagString, "string", "s", "", "string to sign and verify")
	pflag.StringVarP(&flagFile, "file", "f", "", "file to sign and verify")
	pflag.StringVar(&flagPublicKey, "pubkey", "", "Base64 encoded public key for verification")
	pflag.StringVar(&flagMessage, "message", "", "The original message to verify")
	pflag.StringVar(&flagSignature, "signature", "", "Base64 encoded signature to verify")

	pflag.Parse()

	// Initialize output directory
	err := os.MkdirAll(flagOutputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create output directory: %s", err)
	}

	privKeyFile := filepath.Join(flagOutputDir, privKeyName)

	priv, pub, err := LoadOrCreateKeys(privKeyFile, flagOutputDir)
	if err != nil {
		log.Fatalf("Error loading or creating keys: %s", err)
	}

	if flagString != "" || flagFile != "" {
		HandleSignAndVerify(priv, pub, flagString, flagFile, flagOutputDir)
	}

	if flagPublicKey != "" && flagMessage != "" && flagSignature != "" {
		VerifyGivenSignature(flagPublicKey, flagMessage, flagSignature)
	}

	if flagPeerID != "" && flagMessage != "" && flagSignature != "" {
		VerifyGivenSignatureWithPeerID(flagPeerID, flagMessage, flagSignature)
	}
}
