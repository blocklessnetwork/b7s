package main

import (
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func peerIDFromKey(keyPath string) (string, error) {

	key, err := readPrivateKey(keyPath)
	if err != nil {
		log.Fatalf("could not read key file: %s", err)
	}

	id, err := peer.IDFromPrivateKey(key)
	if err != nil {
		log.Fatalf("could not determine identity: %s", err)
	}

	return id.String(), nil
}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}
