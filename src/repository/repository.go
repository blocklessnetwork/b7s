package repository

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/http"
	"github.com/blocklessnetworking/b7s/src/models"
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

	if functionManifest.Runtime.Url != "" {
		DeploymentUrl, _ := url.Parse(functionManifest.Runtime.Url)
		if err != nil {
			log.Warn(err)
		}
		functionManifest.Deployment = models.Deployment{
			Uri:      DeploymentUrl.String(),
			Checksum: functionManifest.Runtime.Checksum,
		}
	}

	cachedFunction, err := db.GetString(ctx, functionManifest.Function.ID)
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

		db.Set(ctx, functionManifest.Function.ID, string(functionManifestJson))

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

func UnGzip(archive, destination string) error {
	dest := destination
	if len(dest) == 0 {
		log.WithFields(log.Fields{
			"archive": archive,
		}).Warn("extract archive destination is not specified and cwd is used.")
		dest = "./"
	} else {

		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}
	gzipStream, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer gzipStream.Close()
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.WithFields(log.Fields{
			"archive": archive,
			"err":     err,
		}).Error("extract archive failed.")
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.WithFields(log.Fields{
				"archive": archive,
				"err":     err,
			}).Error("extract archive failed.")
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			dir := filepath.Join(dest, header.Name)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.WithFields(log.Fields{
					"archive": archive,
					"err":     err,
					"dir":     dir,
				}).Error("extract archive mkdir failed")
				return err
			}
		case tar.TypeReg:
			file := filepath.Join(dest, header.Name)
			outFile, err := os.Create(file)
			if err != nil {
				log.WithFields(log.Fields{
					"archive": archive,
					"err":     err,
					"file":    file,
				}).Error("extract archive create new file failed")
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.WithFields(log.Fields{
					"archive": archive,
					"err":     err,
					"file":    file,
				}).Error("extract archive copy content to new file failed")
				outFile.Close()
				return err
			}
			outFile.Close()

		default:
			log.WithFields(log.Fields{
				"archive":  archive,
				"typeflag": header.Typeflag,
				"header":   header.Name,
			}).Error("extract archive unknown header")
			return errors.New("extract tgz file failed: unknown header")
		}

	}

	return nil
}
