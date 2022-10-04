package repository

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/src/models"
)

func TestGetPackage(t *testing.T) {
	ctx := context.Background()

	config := models.Config{}
	config.Node.WorkSpaceRoot = "/tmp/b7s_test"
	ctx = context.WithValue(ctx, "config", config)

	// file uri reference manifest
	p := GetPackage(ctx, "https://bafybeiho3scwi3njueloobzhg7ndn7yjb5rkcaydvsoxmnhmu2adv6oxzq.ipfs.w3s.link/manifest.json")
	t.Log(p)
}
