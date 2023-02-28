package node

import (
	"sync"

	"github.com/blocklessnetworking/b7s/models/response"
)

type rollCallQueue struct {
	sync.Mutex

	size uint
	m    map[string]chan response.RollCall
}

func newQueue(bufSize uint) *rollCallQueue {

	q := rollCallQueue{
		size: bufSize,
		m:    make(map[string]chan response.RollCall),
	}

	return &q
}

// add records a new response to a roll call.
func (q *rollCallQueue) add(id string, res response.RollCall) {
	q.Lock()
	defer q.Unlock()

	_, ok := q.m[id]
	if !ok {
		q.m[id] = make(chan response.RollCall, q.size)
	}

	q.m[id] <- res
}

// responses will return a channel that can be used to iterate through all of the responses.
func (q *rollCallQueue) responses(id string) <-chan response.RollCall {
	q.Lock()

	_, ok := q.m[id]
	if !ok {
		q.m[id] = make(chan response.RollCall, q.size)
	}

	q.Unlock()

	return q.m[id]
}
