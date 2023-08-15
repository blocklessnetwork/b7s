package pbft

import (
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (r Request) MarshalJSON() ([]byte, error) {

	// Define an alias without the JSON marshaller.
	type alias Request
	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessageRequest,
			Data: alias(r),
		})
}

func (r *Request) UnmarshalJSON(data []byte) error {

	request := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &request)
	if err != nil {
		return err
	}

	type alias *Request
	return json.Unmarshal(request.Data, alias(r))
}

func (p PrePrepare) MarshalJSON() ([]byte, error) {

	// Define aliases without the JSON marshaller.
	type alias PrePrepare
	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessagePrePrepare,
			Data: alias(p),
		})
}

func (p *PrePrepare) UnmarshalJSON(data []byte) error {

	prePrepare := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &prePrepare)
	if err != nil {
		return err
	}

	type alias *PrePrepare
	return json.Unmarshal(prePrepare.Data, alias(p))
}

func (p Prepare) MarshalJSON() ([]byte, error) {
	type alias Prepare
	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessagePrepare,
			Data: alias(p),
		})
}

func (p *Prepare) UnmarshalJSON(data []byte) error {

	prepare := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &prepare)
	if err != nil {
		return err
	}

	// Aliased to prevent recursive call to UnmarshalJSON.
	type aliased *Prepare
	return json.Unmarshal(prepare.Data, aliased(p))
}

func (c Commit) MarshalJSON() ([]byte, error) {
	type alias Commit
	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessageCommit,
			Data: alias(c),
		})
}

func (c *Commit) UnmarshalJSON(data []byte) error {

	commit := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &commit)
	if err != nil {
		return err
	}

	// Aliased to prevent recursive call to UnmarshalJSON.
	type aliased *Commit
	return json.Unmarshal(commit.Data, aliased(c))
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
	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessageViewChange,
			Data: alias(v),
		})
}

func (v *ViewChange) UnmarshalJSON(data []byte) error {

	vc := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &vc)
	if err != nil {
		return err
	}

	// Aliased to prevent recursive call to UnmarshalJSON.
	type aliased *ViewChange
	return json.Unmarshal(vc.Data, aliased(v))
}

type newViewEncode struct {
	View        uint                  `json:"view"`
	Messages    map[string]ViewChange `json:"messages"`
	PrePrepares []PrePrepare          `json:"preprepares"`
}

func (v NewView) MarshalJSON() ([]byte, error) {
	// To properly handle `peer.ID` serialization, this is a bit more involved.
	// See documentation for `ResultMap.MarshalJSON` in `models/execute/response.go`.

	nv := make(map[string]ViewChange)
	for replica, vc := range v.Messages {
		nv[replica.String()] = vc
	}

	return json.Marshal(
		struct {
			Type MessageType   `json:"type"`
			Data newViewEncode `json:"data"`
		}{
			Type: MessageNewView,
			Data: newViewEncode{
				View:        v.View,
				Messages:    nv,
				PrePrepares: v.PrePrepares,
			},
		})
}

func (n *NewView) UnmarshalJSON(data []byte) error {

	newView := struct {
		Type MessageType     `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	err := json.Unmarshal(data, &newView)
	if err != nil {
		return err
	}

	var nv newViewEncode
	err = json.Unmarshal(newView.Data, &nv)
	if err != nil {
		return err
	}

	viewChangeMap := make(map[peer.ID]ViewChange)
	for idStr, vc := range nv.Messages {
		id, err := peer.Decode(idStr)
		if err != nil {
			return fmt.Errorf("could not decode peer.ID (str: %s): %w", idStr, err)
		}
		viewChangeMap[id] = vc
	}

	*n = NewView{
		View:        nv.View,
		Messages:    viewChangeMap,
		PrePrepares: nv.PrePrepares,
	}

	return nil
}
