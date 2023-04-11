package node

import (
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// FStore provides retrieval of function manifest.
type FStore interface {
	// Get retrieves a function manifest based on the address or CID. `useCached` boolean
	// determines if function manifest should be refetched or previously cached data can be returned.
	Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error)

	// Installed returns info if the function is installed or not.
	Installed(cid string) (bool, error)

	// InstalledFunction returns the list of CIDs of installed functions.
	InstalledFunctions() []string

	// Sync will recheck if function installation is found in local storage, and redownload it if it isn't.
	Sync(cid string) error
}
