package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// LoadOrCreateECDSAKeys loads existing ECDSA keys or creates new ones if not present
func LoadOrCreateECDSAKeys(ecdsaPrivKeyFile string, outputDir string) (crypto.PrivKey, crypto.PubKey, error) {
	var priv crypto.PrivKey
	var pub crypto.PubKey
	var err error

	if _, err := os.Stat(ecdsaPrivKeyFile); os.IsNotExist(err) {
		priv, pub, err = crypto.GenerateKeyPair(crypto.ECDSA, 256)
		if err != nil {
			return nil, nil, err
		}
	} else {
		privBytes, err := ioutil.ReadFile(ecdsaPrivKeyFile)
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
		log.Fatalf("Could not marshal private ECDSA key: %s", err)
	}

	pubPayload, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		log.Fatalf("Could not marshal public ECDSA key: %s", err)
	}

	identity, err := peer.IDFromPublicKey(pub)
	if err != nil {
		log.Fatalf("Could not generate identity: %s", err)
	}

	ecdsaPubKeyFile := filepath.Join(outputDir, ecdsaPubKeyName)
	err = ioutil.WriteFile(ecdsaPubKeyFile, pubPayload, ecdsaPubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write public ECDSA key to file: %s", err)
	}

	ecdsaPubKeyTextFile := filepath.Join(outputDir, ecdsaPubKeyTxtName)
	ecdsaPubKeyBase64 := base64.StdEncoding.EncodeToString(pubPayload)
	err = ioutil.WriteFile(ecdsaPubKeyTextFile, []byte(ecdsaPubKeyBase64), ecdsaPubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write public ECDSA key text to file: %s", err)
	}

	ecdsaIdentityFile := filepath.Join(outputDir, ecdsaIdentityName)
	err = ioutil.WriteFile(ecdsaIdentityFile, []byte(identity.Pretty()), ecdsaPubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write identity to file: %s", err)
	}

	ecdsaPrivKeyFilePath := filepath.Join(outputDir, ecdsaPrivKeyName)
	err = ioutil.WriteFile(ecdsaPrivKeyFilePath, privPayload, ecdsaPrivKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write private ECDSA key to file: %s", err)
	}

	ecdsaPeerIDFile := filepath.Join(outputDir, ecdsaPeerIDFileName)
	err = ioutil.WriteFile(ecdsaPeerIDFile, []byte(identity.Pretty()), ecdsaPubKeyPermissions)
	if err != nil {
		log.Fatalf("Could not write peer ID to file: %s", err)
	}

	return priv, pub, nil
}
