package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

var RestClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := RestClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func Download(ctx context.Context, functionManifest models.FunctionManifest) error {
	WorkSpaceRoot := ctx.Value("config").(models.Config).Node.WorkSpaceRoot
	WorkSpaceDirectory := WorkSpaceRoot + "/" + functionManifest.Function.ID
	client := grab.NewClient()

	// ensure path exists
	os.MkdirAll(WorkSpaceDirectory, os.ModePerm)
	// download function
	req, _ := grab.NewRequest(WorkSpaceDirectory, functionManifest.Deployment.URI)
	resp := client.Do(req)

	log.WithFields(log.Fields{
		"uri": functionManifest.Deployment.URI,
	}).Info("function scheduled for sync")

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			log.WithFields(log.Fields{
				"uri": functionManifest.Deployment.URI,
			}).Info("function sync progress")

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		log.WithFields(log.Fields{
			"uri": functionManifest.Deployment.URI,
		}).Info("function sync field will try again")
		return err
	}

	return nil
}
