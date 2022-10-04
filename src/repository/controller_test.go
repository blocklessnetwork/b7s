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

	p := GetPackage(ctx, "http://localhost:8080/someid/manifest.json")
	t.Log(p)
}
