package node

import (
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// FunctionStore provides retrieval of function manifest.
// TODO: Somewhat of a name discrepancy - function store vs function handler.. settle down on a name.
type FunctionStore interface {
	// Get retrieves a function manifest based on the address or CID. `useCached` boolean
	// determines if function manifest should be refetched or previously cached data can be returned.
	Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error)
}
