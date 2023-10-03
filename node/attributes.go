package node

import (
	"fmt"
	"net/http"

	"github.com/ipfs/boxo/ipns"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

const (
	defaultAttributesFilename = "attributes.bin"
)

func loadAttributes(key crypto.PubKey) (attributes.Attestation, error) {

	name, err := getAttributesIPNSName(key)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not get name from key: %w", err)
	}

	attributeURL := ipnsGatewayURL(name)

	res, err := http.Get(attributeURL)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not get attribute file from URL: %w", err)
	}
	defer res.Body.Close()

	att, err := attributes.ImportAttestation(res.Body)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not load attestation from file: %w", err)
	}

	return att, nil
}

func getAttributesIPNSName(key crypto.PubKey) (string, error) {

	id, err := peer.IDFromPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("could not get peer ID for key: %w", err)
	}

	name := ipns.NameFromPeer(id)
	return name.String(), nil
}

func ipnsGatewayURL(name string) string {
	return fmt.Sprintf("https://%s.ipns.cf-ipfs.com/%s", name, defaultAttributesFilename)
}
