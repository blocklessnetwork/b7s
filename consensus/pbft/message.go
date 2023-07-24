package pbft

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/libp2p/go-libp2p/core/peer"
)

type MessageType uint

const (
	MessageRequest MessageType = iota + 1
	MessagePrePrepare
	MessagePrepare
	MessageCommit
)

func (m MessageType) String() string {
	switch m {
	case MessagePrePrepare:
		return "MessagePrePrepare"
	case MessagePrepare:
		return "MessagePrepare"
	case MessageCommit:
		return "MessageCommit"
	default:
		return fmt.Sprintf("unknown: %d", m)
	}
}

type Request struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Origin    peer.ID         `json:"origin"`
	Execute   execute.Request `json:"execute"`
}

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

// TODO: In fabric code, all messages have a `replicaID` field.

type PrePrepare struct {
	View           uint    `json:"view"`
	SequenceNumber uint    `json:"sequence_number"`
	Digest         string  `json:"digest"`
	Request        Request `json:"request"`
}

func (p PrePrepare) MarshalJSON() ([]byte, error) {

	// Define an alias without the JSON marshaller.
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

type Prepare struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`
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

type Commit struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`
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

// messageRecord is used as an interim format to supplement the original type with its type.
// Useful for serialization to automatically include the message type field.
type messageRecord struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func unpackMessage(payload []byte) (any, error) {

	var msg messageRecord
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return nil, fmt.Errorf("could not unpack base message: %w", err)
	}

	switch msg.Type {
	case MessageRequest:
		var request Request
		err = json.Unmarshal(msg.Data, &request)
		if err != nil {
			return nil, fmt.Errorf("could not unpack request: %w", err)
		}
		return request, nil

	case MessagePrePrepare:
		var preprepare PrePrepare
		err = json.Unmarshal(msg.Data, &preprepare)
		if err != nil {
			return nil, fmt.Errorf("could not unpack pre-prepare message: %w", err)
		}
		return preprepare, nil

	case MessagePrepare:
		var prepare Prepare
		err = json.Unmarshal(msg.Data, &prepare)
		if err != nil {
			return nil, fmt.Errorf("could not unpack prepare message: %w", err)
		}
		return prepare, nil

	case MessageCommit:
		var commit Commit
		err = json.Unmarshal(msg.Data, &commit)
		if err != nil {
			return nil, fmt.Errorf("could not unpack commit message: %w", err)
		}
		return commit, nil
	}

	return nil, fmt.Errorf("unexpected message type (type: %v)", msg.Type)
}
