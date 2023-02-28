package mocks

import (
	"testing"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type FunctionHandler struct {
	GetFunc func(string, string, bool) (*blockless.FunctionManifest, error)
}

func BaselineFunctionHandler(t *testing.T) *FunctionHandler {
	t.Helper()

	fh := FunctionHandler{
		GetFunc: func(string, string, bool) (*blockless.FunctionManifest, error) {
			return &GenericManifest, nil
		},
	}

	return &fh
}

func (f *FunctionHandler) Get(address string, cid string, useCached bool) (*blockless.FunctionManifest, error) {
	return f.GetFunc(address, cid, useCached)
}
