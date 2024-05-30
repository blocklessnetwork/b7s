package node

import (
	"context"
)

// FStore provides retrieval of function manifest.
type FStore interface {
	// Install will install a function based on the address and CID.
	Install(ctx context.Context, address string, cid string) error

	// Installed returns info if the function is installed or not.
	Installed(cid string) (bool, error)

	// TODO: Refactor the sync code - move the logic outside of the package
	// Sync will ensure function installations are correct, redownloading functions if needed.
	Sync(ctx context.Context, haltOnError bool) error
}
