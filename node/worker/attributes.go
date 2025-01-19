package worker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ipfs/boxo/ipns"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/execute"
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

func haveAttributes(have attributes.Attestation, want execute.Attributes) error {

	if want.AttestationRequired && len(have.Attestors) == 0 {
		return errors.New("attestors required but none found")
	}

	// If we need to check attestors, create a map of them now.
	var attestors map[peer.ID]struct{}
	if len(want.Attestors.Each) > 0 || len(want.Attestors.OneOf) > 0 {
		attestors = make(map[peer.ID]struct{}, len(have.Attestors))
		for _, attestor := range have.Attestors {
			attestors[attestor.Signer] = struct{}{}
		}
	}

	// If the client wants specific attestors, check if they're present.
	if len(want.Attestors.Each) > 0 {
		for _, wa := range want.Attestors.Each {
			_, ok := attestors[wa]
			if !ok {
				return fmt.Errorf("attestor %s explicitly requested but not found", wa.String())
			}
		}
	}

	// If the client wants some of these attestors, check if at least one if found.
	if len(want.Attestors.OneOf) > 0 {
		var found bool
		for _, wa := range want.Attestors.OneOf {
			_, ok := attestors[wa]
			if ok {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("at least one attestor wanted but none found (wanted: %s)", bls.PeerIDsToStr(want.Attestors.OneOf))
		}
	}

	// It doesn't make a lot of sense to require attestors without wanting specific attributes,
	// but if that's the case, and there's no attributes wanted, we're done now.
	if len(want.Values) == 0 {
		return nil
	}

	attrs := make(map[string]string, len(have.Attributes))
	for _, attr := range have.Attributes {
		attrs[attr.Name] = attr.Value
	}

	for _, wantAttr := range want.Values {

		value, ok := attrs[wantAttr.Name]
		if !ok {
			return fmt.Errorf("attribute wanted but not found (attr: %v, value: %v)", wantAttr.Name, wantAttr.Value)
		}

		if value != wantAttr.Value {
			return fmt.Errorf("attribute wanted but value doesn't match (attr: %v, want: %v, have: %v)", wantAttr.Name, wantAttr.Value, value)
		}
	}

	return nil
}
