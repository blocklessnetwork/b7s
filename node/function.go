package node

// FStore provides retrieval of function manifest.
type FStore interface {
	// Install will install a function based on the address and CID.
	Install(address string, cid string) error

	// Installed returns info if the function is installed or not.
	Installed(cid string) (bool, error)

	// InstalledFunction returns the list of CIDs of installed functions.
	InstalledFunctions() ([]string, error)

	// Sync will recheck if function installation is found in local storage, and redownload it if it isn't.
	Sync(cid string) error
}
