package node

import (
	"fmt"
	"io"

	"github.com/hashicorp/raft"
)

type NotImplementedFSM struct{}

func (n NotImplementedFSM) Apply(*raft.Log) interface{} {
	return nil
}

func (n NotImplementedFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, fmt.Errorf("TBD: not implemented")
}

func (n NotImplementedFSM) Restore(snapshot io.ReadCloser) error {
	return fmt.Errorf("TBD: not implemented")
}
