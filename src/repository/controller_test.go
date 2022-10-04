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

	p := GetPackage(ctx, "https://bafybeibyniiukxqmb7qae7ljif6atvo7ipg6wnpwvtqb4stf4ubjjterha.ipfs.w3s.link/manifest.json")
	t.Log(p)
}
