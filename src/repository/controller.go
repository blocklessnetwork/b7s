package repository

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/models"
)

func WithEndpoint(ctx context.Context, endPoint string) JSONRepository {
	repo := JSONRepository{
		Endpoint: endPoint,
	}
	return repo
}

func GetPackage(ctx context.Context) []models.RepoPackage {
	repo := JSONRepository{
		Endpoint: ctx.Value("config").(models.Config).Repository.Url,
	}
	return repo.List()
}

func ListPackages(ctx context.Context) models.RepoPackage {
	repo := JSONRepository{
		Endpoint: ctx.Value("config").(models.Config).Repository.Url,
	}
	return repo.Get()
}
