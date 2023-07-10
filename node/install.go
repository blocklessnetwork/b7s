package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

func (n *Node) processInstallFunction(ctx context.Context, from peer.ID, payload []byte) error {

	// Only workers should respond to function install requests.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Debug().Msg("received function install request, ignoring")
		return nil
	}

	// Unpack the request.
	var req request.InstallFunction
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack request: %w", err)
	}
	req.From = from

	// Install function.
	err = n.installFunction(req.CID, req.ManifestURL)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	// Create the response.
	res := response.InstallFunction{
		Type:    blockless.MessageInstallFunctionResponse,
		Code:    codes.Accepted,
		Message: "installed",
	}

	// Reply to the caller.
	err = n.send(ctx, from, res)
	if err != nil {
		return fmt.Errorf("could not send the response (peer: %s): %w", from, err)
	}

	return nil
}

// installFunction will check if the function is installed first, and install it if not.
func (n *Node) installFunction(cid string, manifestURL string) error {

	// Check if the function is installed.
	installed, err := n.fstore.Installed(cid)
	if err != nil {
		return fmt.Errorf("could not check if function is installed: %w", err)
	}

	if installed {
		return nil
	}

	// If the function was not installed already, install it now.
	err = n.fstore.Install(manifestURL, cid)
	if err != nil {
		return fmt.Errorf("could not install function: %w", err)
	}

	return nil
}

func manifestURLFromCID(cid string) string {
	return fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid)
}
