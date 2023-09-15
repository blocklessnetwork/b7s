package main

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// LoadOrCreateKeys loads existing keys or creates new ones if not present
func LoadOrCreateKeys(privKeyFile string, outputDir string) (crypto.PrivKey, crypto.PubKey, error) {
	var priv crypto.PrivKey
	var pub crypto.PubKey
	var err error

	if _, err := os.Stat(privKeyFile); os.IsNotExist(err) {
		priv, pub, err = crypto.GenerateKeyPair(crypto.Ed25519, 0)
		if err != nil {
			return nil, nil, err
		}
	} else {
		privBytes, err := os.ReadFile(privKeyFile)
		if err != nil {
			return nil, nil, err
		}

		priv, err = crypto.UnmarshalPrivateKey(privBytes)
		if err != nil {
			return nil, nil, err
		}
		pub = priv.GetPublic()
	}

	privPayload, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		log.Fatalf("Could not marshal private key: %s", err)
	}

	pubPayload, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		log.Fatalf("Could not marshal public key: %s", err)
	}

	identity, err := peer.IDFromPublicKey(pub)
	if err != nil {
		log.Fatalf("Could not generate identity: %s", err)
	}

	pubKeyFile := filepath.Join(outputDir, pubKeyName)
	err = os.WriteFile(pubKeyFile, pubPayload, pubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write public key to file: %s", err)
	}

	pubKeyTextFile := filepath.Join(outputDir, pubKeyTxtName)
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubPayload)
	err = os.WriteFile(pubKeyTextFile, []byte(pubKeyBase64), pubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write public key text to file: %s", err)
	}

	identityFile := filepath.Join(outputDir, identityName)
	err = os.WriteFile(identityFile, []byte(identity.Pretty()), pubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write identity to file: %s", err)
	}

	privKeyFile = filepath.Join(outputDir, privKeyName)
	err = os.WriteFile(privKeyFile, privPayload, privKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write private key to file: %s", err)
	}

	// Write peer ID to file
	peerIDFile := filepath.Join(outputDir, peerIDFileName)
	err = os.WriteFile(peerIDFile, []byte(identity.Pretty()), pubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write peer ID to file: %s", err)
	}

	return priv, pub, nil
}
