package node

import (
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// Function provides retrieval of function manifest.
type Function interface {
	// Get retrieves a function manifest based on the address or CID. `useCached` boolean
	// determines if function manifest should be refetched or previously cached data can be returned.
	Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error)
}
