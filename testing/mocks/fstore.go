package mocks

import (
	"testing"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type FStore struct {
	InstallFunc            func(string, string) error
	GetFunc                func(string) (*blockless.FunctionManifest, error)
	InstalledFunc          func(string) (bool, error)
	InstalledFunctionsFunc func() ([]string, error)
	SyncFunc               func(bool) error
}

func BaselineFStore(t *testing.T) *FStore {
	t.Helper()

	fh := FStore{
		GetFunc: func(string) (*blockless.FunctionManifest, error) {
			return &GenericManifest, nil
		},
		InstallFunc: func(string, string) error {
			return nil
		},
		InstalledFunc: func(string) (bool, error) {
			return true, nil
		},
		InstalledFunctionsFunc: func() ([]string, error) {
			return nil, nil
		},
		SyncFunc: func(bool) error {
			return nil
		},
	}

	return &fh
}

func (f *FStore) Install(address string, cid string) error {
	return f.InstallFunc(address, cid)
}

func (f *FStore) Get(cid string) (*blockless.FunctionManifest, error) {
	return f.GetFunc(cid)
}

func (f *FStore) Installed(cid string) (bool, error) {
	return f.InstalledFunc(cid)
}

func (f *FStore) InstalledFunctions() ([]string, error) {
	return f.InstalledFunctionsFunc()
}

func (f *FStore) Sync(haltOnError bool) error {
	return f.SyncFunc(haltOnError)
}
