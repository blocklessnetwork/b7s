package pbft

type pbftCore struct {
	// Number of replicas in the cluster.
	n uint

	// Number of byzantine replicas we can tolerate.
	f uint

	// Sequence number.
	sequence uint

	// ViewNumber.
	view uint
}

func newPbftCore(total uint) pbftCore {

	return pbftCore{
		sequence: 0,
		view:     0,
		n:        total,
		f:        calcByzantineTolerance(total),
	}
}

// given a view number, return the index of the expected primary.
func (c pbftCore) primary(v uint) uint {
	return v % c.n
}

// return the index of the expected primary for the current view.
func (c pbftCore) currentPrimary() uint {
	return c.view % c.n
}

func (c pbftCore) prepareQuorum() uint {
	// TODO (pbft): Not sure about this quorum - is it 2f total or 2f prepares from OTHER replicas?
	// We need 2f+1 for a quorum; but, the primary's `preprepare` message counts as it's vote, so we may need 2f more.
	return 2 * c.f
}

func (c pbftCore) commitQuorum() uint {
	return 2*c.f + 1
}

// based on the number of replicas, determine how many byzantine replicas we can tolerate.
func calcByzantineTolerance(n uint) uint {

	if n <= 1 {
		return 0
	}

	f := (n - 1) / 3
	return f
}

// messageID is used to identify a specific point in time as view + sequence number combination.
type messageID struct {
	view     uint
	sequence uint
}

func getMessageID(view uint, sequenceNo uint) messageID {
	return messageID{
		view:     view,
		sequence: sequenceNo,
	}
}
