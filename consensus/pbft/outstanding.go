package pbft

// outstandingRequests returns the list of requests that have been seen by the replica, but are not already in the pipeline.
// This is called on view change to check if replicas should re-start their view change timers to keep the new primary honest.
// New primary should do the same thing to see if it has seen some requests that the previous replica has not made progress on and,
// if there are any, make actions related to these requests (by issuing preprepares).
func (r *Replica) outstandingRequests() []Request {

	r.log.Debug().Msg("checking if there are any requests not yet in the pipeline")

	var requests []Request

	for digest, request := range r.requests {

		log := r.log.With().Str("digest", digest).Str("request", request.ID).Logger()

		_, executed := r.executions[request.ID]
		if executed {
			log.Debug().Msg("request already executed, skipping")
			continue
		}

		log.Info().Msg("request not yet in the pipeline nor executed")

		// This means there's a request we've seen that hasn't been executed and not in the pipeline.
		requests = append(requests, request)
	}

	return requests
}
