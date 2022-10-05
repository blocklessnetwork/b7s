package repository

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/http"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/cockroachdb/pebble"
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

	appDb := ctx.Value("appDb").(*pebble.DB)
	cachedFunction := db.Value(appDb, functionManifest.Function.ID)

	if cachedFunction == "" {
		// download the function
		fileName, err := http.Download(ctx, functionManifest)

		if err != nil {
			log.Warn(err)
		}

		db.Set(appDb, functionManifest.Function.ID, fileName)

		log.WithFields(log.Fields{
			"uri": functionManifest.Deployment.Uri,
		}).Info("function sync completed")

	} else {
		log.WithFields(log.Fields{
			"uri": functionManifest.Deployment.Uri,
		}).Info("function sync skipped, already present")
	}

	return functionManifest
}
