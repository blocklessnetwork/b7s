package node

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/models/execute"
)

type fsmLogEntry struct {
	RequestID string          `json:"request_id,omitempty"`
	Origin    peer.ID         `json:"origin,omitempty"`
	Execute   execute.Request `json:"execute,omitempty"`
}

type fsmProcessFunc func(req fsmLogEntry, res execute.Result)

type fsmExecutor struct {
	log        zerolog.Logger
	executor   Executor
	processors []fsmProcessFunc
}

func newFsmExecutor(log zerolog.Logger, executor Executor, processors ...fsmProcessFunc) *fsmExecutor {

	ps := make([]fsmProcessFunc, 0, len(processors))
	ps = append(ps, processors...)

	fsm := fsmExecutor{
		log:        log.With().Str("component", "fsm").Logger(),
		executor:   executor,
		processors: ps,
	}

	return &fsm
}

func (f fsmExecutor) Apply(log *raft.Log) interface{} {

	f.log.Info().Msg("applying log entry")

	// Unpack the execution request.
	payload := log.Data

	var logEntry fsmLogEntry
	err := json.Unmarshal(payload, &logEntry)
	if err != nil {
		return fmt.Errorf("could not unmarshal request: %w", err)
	}

	f.log.Info().Str("request", logEntry.RequestID).Str("function", logEntry.Execute.FunctionID).Msg("FSM executing function")

	res, err := f.executor.ExecuteFunction(logEntry.RequestID, logEntry.Execute)
	if err != nil {
		return fmt.Errorf("could not execute function: %w", err)
	}

	// Execute processors.
	for _, proc := range f.processors {
		proc(logEntry, res)
	}

	f.log.Info().Str("request", logEntry.RequestID).Msg("FSM successfully executed function")

	return res
}

func (f fsmExecutor) Snapshot() (raft.FSMSnapshot, error) {
	f.log.Info().Msg("received snapshot request")
	return nil, fmt.Errorf("TBD: not implemented")
}

func (f fsmExecutor) Restore(snapshot io.ReadCloser) error {
	f.log.Info().Msg("received snapshot restore request")
	return fmt.Errorf("TBD: not implemented")
}
