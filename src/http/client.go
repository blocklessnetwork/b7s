package http

import (
	"context"
	"encoding/json"
	"net/http"
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
		return err
	}

	return nil
}
