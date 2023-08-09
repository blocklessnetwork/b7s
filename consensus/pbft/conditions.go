package pbft

func (r *Replica) prePrepared(view uint, sequenceNo uint, digest string) bool {

	// TODO (pbft): This is now okay as it's a null request.
	if digest == "" {
		return false
	}

	// Have we seen this request before?
	_, seen := r.requests[digest]
	if !seen {
		return false
	}

	// Do we have a pre-prepare for this request?
	preprepare, seen := r.preprepares[getMessageID(view, sequenceNo)]
	if !seen {
		return false
	}

	if preprepare.Digest != digest {
		return false
	}

	return true
}

func (r *Replica) prepared(view uint, sequenceNo uint, digest string) bool {

	// Check if we have seen this request before.
	// NOTE: This is also checked as part of the pre-prepare check.
	_, seen := r.requests[digest]
	if !seen {
		return false
	}

	// Is the pre-prepare condition met for this request?
	if !r.prePrepared(view, sequenceNo, digest) {
		return false
	}

	prepares, ok := r.prepares[getMessageID(view, sequenceNo)]
	if !ok {
		return false
	}

	prepareCount := uint(len(prepares.m))
	haveQuorum := prepareCount >= r.prepareQuorum()

	r.log.Debug().Str("digest", digest).Uint("view", view).Uint("sequence_no", sequenceNo).
		Uint("quorum", prepareCount).Bool("have_quorum", haveQuorum).
		Msg("number of prepares for a request")

	return haveQuorum
}

func (r *Replica) committed(view uint, sequenceNo uint, digest string) bool {

	// Is the prepare condition met for this request?
	if !r.prepared(view, sequenceNo, digest) {
		return false
	}

	commits, ok := r.commits[getMessageID(view, sequenceNo)]
	if !ok {
		return false
	}

	commitCount := uint(len(commits.m))
	haveQuorum := commitCount > r.commitQuorum()

	r.log.Debug().Str("digest", digest).Uint("view", view).Uint("sequence_no", sequenceNo).
		Uint("quorum", commitCount).Bool("have_quorum", haveQuorum).
		Msg("number of commits for a request")

	return haveQuorum
}

func (r *Replica) viewChangeReady(view uint) bool {

	vc, ok := r.viewChanges[view]
	if !ok {
		return false
	}

	vc.Lock()
	defer vc.Unlock()

	vcCount := uint(len(vc.m))
	haveQuorum := vcCount >= r.commitQuorum()

	r.log.Debug().Uint("view", view).Uint("quorum", vcCount).Bool("have_quorum", haveQuorum).Msg("number of view change messages for a view")

	return haveQuorum
}
