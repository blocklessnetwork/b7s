package pbft

import (
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// messageRecord is used as an interim format to supplement the original type with its type.
// Useful for serialization to automatically include the message type field.
type messageRecord struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}

// messageEnvelope is used as an interim format to extract the original type from the `messageRecord` format.
type messageEnvelope struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (r Request) MarshalJSON() ([]byte, error) {
	type alias Request
	rec := messageRecord{
		Type: MessageRequest,
		Data: alias(r),
	}
	return json.Marshal(rec)
}

func (r *Request) UnmarshalJSON(data []byte) error {
	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	type alias *Request
	return json.Unmarshal(rec.Data, alias(r))
}

func (p PrePrepare) MarshalJSON() ([]byte, error) {
	type alias PrePrepare
	rec := messageRecord{
		Type: MessagePrePrepare,
		Data: alias(p),
	}
	return json.Marshal(rec)
}

func (p *PrePrepare) UnmarshalJSON(data []byte) error {
	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	type alias *PrePrepare
	return json.Unmarshal(rec.Data, alias(p))
}

func (p Prepare) MarshalJSON() ([]byte, error) {
	type alias Prepare
	rec := messageRecord{
		Type: MessagePrepare,
		Data: alias(p),
	}
	return json.Marshal(rec)
}

func (p *Prepare) UnmarshalJSON(data []byte) error {
	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	type alias *Prepare
	return json.Unmarshal(rec.Data, alias(p))
}

func (c Commit) MarshalJSON() ([]byte, error) {
	type alias Commit
	rec := messageRecord{
		Type: MessageCommit,
		Data: alias(c),
	}
	return json.Marshal(rec)
}

func (c *Commit) UnmarshalJSON(data []byte) error {
	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	type alias *Commit
	return json.Unmarshal(rec.Data, alias(c))
}

type prepareInfoEncoded struct {
	View           uint               `json:"view"`
	SequenceNumber uint               `json:"sequence_number"`
	Digest         string             `json:"digest"`
	PrePrepare     PrePrepare         `json:"preprepare"`
	Prepares       map[string]Prepare `json:"prepares"`
}

func (p PrepareInfo) MarshalJSON() ([]byte, error) {

	encodedPrepareMap := make(map[string]Prepare)
	for id, prepare := range p.Prepares {
		encodedPrepareMap[id.String()] = prepare
	}

	rec := prepareInfoEncoded{
		View:           p.View,
		SequenceNumber: p.SequenceNumber,
		Digest:         p.Digest,
		PrePrepare:     p.PrePrepare,
		Prepares:       encodedPrepareMap,
	}

	return json.Marshal(rec)
}

func (p *PrepareInfo) UnmarshalJSON(data []byte) error {

	var info prepareInfoEncoded
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}

	prepareMap := make(map[peer.ID]Prepare)
	for idStr, prepare := range info.Prepares {
		id, err := peer.Decode(idStr)
		if err != nil {
			return fmt.Errorf("could not decode peer.ID (str: %s): %w", idStr, err)
		}
		prepareMap[id] = prepare
	}

	*p = PrepareInfo{
		View:           info.View,
		SequenceNumber: info.SequenceNumber,
		Digest:         info.Digest,
		PrePrepare:     info.PrePrepare,
		Prepares:       prepareMap,
	}

	return nil
}

func (v ViewChange) MarshalJSON() ([]byte, error) {
	type alias ViewChange
	rec := messageRecord{
		Type: MessageViewChange,
		Data: alias(v),
	}
	return json.Marshal(rec)
}

func (v *ViewChange) UnmarshalJSON(data []byte) error {
	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	type alias *ViewChange
	return json.Unmarshal(rec.Data, alias(v))
}

type newViewEncode struct {
	View        uint                  `json:"view"`
	Messages    map[string]ViewChange `json:"messages"`
	PrePrepares []PrePrepare          `json:"preprepares"`
}

func (v NewView) MarshalJSON() ([]byte, error) {

	// To properly handle `peer.ID` serialization, this is a bit more involved.
	// See documentation for `ResultMap.MarshalJSON` in `models/execute/response.go`.
	messages := make(map[string]ViewChange)
	for replica, vc := range v.Messages {
		messages[replica.String()] = vc
	}

	nv := newViewEncode{
		View:        v.View,
		Messages:    messages,
		PrePrepares: v.PrePrepares,
	}

	rec := messageRecord{
		Type: MessageNewView,
		Data: nv,
	}

	return json.Marshal(rec)
}

func (n *NewView) UnmarshalJSON(data []byte) error {

	var rec messageEnvelope
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}

	var nv newViewEncode
	err = json.Unmarshal(rec.Data, &nv)
	if err != nil {
		return err
	}

	messages := make(map[peer.ID]ViewChange)
	for idStr, vc := range nv.Messages {
		id, err := peer.Decode(idStr)
		if err != nil {
			return fmt.Errorf("could not decode peer.ID (str: %s): %w", idStr, err)
		}
		messages[id] = vc
	}

	*n = NewView{
		View:        nv.View,
		Messages:    messages,
		PrePrepares: nv.PrePrepares,
	}

	return nil
}

func unpackMessage(payload []byte) (any, error) {

	var msg messageEnvelope
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return nil, fmt.Errorf("could not unpack base message: %w", err)
	}

	switch msg.Type {
	case MessageRequest:
		var request Request
		err = json.Unmarshal(payload, &request)
		if err != nil {
			return nil, fmt.Errorf("could not unpack request: %w", err)
		}
		return request, nil

	case MessagePrePrepare:
		var preprepare PrePrepare
		err = json.Unmarshal(payload, &preprepare)
		if err != nil {
			return nil, fmt.Errorf("could not unpack pre-prepare message: %w", err)
		}
		return preprepare, nil

	case MessagePrepare:
		var prepare Prepare
		err = json.Unmarshal(payload, &prepare)
		if err != nil {
			return nil, fmt.Errorf("could not unpack prepare message: %w", err)
		}
		return prepare, nil

	case MessageCommit:
		var commit Commit
		err = json.Unmarshal(payload, &commit)
		if err != nil {
			return nil, fmt.Errorf("could not unpack commit message: %w", err)
		}
		return commit, nil

	case MessageViewChange:
		var viewChange ViewChange
		err = json.Unmarshal(payload, &viewChange)
		if err != nil {
			return nil, fmt.Errorf("could not unpack view change message: %w", err)
		}
		return viewChange, nil

	case MessageNewView:
		var newView NewView
		err = json.Unmarshal(payload, &newView)
		if err != nil {
			return nil, fmt.Errorf("could not unpack new view message: %w", err)
		}
		return newView, nil
	}

	return nil, fmt.Errorf("unexpected message type (type: %v)", msg.Type)
}
