package mocks

import (
	"testing"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type FStore struct {
	GetFunc                func(string, string, bool) (*blockless.FunctionManifest, error)
	InstalledFunc          func(string) (bool, error)
	InstalledFunctionsFunc func() []string
	SyncFunc               func(string) error
}

func BaselineFunctionHandler(t *testing.T) *FStore {
	t.Helper()

	fh := FStore{
		GetFunc: func(string, string, bool) (*blockless.FunctionManifest, error) {
			return &GenericManifest, nil
		},
		InstalledFunc: func(string) (bool, error) {
			return true, nil
		},
		InstalledFunctionsFunc: func() []string {
			return nil
		},
		SyncFunc: func(string) error {
			return nil
		},
	}

	return &fh
}

func (f *FStore) Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error) {
	return f.GetFunc(address, cid, useCached)
}

func (f *FStore) Installed(cid string) (bool, error) {
	return f.InstalledFunc(cid)
}

func (f *FStore) InstalledFunctions() []string {
	return f.InstalledFunctionsFunc()
}

func (f *FStore) Sync(cid string) error {
	return f.SyncFunc(cid)
}
