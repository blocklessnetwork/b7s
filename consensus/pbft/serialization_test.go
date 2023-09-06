package pbft

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/testing/mocks"
)

var (
	genericRequest = Request{
		ID:        mocks.GenericUUID.String(),
		Timestamp: time.Now().UTC(),
		Origin:    mocks.GenericPeerID,
		Execute:   mocks.GenericExecutionRequest,
	}

	genericPrepareInfo = PrepareInfo{
		View:           1,
		SequenceNumber: 14,
		Digest:         "abcdef123456789",
		PrePrepare: PrePrepare{
			View:           1,
			SequenceNumber: 14,
			Digest:         "abcdef123456789",
			Request: Request{
				ID:        mocks.GenericUUID.String(),
				Timestamp: time.Now().UTC(),
				Origin:    mocks.GenericPeerID,
				Execute:   mocks.GenericExecutionRequest,
			},
		},
		// These are all different but it's test data so it's fine.
		Prepares: map[peer.ID]Prepare{
			peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x91}): {
				View:           45,
				SequenceNumber: 19,
				Digest:         "abcdef123456789",
			},
			peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x92}): {
				View:           98,
				SequenceNumber: 12,
				Digest:         "987654321",
			},
			peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x93}): {
				View:           100,
				SequenceNumber: 91,
				Digest:         "abc123def456",
			},
		},
	}
)

func TestRequest_Serialization(t *testing.T) {

	orig := genericRequest

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked Request
	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)

	require.Equal(t, orig, unpacked)
}

func TestPrePrepare_Serialization(t *testing.T) {

	orig := PrePrepare{
		View:           1,
		SequenceNumber: 2,
		Digest:         "123456789",
		Request:        genericRequest,
	}

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked PrePrepare
	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)

	require.Equal(t, orig, unpacked)
}

func TestPrepare_Serialization(t *testing.T) {

	orig := Prepare{
		View:           14,
		SequenceNumber: 45,
		Digest:         "abc123def",
	}

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked Prepare
	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)
	require.Equal(t, orig, unpacked)
}

func TestCommit_Serialization(t *testing.T) {

	orig := Commit{
		View:           23,
		SequenceNumber: 51,
		Digest:         "987xyz",
	}

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked Commit
	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)
	require.Equal(t, orig, unpacked)
}

func TestPrepareInfo_Serialization(t *testing.T) {

	orig := genericPrepareInfo

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked PrepareInfo
	err = json.Unmarshal(encoded, &unpacked)

	require.NoError(t, err)
	require.Equal(t, orig, unpacked)
}

func TestViewChange_Serialization(t *testing.T) {

	// Use the generic type as the baseline but change individual fields.
	info1 := genericPrepareInfo
	info1.View = 973

	info2 := genericPrepareInfo
	info2.SequenceNumber = 175

	info3 := genericPrepareInfo
	info3.PrePrepare.Digest = "xyz"

	orig := ViewChange{
		View:     15,
		Prepares: []PrepareInfo{info1, info2, info3},
	}

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked ViewChange

	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)
	require.Equal(t, orig, unpacked)
}

func TestNewView_Serialization(t *testing.T) {

	var (
		genericPeerID1 = peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x91})
		genericPeerID2 = peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x92})
		genericPeerID3 = peer.ID([]byte{0x0, 0x24, 0x8, 0x1, 0x12, 0x20, 0x56, 0x77, 0x86, 0x82, 0x76, 0xa, 0xc5, 0x9, 0x63, 0xde, 0xe4, 0x31, 0xfc, 0x44, 0x75, 0xdd, 0x5a, 0x27, 0xee, 0x6b, 0x94, 0x13, 0xed, 0xe2, 0xa3, 0x6d, 0x8a, 0x1d, 0x57, 0xb6, 0xb8, 0x93})

		orig = NewView{
			View: 654,
			PrePrepares: []PrePrepare{
				{
					View:           12,
					SequenceNumber: 32,
					Request:        genericRequest,
				},
				{
					View:           32,
					SequenceNumber: 45,
					Request:        genericRequest,
				},
				{
					View:           78,
					SequenceNumber: 32,
					Request:        genericRequest,
				},
			},
			Messages: map[peer.ID]ViewChange{
				genericPeerID1: {
					View: 13,
					Prepares: []PrepareInfo{
						genericPrepareInfo,
					},
				},
				genericPeerID2: {
					View: 13,
					Prepares: []PrepareInfo{
						genericPrepareInfo,
					},
				},
				genericPeerID3: {
					View: 13,
					Prepares: []PrepareInfo{
						genericPrepareInfo,
					},
				},
			},
		}
	)

	encoded, err := json.Marshal(orig)
	require.NoError(t, err)

	var unpacked NewView
	err = json.Unmarshal(encoded, &unpacked)
	require.NoError(t, err)
	require.Equal(t, orig, unpacked)
}
