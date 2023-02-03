package controller

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/messaging"
	"github.com/blocklessnetworking/b7s/src/models"
)

// MsgInstallFunction publishes the function install message to pubsub.
func MsgInstallFunction(ctx context.Context, req models.RequestFunctionInstall) error {

	if req.Uri == "" && req.Cid == "" {
		return errors.New("invalid request - URI and CID are empty")
	}

	var msg models.MsgInstallFunction
	if req.Uri != "" {
		var err error
		msg, err = createInstallMessageFromURI(req.Uri)
		if err != nil {
			return fmt.Errorf("could not create install message from URI: %W", err)
		}
	} else {
		msg = createInstallMessageFromCID(req.Cid)
	}

	log.WithField("url", msg.ManifestUrl).
		Info("Requesting to message peer for function installation")

	// Get the pubsub topic from the context.
	topic, ok := ctx.Value("topic").(*pubsub.Topic)
	if !ok {
		return fmt.Errorf("unexpected value for pubsub topic (got: %T)", topic)
	}

	// Write the message to pubsub topic.
	messaging.PublishMessage(ctx, topic, msg)

	return nil
}

// createInstallMessageFromURI creates a MsgInstallFunction from the given URI.
// CID is calculated as a SHA-256 hash of the URI.
func createInstallMessageFromURI(uri string) (models.MsgInstallFunction, error) {

	h := sha256.New()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return models.MsgInstallFunction{}, fmt.Errorf("could not calculate hash: %w", err)
	}
	cid := fmt.Sprintf("%x", h.Sum(nil))

	msg := models.MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: uri,
		Cid:         cid,
	}

	return msg, nil
}

// createInstallMessageFromCID creates the MsgInstallFunction from the given CID.
func createInstallMessageFromCID(cid string) models.MsgInstallFunction {

	msg := models.MsgInstallFunction{
		Type:        enums.MsgInstallFunction,
		ManifestUrl: fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid),
		Cid:         cid,
	}

	return msg
}
