package execute

import (
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
)

// Result describes an execution result.
type Result struct {
	Code   codes.Code    `json:"code"`
	Result RuntimeOutput `json:"result"`
	Usage  Usage         `json:"usage,omitempty"`
}

// Cluster represents the set of peers that executed the request.
type Cluster struct {
	Main  peer.ID   `json:"main,omitempty"`
	Peers []peer.ID `json:"peers,omitempty"`
}

// RuntimeOutput describes the output produced by the Blockless Runtime during execution.
type RuntimeOutput struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Log      string `json:"-"`
}

// Usage represents the resource usage information for a particular execution.
type Usage struct {
	WallClockTime time.Duration `json:"wall_clock_time,omitempty"`
	CPUUserTime   time.Duration `json:"cpu_user_time,omitempty"`
	CPUSysTime    time.Duration `json:"cpu_sys_time,omitempty"`
	MemoryMaxKB   int64         `json:"memory_max_kb,omitempty"`
}

// ResultMap contains execution results from multiple peers.
type ResultMap map[peer.ID]Result

// MarshalJSON provides means to correctly handle JSON serialization/deserialization.
// See:
//
//	https://github.com/libp2p/go-libp2p/pull/2156
//	https://github.com/libp2p/go-libp2p-resource-manager/pull/67#issuecomment-1176820561
func (m ResultMap) MarshalJSON() ([]byte, error) {

	em := make(map[string]Result, len(m))
	for p, v := range m {
		em[p.String()] = v
	}

	return json.Marshal(em)
}
