package response

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
)

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the response to the `MessageExecute` message.
type Execute struct {
	blockless.BaseMessage
	RequestID string            `json:"request_id,omitempty"`
	Code      codes.Code        `json:"code,omitempty"`
	Results   execute.ResultMap `json:"results,omitempty"`
	Cluster   execute.Cluster   `json:"cluster,omitempty"`

	PBFT PBFTResultInfo `json:"pbft,omitempty"`
	// Signed digest of the response.
	Signature string `json:"signature,omitempty"`

	// Used to communicate the reason for failure to the user.
	ErrorMessage string `json:"message,omitempty"`
}

func (e *Execute) WithResults(r execute.ResultMap) *Execute {
	e.Results = r
	return e
}

func (e *Execute) WithCluster(c execute.Cluster) *Execute {
	e.Cluster = c
	return e
}

func (e *Execute) WithErrorMessage(err error) *Execute {
	e.ErrorMessage = err.Error()
	return e
}

func (Execute) Type() string { return blockless.MessageExecuteResponse }

func (e Execute) MarshalJSON() ([]byte, error) {
	type Alias Execute
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(e),
		Type:  e.Type(),
	}
	return json.Marshal(rec)
}

type PBFTResultInfo struct {
	View             uint      `json:"view"`
	RequestTimestamp time.Time `json:"request_timestamp,omitempty"`
	Replica          peer.ID   `json:"replica,omitempty"`
}

func (e *Execute) Sign(key crypto.PrivKey) error {

	cp := *e
	// Exclude some of the fields from the signature.
	cp.Signature = ""
	cp.BaseMessage = blockless.BaseMessage{}

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

	cp := e
	// Exclude some of the fields from the signature.
	cp.Signature = ""
	cp.BaseMessage = blockless.BaseMessage{}

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
