package pbft

import (
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r *Replica) sign(rec signable) error {

	digest := getDigest(rec.signableRecord())
	sig, err := r.host.PrivateKey().Sign([]byte(digest))
	if err != nil {
		return fmt.Errorf("could not sign preprepare: %w", err)
	}

	rec.setSignature(hex.EncodeToString(sig))

	return nil
}

func (r *Replica) verifySignature(rec signable, signer peer.ID) error {

	// Get the digest of the message, excluding the signature.
	digest := getDigest(rec.signableRecord())

	pub, err := signer.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf("could not extract public key from peer ID (id: %s): %w", signer.String(), err)
	}

	sig, err := hex.DecodeString(rec.getSignature())
	if err != nil {
		return fmt.Errorf("could not decode signature from hex: %w", err)
	}

	ok, err := pub.Verify([]byte(digest), sig)
	if err != nil {
		return fmt.Errorf("could not verify signature: %w", err)
	}

	if !ok {
		return ErrInvalidSignature
	}

	return nil
}
