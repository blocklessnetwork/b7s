package pbft

func (r *Replica) execute(digest string) error {

	// Remove this request from the list of outstanding requests.
	delete(r.pending, digest)

	// TODO (pbft): Implement.
	r.log.Warn().Str("digest", digest).Msg("executing the request")

	return nil
}
