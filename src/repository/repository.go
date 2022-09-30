package repository

import (
	"context"
	"time"

	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

func getRepoUpdates(ctx context.Context) {
	var cfg = ctx.Value("config").(models.Config)
	repo := WithEndpoint(ctx, cfg.Repository.Url)

	log.Debug("fetching list of repository packages")
	packages := repo.List()

	log.WithFields(log.Fields{
		"packages": packages,
	}).Info("repository packages")
}

func startClient(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			getRepoUpdates(ctx)
		}
	}
}

func Start(ctx context.Context, ticker *time.Ticker) {
	// var config = ctx.Value("config").(models.Config)
	log.Debug("starting repository client")

	getRepoUpdates(ctx)

	// start health monitoring
	go startClient(ctx, ticker)
}
