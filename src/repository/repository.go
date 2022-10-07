package repository

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

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
	WorkSpaceRoot := ctx.Value("config").(models.Config).Node.WorkspaceRoot

	functionManifest := models.FunctionManifest{}
	err := http.GetJson(manifestPath, &functionManifest)

	if err != nil {
		log.Warn(err)
	}

	appDb := ctx.Value("appDb").(*pebble.DB)
	cachedFunction, err := db.Value(appDb, functionManifest.Function.ID)
	WorkSpaceDirectory := WorkSpaceRoot + "/" + functionManifest.Function.ID

	if err != nil {
		if err.Error() == "pebble: not found" {
			log.Info("function not found in cache, syncing")
		} else {
			log.Warn(err)
		}
	}

	if cachedFunction == "" {
		// download the function
		fileName, err := http.Download(ctx, functionManifest)

		UnGzip(fileName, WorkSpaceDirectory)
		if err != nil {
			log.Warn(err)
		}

		functionManifest.Deployment.File = fileName
		functionManifest.Cached = true

		functionManifestJson, error := json.Marshal(functionManifest)

		if error != nil {
			log.Warn(error)
		}

		db.Set(appDb, functionManifest.Function.ID, string(functionManifestJson))

		log.WithFields(log.Fields{
			"uri": functionManifest.Deployment.Uri,
		}).Info("function sync completed")

	} else {

		err := json.Unmarshal([]byte(cachedFunction), &functionManifest)
		if err != nil {
			log.Warn(err)
		}

		log.WithFields(log.Fields{
			"uri": functionManifest.Deployment.Uri,
		}).Info("function sync skipped, already present")
	}

	return functionManifest
}

func UnGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, "inflated")
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}
