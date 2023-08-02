package pbft

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
)

// execute executes the request AND sends the result back to origin.
func (r *Replica) execute(digest string) error {

	defer r.stopRequestTimer()
	// TODO (pbft): Start new timer after this, as we're now waiting for a new request to execute - IF there are any pending requests.

	// Sanity check, should not happen.
	request, ok := r.requests[digest]
	if !ok {
		return fmt.Errorf("unknown request (digest: %s)", digest)
	}

	log := r.log.With().Str("digest", digest).Str("request", request.ID).Logger()

	// We don't want to execute a job multiple times.
	_, havePending := r.pending[digest]
	if !havePending {
		r.log.Warn().Str("digest", digest).Str("request", request.ID).Msg("no pending request with matching info - likely already executed")
		return nil
	}
	// Remove this request from the list of outstanding requests.
	delete(r.pending, digest)

	log.Info().Msg("executing request")

	res, err := r.executor.ExecuteFunction(request.ID, request.Execute)
	if err != nil {
		log.Error().Err(err).Msg("execution failed")
	}

	log.Info().Msg("executed request")

	msg := response.Execute{
		Type:      blockless.MessageExecuteResponse,
		Code:      res.Code,
		RequestID: request.ID,
		Results: execute.ResultMap{
			r.id: res,
		},
	}

	err = r.send(request.Origin, msg, blockless.ProtocolID)
	if err != nil {
		return fmt.Errorf("could not send execution response to node (target: %s, request: %s): %w", request.Origin.String(), request.ID, err)
	}

	return nil
}
