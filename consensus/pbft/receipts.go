package pbft

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
)

// prepareReceipts maps a peer/replica ID to the prepare message it sent.
type prepareReceipts struct {
	m map[peer.ID]Prepare
	*sync.Mutex
}

func newPrepareReceipts() *prepareReceipts {

	pr := prepareReceipts{
		m:     make(map[peer.ID]Prepare),
		Mutex: &sync.Mutex{},
	}

	return &pr
}

type commitReceipts struct {
	m map[peer.ID]Commit
	*sync.Mutex
}

func newCommitReceipts() *commitReceipts {

	cr := commitReceipts{
		m:     make(map[peer.ID]Commit),
		Mutex: &sync.Mutex{},
	}

	return &cr
}

type viewChangeReceipts struct {
	m map[peer.ID]ViewChange
	*sync.Mutex
}

func newViewChangeReceipts() *viewChangeReceipts {

	vcr := viewChangeReceipts{
		m:     make(map[peer.ID]ViewChange),
		Mutex: &sync.Mutex{},
	}

	return &vcr
}
