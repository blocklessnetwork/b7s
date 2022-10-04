package repository

import (
	"context"
	"time"

	"github.com/blocklessnetworking/b7s/src/http"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

type JSONRepository struct{}

// get the manifest from the repository
// downloads the binary
func (r JSONRepository) Get(ctx context.Context, manifestPath string) models.FunctionManifest {
	WorkSpaceRoot := ctx.Value("config").(models.Config).Node.WorkSpaceRoot

	functionManifest := models.FunctionManifest{}
	err := http.GetJson(manifestPath, &functionManifest)

	if err != nil {
		log.Warn(err)
	}

	client := grab.NewClient()
	req, _ := grab.NewRequest(WorkSpaceRoot+"/"+functionManifest.Id, functionManifest.Runtime.Uri)
	resp := client.Do(req)

	log.WithFields(log.Fields{
		"uri": functionManifest.Runtime.Uri,
	}).Info("function scheduled for sync")

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			log.WithFields(log.Fields{
				"uri": functionManifest.Runtime.Uri,
			}).Info("function sync progress")

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		log.WithFields(log.Fields{
			"uri": functionManifest.Runtime.Uri,
		}).Info("function sync field will try again")
	}

	log.WithFields(log.Fields{
		"uri": functionManifest.Runtime.Uri,
	}).Info("function sync completed")

	return functionManifest
}
