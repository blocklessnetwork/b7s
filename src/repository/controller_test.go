package repository

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPackage(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	// set test context and test appdb
	config := models.Config{}
	config.Node.WorkSpaceRoot = "/tmp/b7s_test"
	ctx = context.WithValue(ctx, "config", config)
	appDb := db.Get("/tmp/b7s_test/controller_testdb")
	ctx = context.WithValue(ctx, "appDb", appDb)

	// file uri reference manifest
	manifest := GetPackage(ctx, "https://bafybeiho3scwi3njueloobzhg7ndn7yjb5rkcaydvsoxmnhmu2adv6oxzq.ipfs.w3s.link/manifest.json")

	assert.Equal(manifest.Function.ID, "org.blockless.functions.myfunction", "manifest with known function id returned")

	// ask for the file again, should be cached
	manifest = GetPackage(ctx, "https://bafybeiho3scwi3njueloobzhg7ndn7yjb5rkcaydvsoxmnhmu2adv6oxzq.ipfs.w3s.link/manifest.json")
	assert.Equal(manifest.Cached, true, "manifest is marked as cached")

	db.Close(appDb)
}
