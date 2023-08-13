package pbft

import (
	"encoding/json"
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

func (p PrePrepare) MarshalJSON() ([]byte, error) {

	// Define aliases without the JSON marshaller.
	type arequest Request
	type alias struct {
		View           uint     `json:"view"`
		SequenceNumber uint     `json:"sequence_number"`
		Digest         string   `json:"digest"`
		Request        arequest `json:"request"`
	}

	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data alias       `json:"data"`
		}{
			Type: MessagePrePrepare,
			Data: alias{
				View:           p.View,
				SequenceNumber: p.SequenceNumber,
				Digest:         p.Digest,
				Request:        arequest(p.Request),
			},
		})
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

func (p PrepareInfo) MarshalJSON() ([]byte, error) {

	encodedPrepareMap := make(map[string]Prepare)
	for id, prepare := range p.Prepares {
		encodedPrepareMap[id.String()] = prepare
	}

	type aliasPrePrepare PrePrepare
	return json.Marshal(
		struct {
			View           uint               `json:"view"`
			SequenceNumber uint               `json:"sequence_number"`
			Digest         string             `json:"digest"`
			PrePrepare     aliasPrePrepare    `json:"preprepare"`
			Prepares       map[string]Prepare `json:"prepares"`
		}{
			View:           p.View,
			SequenceNumber: p.SequenceNumber,
			Digest:         p.Digest,
			PrePrepare:     aliasPrePrepare(p.PrePrepare),
			Prepares:       encodedPrepareMap,
		},
	)
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

func (v NewView) MarshalJSON() ([]byte, error) {
	// To properly handle `peer.ID` serialization, this is a bit more involved.
	// See documentation for `ResultMap.MarshalJSON` in `models/execute/response.go`.
	type preprepareAlias PrePrepare
	type newView struct {
		View        uint                  `json:"view"`
		Messages    map[string]ViewChange `json:"messages"`
		PrePrepares []preprepareAlias     `json:"preprepares"`
	}

	nv := make(map[string]ViewChange)
	for replica, vc := range v.Messages {
		nv[replica.String()] = vc
	}

	preprepares := make([]preprepareAlias, 0, len(v.PrePrepares))
	for _, pp := range v.PrePrepares {
		preprepares = append(preprepares, preprepareAlias(pp))
	}

	return json.Marshal(
		struct {
			Type MessageType `json:"type"`
			Data newView     `json:"data"`
		}{
			Type: MessageNewView,
			Data: newView{
				View:        v.View,
				Messages:    nv,
				PrePrepares: preprepares,
			},
		})
}
