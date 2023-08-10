package pbft

import (
	"sync"
)

type replicaState struct {
	// State lock. This is a global lock for all state maps (pre-prepares, prepares, commits etc).
	// Even though access to those could be managed on a more granular level, at the moment I'm not
	// sure it's worth it.
	sl *sync.Mutex

	// False if view change is in progress.
	activeView bool

	// Sequence number of last execution.
	lastExecuted uint

	// Keep track of seen requests. Map request to the digest.
	requests map[string]Request
	// Keep track of requests queued for execution. Could also be tracked via a single map.
	pending map[string]Request

	// Keep track of seen pre-prepare messages.
	preprepares map[messageID]PrePrepare
	// Keep track of seen prepare messages.
	prepares map[messageID]*prepareReceipts
	// Keep track of seen commit messages.
	commits map[messageID]*commitReceipts
	// Keep track of view change messages.
	viewChanges map[uint]*viewChangeReceipts
}

func newState() replicaState {

	state := replicaState{
		sl: &sync.Mutex{},

		activeView:  true,
		requests:    make(map[string]Request),
		pending:     make(map[string]Request),
		preprepares: make(map[messageID]PrePrepare),
		prepares:    make(map[messageID]*prepareReceipts),
		commits:     make(map[messageID]*commitReceipts),
		viewChanges: make(map[uint]*viewChangeReceipts),
	}

	return state
}
