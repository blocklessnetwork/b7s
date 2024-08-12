package host

const (
	// Sentinel error for DHT.
	errNoGoodAddresses = "no good addresses"
)

var (
	messagesSentMetric      = []string{"messages", "sent"}
	messagesPublishedMetric = []string{"messages", "published"}
)
