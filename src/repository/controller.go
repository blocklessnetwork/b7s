package repository

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/models"
)

func GetPackage(ctx context.Context) models.FunctionManifest {
	repo := JSONRepository{
		Endpoint: ctx.Value("config").(models.Config).Repository.Url,
	}
	return repo.Get()
}
