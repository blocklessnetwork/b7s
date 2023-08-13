package pbft

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// NOTE: JSON encoding related code is in serialization.go

type MessageType uint

const (
	MessageRequest MessageType = iota + 1
	MessagePrePrepare
	MessagePrepare
	MessageCommit
	MessageViewChange
	MessageNewView
)

func (m MessageType) String() string {
	switch m {
	case MessagePrePrepare:
		return "MessagePrePrepare"
	case MessagePrepare:
		return "MessagePrepare"
	case MessageCommit:
		return "MessageCommit"
	case MessageViewChange:
		return "MessageViewChange"
	case MessageNewView:
		return "MessageNewView"
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

type PrePrepare struct {
	View           uint    `json:"view"`
	SequenceNumber uint    `json:"sequence_number"`
	Digest         string  `json:"digest"`
	Request        Request `json:"request"`
}

type Prepare struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`
}

type Commit struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`
}

type ViewChange struct {
	View     uint          `json:"view"`
	Prepares []PrepareInfo `json:"prepares"`

	// Technically, view change message also includes:
	//	- n - sequence number of the last stable checkpoint => not needed here since we don't support checkpoints
	//  - C - 2f+1 checkpoint messages proving the correctness of s => see above
	//	- P - set Pm for each request m prepared at replica i with a sequence number higher than n; Pm includes a valid pre-prepare message and 2f matching, valid prepared messages (same view, sequence number and digest of m). Because we don't support checkpoints, this means everything from sequence number 0.
}

type PrepareInfo struct {
	View           uint                `json:"view"`
	SequenceNumber uint                `json:"sequence_number"`
	Digest         string              `json:"digest"`
	PrePrepare     PrePrepare          `json:"preprepare"`
	Prepares       map[peer.ID]Prepare `json:"prepares"`
}
type NewView struct {
	View        uint                   `json:"view"`
	Messages    map[peer.ID]ViewChange `json:"messages"`
	PrePrepares []PrePrepare           `json:"preprepares"`
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

	case MessageViewChange:
		var viewChange ViewChange
		err = json.Unmarshal(msg.Data, &viewChange)
		if err != nil {
			return nil, fmt.Errorf("could not unpack view change message: %w", err)
		}
		return viewChange, nil

	case MessageNewView:
		var newView NewView
		err = json.Unmarshal(msg.Data, &newView)
		if err != nil {
			return nil, fmt.Errorf("could not unpack new view message: %w", err)
		}
		return newView, nil
	}

	return nil, fmt.Errorf("unexpected message type (type: %v)", msg.Type)
}
