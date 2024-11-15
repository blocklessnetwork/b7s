package worker

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
)

func (w *Worker) processInstallFunction(ctx context.Context, from peer.ID, req request.InstallFunction) error {

	// Install function.
	err := w.installFunction(ctx, req.CID, req.ManifestURL)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	// Reply to the caller.
	err = w.Send(ctx, from, req.Response(codes.Accepted))
	if err != nil {
		return fmt.Errorf("could not send the response (peer: %s): %w", from, err)
	}

	return nil
}

// installFunction will check if the function is installed first, and install it if not.
func (w *Worker) installFunction(ctx context.Context, cid string, manifestURL string) error {

	// Check if the function is installed.
	installed, err := w.fstore.IsInstalled(cid)
	if err != nil {
		return fmt.Errorf("could not check if function is installed: %w", err)
	}

	if installed {
		return nil
	}

	// If the function was not installed already, install it now.
	err = w.fstore.Install(ctx, manifestURL, cid)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	return nil
}
