package keygen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

var (
	bits int
)

func GenerateKeys(outputFolder string) error {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exPath := filepath.Dir(ex)

	privKey, pubKey, err := crypto.GenerateKeyPair(
		crypto.Ed25519,
		bits,
	)
	if err != nil {
		log.Fatal("failed to generate key:" + err.Error())
	}

	privBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		log.Fatal("failed to marshal private key:" + err.Error())
	}

	pubBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		log.Fatal("failed to marshal public key:" + err.Error())

	}

	identity, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		log.Fatal("failed to get peer identity from public key:" + err.Error())
	}

	if err := os.MkdirAll(exPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(exPath)
	err = os.MkdirAll(dir, 0777)

	if err != nil {
		log.Fatal("failed to write to folder" + err.Error())
	}

	pubKeyFile := fmt.Sprintf("%s/pub.bin", exPath)
	privKeyFile := fmt.Sprintf("%s/priv.bin", exPath)
	peerIdFile := fmt.Sprintf("%s/identity", exPath)

	if err := ioutil.WriteFile(pubKeyFile, pubBytes, 0644); err != nil {
		log.Fatal("failed to save pub key to file:" + err.Error())
	}

	if err := ioutil.WriteFile(privKeyFile, privBytes, 0644); err != nil {
		log.Fatal("failed to save private key to file:" + err.Error())
	}

	if err := ioutil.WriteFile(peerIdFile, []byte(identity.String()), 0644); err != nil {
		log.Fatal("failed to save identity to file:" + err.Error())
	}

	log.Info("Keys are generated at: ", exPath)
	log.Info("identity:", identity)

	return nil
}
