package node

import (
	"fmt"
	"sync"
)

// Subgroups are (optional) groups of nodes that can work on specific things.
// Generally all nodes subscribe to the B7S general topic and can receive work from there.
// However, nodes can also be part of smaller groups, where they join a specific topic where
// some specific work (roll calls) may be published to.
type workSubgroups struct {
	*sync.RWMutex
	topics map[string]*topicInfo
}

// wrapper around topic joining + housekeeping.
func (n *Node) joinTopic(topic string) (*topicInfo, error) {

	n.subgroups.Lock()
	defer n.subgroups.Unlock()

	th, err := n.host.JoinTopic(topic)
	if err != nil {
		return nil, fmt.Errorf("could not join topic (topic: %s): %w", topic, err)
	}

	// NOTE: No subscription, joining topic only.
	ti := &topicInfo{
		handle: th,
	}

	n.subgroups.topics[topic] = ti

	return ti, nil
}
