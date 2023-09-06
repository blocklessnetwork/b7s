package execute

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func (e *Request) Sign(key crypto.PrivKey) error {

	cp := *e
	e.Signature = ""

	payload, err := json.Marshal(cp)
	if err != nil {
		return fmt.Errorf("could not get byte representation of the record: %w", err)
	}

	sig, err := key.Sign(payload)
	if err != nil {
		return fmt.Errorf("could not sign digest: %w", err)
	}

	e.Signature = hex.EncodeToString(sig)
	return nil
}

func (e Request) VerifySignature(key crypto.PubKey) error {

	// Exclude signature and the `from` field from the signature.
	cp := e
	cp.Signature = ""

	payload, err := json.Marshal(cp)
	if err != nil {
		return fmt.Errorf("could not get byte representation of the record: %w", err)
	}

	sig, err := hex.DecodeString(e.Signature)
	if err != nil {
		return fmt.Errorf("could not decode signature from hex: %w", err)
	}

	ok, err := key.Verify(payload, sig)
	if err != nil {
		return fmt.Errorf("could not verify signature: %w", err)
	}

	if !ok {
		return errors.New("invalid signature")
	}

	return nil
}
