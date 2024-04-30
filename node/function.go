package node

// FStore provides retrieval of function manifest.
type FStore interface {
	// Install will install a function based on the address and CID.
	Install(address string, cid string) error

	// Installed returns info if the function is installed or not.
	Installed(cid string) (bool, error)

	// Sync will ensure function installations are correct, redownloading functions if needed.
	Sync(haltOnError bool) error
}
