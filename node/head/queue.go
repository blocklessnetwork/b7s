package head

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/models/response"
)

type rollCallQueue struct {
	sync.Mutex

	size uint
	m    map[string]chan rollCallResponse
}

type rollCallResponse struct {
	From peer.ID
	response.RollCall
}

// newQueue is used to record per-request roll call responses.
func newQueue(bufSize uint) *rollCallQueue {

	q := rollCallQueue{
		size: bufSize,
		m:    make(map[string]chan rollCallResponse),
	}

	return &q
}

// create will create a response queue for the given requestID.
// Needs to be called before receiving/reading roll call responses.
func (q *rollCallQueue) create(reqID string) {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[reqID]
	if ok {
		return
	}

	q.m[reqID] = make(chan rollCallResponse, q.size)
}

// add records a new response to a roll call.
func (q *rollCallQueue) add(id string, res rollCallResponse) {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[id]
	if !ok {
		return
	}

	q.m[id] <- res
}

// exists returns true if a given request ID exists in the roll call map.
func (q *rollCallQueue) exists(reqID string) bool {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[reqID]
	return ok
}

// responses will return a channel that can be used to iterate through all of the responses.
func (q *rollCallQueue) responses(reqID string) <-chan rollCallResponse {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[reqID]
	if !ok {
		// Technically we shouldn't be here since we already called `create`, but there's also no harm in it.
		q.m[reqID] = make(chan rollCallResponse, q.size)
	}

	return q.m[reqID]
}

// Remove will remove the channel with the given ID.
func (q *rollCallQueue) remove(reqID string) {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[reqID]
	if !ok {
		// Should not be done but make it safe for double close.
		return
	}

	// First drain the channel.
	for len(q.m[reqID]) > 0 {
		<-q.m[reqID]
	}

	close(q.m[reqID])
	delete(q.m, reqID)
}
