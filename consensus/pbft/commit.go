package pbft

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// TODO (pbft): Split into `shouldSendCommit` and `sendCommit`.
func (r *Replica) maybeSendCommit(view uint, sequenceNo uint, digest string) error {

	log := r.log.With().Uint("view", view).Uint("sequence_number", sequenceNo).Str("digest", digest).Logger()

	if !r.prepared(view, sequenceNo, digest) {
		log.Info().Msg("request not yet prepared, not committing")
		return nil
	}

	// Have we already sent a commit message?
	msgID := getMessageID(view, sequenceNo)
	commits, ok := r.commits[msgID]
	if ok {
		_, sent := commits.m[r.id]
		if sent {
			log.Info().Msg("already have broadcast commit for this request, stopping now")
			return nil
		}
	}

	log.Info().Msg("request prepared, broadcasting commit")

	err := r.sendCommit(view, sequenceNo, digest)
	if err != nil {
		return fmt.Errorf("could not send commit message: %w", err)
	}

	log.Info().Msg("commit successfuly broadcast")

	// TODO (pbft): This function does too much, split.
	if !r.committed(view, sequenceNo, digest) {
		log.Info().Msg("request is not yet committed")
		return nil
	}

	log.Info().Msg("request committed, executing")

	return r.execute(view, sequenceNo, digest)
}

func (r *Replica) sendCommit(view uint, sequenceNo uint, digest string) error {

	log := r.log.With().Uint("view", view).Uint("sequence_number", sequenceNo).Str("digest", digest).Logger()

	log.Info().Msg("broadcasting commit message")

	commit := Commit{
		View:           view,
		SequenceNumber: sequenceNo,
		Digest:         digest,
	}

	err := r.broadcast(commit)
	if err != nil {
		return fmt.Errorf("could not broadcast commit message: %w", err)
	}

	log.Info().Msg("commit message successfully broadcast")

	// Record this commit message.
	r.recordCommitReceipt(r.id, commit)

	return nil
}

func (r *Replica) processCommit(replica peer.ID, commit Commit) error {

	log := r.log.With().Str("replica", replica.String()).Uint("view", commit.View).Uint("sequence_no", commit.SequenceNumber).Str("digest", commit.Digest).Logger()

	log.Info().Msg("received commit message")

	if commit.View != r.view {
		return fmt.Errorf("commit has an invalid view value (received: %v, current: %v)", commit.View, r.view)
	}

	r.recordCommitReceipt(replica, commit)

	if !r.committed(commit.View, commit.SequenceNumber, commit.Digest) {
		log.Info().Msg("request is not yet committed")
		return nil
	}

	err := r.execute(commit.View, commit.SequenceNumber, commit.Digest)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)

	}

	return nil
}

func (r *Replica) recordCommitReceipt(replica peer.ID, commit Commit) {

	msgID := getMessageID(commit.View, commit.SequenceNumber)
	commits, ok := r.commits[msgID]
	if !ok {
		r.commits[msgID] = newCommitReceipts()
		commits = r.commits[msgID]
	}

	commits.Lock()
	defer commits.Unlock()

	// Have we already seen this commit?
	_, exists := commits.m[replica]
	if exists {
		r.log.Warn().Uint("view", commit.View).Uint("sequence", commit.SequenceNumber).Str("digest", commit.Digest).Msg("ignoring duplicate commit")
		return
	}

	commits.m[replica] = commit
}
