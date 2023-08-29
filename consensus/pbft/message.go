package pbft

import (
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

	// Signed digest of the pre-prepare message.
	Signature string `json:"signature,omitempty"`
}

type Prepare struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`

	// Signed digest of the prepare message.
	Signature string `json:"signature,omitempty"`
}

type Commit struct {
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`

	// Signed digest of the commit message.
	Signature string `json:"signature,omitempty"`
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
