package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	privKeyName    = "priv.bin"
	pubKeyName     = "pub.bin"
	pubKeyTxtName  = "pub.txt"
	identityName   = "identity"
	peerIDFileName = "peerid.txt"  // New constant for the peer ID file
)

const (
	privKeyPermissions = 0600
	pubKeyPermissions  = 0644
)

func main() {
	var flagOutputDir string

	pflag.StringVarP(&flagOutputDir, "output", "o", ".", "directory where keys should be stored")
	pflag.Parse()

	err := os.MkdirAll(flagOutputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("could not create output directory: %s", err)
	}

	privKeyFile := filepath.Join(flagOutputDir, privKeyName)

	var priv crypto.PrivKey
	var pub crypto.PubKey

	if _, err := os.Stat(privKeyFile); os.IsNotExist(err) {
		priv, pub, err = crypto.GenerateKeyPair(crypto.Ed25519, 0)
		if err != nil {
			log.Fatalf("could not generate key pair: %s", err)
		}
	} else {
		privBytes, err := ioutil.ReadFile(privKeyFile)
		if err != nil {
			log.Fatalf("could not read existing private key: %s", err)
		}

		priv, err = crypto.UnmarshalPrivateKey(privBytes)
		if err != nil {
			log.Fatalf("could not unmarshal private key: %s", err)
		}
		pub = priv.GetPublic()
	}

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

	pubKeyFile := filepath.Join(flagOutputDir, pubKeyName)
	err = os.WriteFile(pubKeyFile, pubPayload, pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write public key to file: %s", err)
	}

	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubPayload)
	pubKeyTxtFile := filepath.Join(flagOutputDir, pubKeyTxtName)
	err = os.WriteFile(pubKeyTxtFile, []byte(pubKeyBase64), pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write public key string to file: %s", err)
	}

	idFile := filepath.Join(flagOutputDir, identityName)
	err = os.WriteFile(idFile, []byte(identity), pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write identity to file: %s", err)
	}

	// Write the Peer ID to a file
	peerIDFile := filepath.Join(flagOutputDir, peerIDFileName)
	err = os.WriteFile(peerIDFile, []byte(identity.Pretty()), pubKeyPermissions)
	if err != nil {
		log.Fatalf("could not write peer ID to file: %s", err)
	}

	err = os.WriteFile(privKeyFile, privPayload, privKeyPermissions)
	if err != nil {
		log.Fatalf("could not write private key to file: %s", err)
	}

	fmt.Printf("private key file: %s\n", privKeyFile)
	fmt.Printf("public key file: %s\n", pubKeyFile)
	fmt.Printf("public key text: %s\n", pubKeyTxtFile)
	fmt.Printf("identity file: %s\n", idFile)
	fmt.Printf("peer ID file: %s\n", peerIDFile)  // Output the Peer ID file path
}
