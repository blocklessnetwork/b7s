package repository

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/http"
	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

type JSONRepository struct{}

// get the manifest from the repository
// downloads the binary
func (r JSONRepository) Get(ctx context.Context, manifestPath string) models.FunctionManifest {

	functionManifest := models.FunctionManifest{}
	err := http.GetJson(manifestPath, &functionManifest)

	if err != nil {
		log.Warn(err)
	}

	http.Download(ctx, functionManifest)

	log.WithFields(log.Fields{
		"uri": functionManifest.Deployment.URI,
	}).Info("function sync completed")

	return functionManifest
}
