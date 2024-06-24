package node

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processInstallFunction(ctx context.Context, from peer.ID, req request.InstallFunction) error {

	// Only workers should respond to function install requests.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Debug().Msg("received function install request, ignoring")
		return nil
	}

	// Install function.
	err := n.installFunction(ctx, req.CID, req.ManifestURL)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	// Create the response.
	res := response.InstallFunction{
		Code:    codes.Accepted,
		Message: "installed",
		CID:     req.CID,
	}

	// Reply to the caller.
	err = n.send(ctx, from, &res)
	if err != nil {
		return fmt.Errorf("could not send the response (peer: %s): %w", from, err)
	}

	return nil
}

// installFunction will check if the function is installed first, and install it if not.
func (n *Node) installFunction(ctx context.Context, cid string, manifestURL string) error {

	// Check if the function is installed.
	installed, err := n.fstore.IsInstalled(cid)
	if err != nil {
		return fmt.Errorf("could not check if function is installed: %w", err)
	}

	if installed {
		return nil
	}

	// If the function was not installed already, install it now.
	err = n.fstore.Install(ctx, manifestURL, cid)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	return nil
}

func manifestURLFromCID(cid string) string {
	return fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid)
}
