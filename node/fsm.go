package node

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/rs/zerolog"

	"github.com/blocklessnetworking/b7s/models/execute"
)

type fsmLogEntry struct {
	RequestID string          `json:"request_id,omitempty"`
	Execute   execute.Request `json:"execute,omitempty"`
}

type fsmExecutor struct {
	log      zerolog.Logger
	executor Executor
}

func newFsmExecutor(log zerolog.Logger, executor Executor) *fsmExecutor {
	fsm := fsmExecutor{
		log:      log.With().Str("component", "fsm").Logger(),
		executor: executor,
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

	f.log.Info().Str("request_id", logEntry.RequestID).Str("function_id", logEntry.Execute.FunctionID).Msg("FSM executing function")

	res, err := f.executor.ExecuteFunction(logEntry.RequestID, logEntry.Execute)
	if err != nil {
		return fmt.Errorf("could not execute function: %w", err)
	}

	f.log.Info().Msg("log entry successfully applied")

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
