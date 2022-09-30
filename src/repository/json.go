package repository

import "github.com/blocklessnetworking/b7s/src/models"

type JSONRepository struct {
	Endpoint string
}

func (r JSONRepository) GetEndpoint() string {
	return ""
}

func (r JSONRepository) SetEndpoint() {

}

func (r JSONRepository) List() []models.RepoPackage {
	return []models.RepoPackage{}
}

func (r JSONRepository) Get() models.RepoPackage {
	return models.RepoPackage{}
}
