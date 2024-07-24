package host

const (
	// Sentinel error for DHT.
	errNoGoodAddresses = "no good addresses"

	defaultMustReachBootNodes = false

	// When we reach this number of connections, we'll prune open connections.
	connLimitHi = 1024
	// Number of connections we will leave after pruning.
	connLimitLo = 768
)
