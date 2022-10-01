package repository

import (
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/blocklessnetworking/b7s/src/rest"
	log "github.com/sirupsen/logrus"
)

type JSONRepository struct {
	Endpoint string
}

func (r JSONRepository) Get() models.FunctionManifest {
	jsonResponse := models.FunctionManifest{}
	err := rest.GetJson(r.Endpoint, &jsonResponse)

	if err != nil {
		log.Warn(err)
	}

	return jsonResponse
}
