package response

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

// Execute describes the response to the `MessageExecute` message.
type Execute struct {
	Type      string            `json:"type,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	From      peer.ID           `json:"from,omitempty"`
	Code      codes.Code        `json:"code,omitempty"`
	Results   execute.ResultMap `json:"results,omitempty"`
	Cluster   execute.Cluster   `json:"cluster,omitempty"`

	PBFT PBFTResultInfo `json:"pbft,omitempty"`
	// Signed digest of the response.
	Signature string `json:"signature,omitempty"`

	// Used to communicate the reason for failure to the user.
	Message string `json:"message,omitempty"`
}

type PBFTResultInfo struct {
	View             uint      `json:"view"`
	RequestTimestamp time.Time `json:"request_timestamp,omitempty"`
	Replica          peer.ID   `json:"replica"`
}

func (e *Execute) Sign(key crypto.PrivKey) error {

	// Exclude signature and the `from` field from the signature.
	cp := *e
	cp.Signature = ""
	cp.From = ""

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

func (e Execute) VerifySignature(key crypto.PubKey) error {

	// Exclude signature and the `from` field from the signature.
	cp := e
	cp.Signature = ""
	cp.From = ""

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
