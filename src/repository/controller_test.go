package repository

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/src/models"
)

func TestGetPackage(t *testing.T) {
	ctx := context.Background()
	config := models.Config{
		Repository: struct {
			Url string "yaml:\"url\""
		}{
			Url: "http://localhost:8080/manifest.json",
		},
	}
	ctx = context.WithValue(ctx, "config", config)
	p := GetPackage(ctx)
	t.Log(p)
}
