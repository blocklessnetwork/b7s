package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
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
	// ECDSA constants
	ecdsaPrivKeyName        = "ecdsaPrivKey.bin"
	ecdsaPubKeyName         = "ecdsaPubKey.bin"
	ecdsaPubKeyTxtName      = "ecdsaPubKey.txt"
	ecdsaIdentityName       = "ecdsaIdentity"
	ecdsaPeerIDFileName     = "ecdsaPeerID.txt"
	ecdsaPrivKeyPermissions = 0600
	ecdsaPubKeyPermissions  = 0644
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
		flagUseECDSA bool
	)

	pflag.StringVar(&flagPeerID, "peerid", "", "PeerID for verification")
	pflag.StringVarP(&flagOutputDir, "output", "o", ".", "directory where keys should be stored")
	pflag.StringVarP(&flagString, "string", "s", "", "string to sign and verify")
	pflag.StringVarP(&flagFile, "file", "f", "", "file to sign and verify")
	pflag.StringVar(&flagPublicKey, "pubkey", "", "Base64 encoded public key for verification")
	pflag.StringVar(&flagMessage, "message", "", "The original message to verify")
	pflag.StringVar(&flagSignature, "signature", "", "Base64 encoded signature to verify")
	pflag.BoolVar(&flagUseECDSA, "ecdsa", false, "Use ECDSA keys instead of the default")

	pflag.Parse()

	// Initialize output directory
	err := os.MkdirAll(flagOutputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create output directory: %s", err)
	}

	var priv, pub interface{}  // use interface{} if the key types are different
	var privKeyFile string

	if flagUseECDSA {
		privKeyFile = filepath.Join(flagOutputDir, ecdsaPrivKeyName)
		priv, pub, err = LoadOrCreateECDSAKeys(privKeyFile, flagOutputDir)
	} else {
		privKeyFile = filepath.Join(flagOutputDir, privKeyName)
		priv, pub, err = LoadOrCreateKeys(privKeyFile, flagOutputDir)
	}
	
	if err != nil {
		log.Fatalf("Error loading or creating keys: %s", err)
	}

	if flagString != "" || flagFile != "" {
		if privKey, ok := priv.(crypto.PrivKey); ok {
			if pubKey, ok := pub.(crypto.PubKey); ok {
				HandleSignAndVerify(privKey, pubKey, flagString, flagFile, flagOutputDir)
			} else {
				log.Fatal("pub is not a valid crypto.PubKey")
			}
		} else {
			log.Fatal("priv is not a valid crypto.PrivKey")
		}
	}

	if flagPublicKey != "" && flagMessage != "" && flagSignature != "" {
		VerifyGivenSignature(flagPublicKey, flagMessage, flagSignature)
	}

	if flagPeerID != "" && flagMessage != "" && flagSignature != "" {
		VerifyGivenSignatureWithPeerID(flagPeerID, flagMessage, flagSignature)
	}
}
