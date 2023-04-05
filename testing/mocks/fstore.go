package mocks

import (
	"testing"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type FStore struct {
	GetFunc func(string, string, bool) (*blockless.FunctionManifest, error)
}

func BaselineFunctionHandler(t *testing.T) *FStore {
	t.Helper()

	fh := FStore{
		GetFunc: func(string, string, bool) (*blockless.FunctionManifest, error) {
			return &GenericManifest, nil
		},
	}

	return &fh
}

func (f *FStore) Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error) {
	return f.GetFunc(address, cid, useCached)
}
