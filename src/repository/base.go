package repository

import (
	"github.com/blocklessnetworking/b7s/src/models"
)

type Repo interface {
	GetEndpoint(endpoint string)
	SetEndpoint(endpoint string)
	List() []models.RepoPackage
	Get() string
}
