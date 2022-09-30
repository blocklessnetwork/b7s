package repository

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func startClient() {

}

func Start(ctx context.Context) {
	// var config = ctx.Value("config").(models.Config)

	log.Info("starting repository client")
	go startClient()
}
