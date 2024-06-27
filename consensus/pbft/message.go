package pbft

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

type TraceableMessage interface {
	SaveTraceContext(t tracing.TraceInfo)
}

type BaseMessage struct {
	tracing.TraceInfo
}

func (m *BaseMessage) SaveTraceContext(t tracing.TraceInfo) {
	m.TraceInfo = t
}

// JSON encoding related code is in serialization.go
// Signature related code is in message_signature.go

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
	BaseMessage
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Origin    peer.ID         `json:"origin"`
	Execute   execute.Request `json:"execute"`
}

type PrePrepare struct {
	BaseMessage
	View           uint    `json:"view"`
	SequenceNumber uint    `json:"sequence_number"`
	Digest         string  `json:"digest"`
	Request        Request `json:"request"`

	// Signed digest of the pre-prepare message.
	Signature string `json:"signature,omitempty"`
}

type Prepare struct {
	BaseMessage
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`

	// Signed digest of the prepare message.
	Signature string `json:"signature,omitempty"`
}

type Commit struct {
	BaseMessage
	View           uint   `json:"view"`
	SequenceNumber uint   `json:"sequence_number"`
	Digest         string `json:"digest"`

	// Signed digest of the commit message.
	Signature string `json:"signature,omitempty"`
}

type ViewChange struct {
	BaseMessage
	View     uint          `json:"view"`
	Prepares []PrepareInfo `json:"prepares"`

	// Signed digest of the view change message.
	Signature string `json:"signature,omitempty"`

	// Technically, view change message also includes:
	//	- n - sequence number of the last stable checkpoint => not needed here since we don't support checkpoints
	//  - C - 2f+1 checkpoint messages proving the correctness of s => see above
}

type PrepareInfo struct {
	View           uint                `json:"view"`
	SequenceNumber uint                `json:"sequence_number"`
	Digest         string              `json:"digest"`
	PrePrepare     PrePrepare          `json:"preprepare"`
	Prepares       map[peer.ID]Prepare `json:"prepares"`
}
type NewView struct {
	BaseMessage
	View        uint                   `json:"view"`
	Messages    map[peer.ID]ViewChange `json:"messages"`
	PrePrepares []PrePrepare           `json:"preprepares"`

	// Signed digest of the new view message.
	Signature string `json:"signature,omitempty"`
}
