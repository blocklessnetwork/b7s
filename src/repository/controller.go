package repository

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/models"
)

func GetPackage(ctx context.Context, manifestPath string) (models.FunctionManifest, error) {
	repo := JSONRepository{}
	return repo.Get(ctx, manifestPath)
}
