package mocks

import (
	"context"
	"testing"
)

type FStore struct {
	InstallFunc     func(context.Context, string, string) error
	IsInstalledFunc func(string) (bool, error)
	SyncFunc        func(context.Context, bool) error
}

func BaselineFStore(t *testing.T) *FStore {
	t.Helper()

	fh := FStore{
		InstallFunc: func(context.Context, string, string) error {
			return nil
		},
		IsInstalledFunc: func(string) (bool, error) {
			return true, nil
		},
		SyncFunc: func(context.Context, bool) error {
			return nil
		},
	}

	return &fh
}

func (f *FStore) Install(ctx context.Context, address string, cid string) error {
	return f.InstallFunc(ctx, address, cid)
}

func (f *FStore) IsInstalled(cid string) (bool, error) {
	return f.IsInstalledFunc(cid)
}

func (f *FStore) Sync(ctx context.Context, haltOnError bool) error {
	return f.SyncFunc(ctx, haltOnError)
}
