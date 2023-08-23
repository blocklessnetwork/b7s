package pbft

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
)

// execute executes the request AND sends the result back to origin.
func (r *Replica) execute(view uint, sequence uint, digest string) error {

	// Sanity check, should not happen.
	request, ok := r.requests[digest]
	if !ok {
		return fmt.Errorf("unknown request (digest: %s)", digest)
	}

	log := r.log.With().Uint("view", view).Uint("sequence", sequence).Str("digest", digest).Str("request", request.ID).Logger()

	// We don't want to execute a job multiple times.
	_, havePending := r.pending[digest]
	if !havePending {
		log.Warn().Msg("no pending request with matching info - likely already executed")
		return nil
	}

	// Requests must be executed in order.
	if sequence != r.lastExecuted+1 {
		log.Warn().Msg("requests with lower sequence number have not been executed")
		// TODO (pbft): Start execution of earlier requests?
		return nil
	}

	// Sanity check - should never happen.
	if sequence < r.lastExecuted {
		log.Error().Uint("last_executed", r.lastExecuted).Msg("requests executed out of order!")
	}

	// Remove this request from the list of outstanding requests.
	delete(r.pending, digest)

	log.Info().Msg("executing request")

	res, err := r.executor.ExecuteFunction(request.ID, request.Execute)
	if err != nil {
		log.Error().Err(err).Msg("execution failed")
	}

	// Stop the timer since we completed an execution.
	r.stopRequestTimer()

	// If we have more pending requests, start a new timer.
	if len(r.pending) > 0 {
		r.startRequestTimer(true)
	}

	log.Info().Msg("executed request")

	r.lastExecuted = sequence

	msg := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      res.Code,
		RequestID: request.ID,
		Results: execute.ResultMap{
			r.id: res,
		},
	}

	// Save this executions in case it's requested again.
	r.executions[request.ID] = msg

	err = r.send(request.Origin, msg, blockless.ProtocolID)
	if err != nil {
		return fmt.Errorf("could not send execution response to node (target: %s, request: %s): %w", request.Origin.String(), request.ID, err)
	}

	return nil
}