package node

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
)

// TODO: Consider introducing an entity - a `delegator`. This could be like an Executor, only
// instead of local execution, it would issue a roll call and delegate work to the worker nodes.
// Problem is that delegator would need to be notified when an execution result has arrived.
// Doing this way would make the execution flow more streamlined and would not differentiate as much between
// worker and head node.
func (n *Node) ExecuteFunction(ctx context.Context, req execute.Request) (execute.Result, error) {

	requestID, err := newRequestID()
	if err != nil {
		return execute.Result{}, fmt.Errorf("could not generate request ID: %w", err)
	}

	switch n.cfg.Role {
	case blockless.WorkerNode:
		return n.workerExecute(ctx, n.host.ID(), requestID, req)

	case blockless.HeadNode:
		return n.headExecute(ctx, n.host.ID(), requestID, req)
	}

	panic(fmt.Errorf("invalid node role: %s", n.cfg.Role))
}

// ExecutionResult fetches the execution result from the node cache.
func (n *Node) ExecutionResult(id string) (execute.Result, bool) {
	res, ok := n.excache.Get(id)
	return res, ok
}

// FunctionInstall initiates function install process.
func (n *Node) FunctionInstall(ctx context.Context, uri string, cid string) error {

	var req request.InstallFunction
	if uri != "" {
		var err error
		req, err = createInstallMessageFromURI(uri)
		if err != nil {
			return fmt.Errorf("could not create install message from URI: %W", err)
		}
	} else {
		req = createInstallMessageFromCID(cid)
	}

	n.log.Debug().
		Str("url", req.ManifestURL).
		Str("cid", req.CID).
		Msg("publishing function install message")

	err := n.publish(ctx, req)
	if err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}

	return nil
}

// createInstallMessageFromURI creates a MsgInstallFunction from the given URI.
// CID is calculated as a SHA-256 hash of the URI.
func createInstallMessageFromURI(uri string) (request.InstallFunction, error) {

	cid, err := deriveCIDFromURI(uri)
	if err != nil {
		return request.InstallFunction{}, fmt.Errorf("could not determine cid: %w", err)
	}

	msg := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		ManifestURL: uri,
		CID:         cid,
	}

	return msg, nil
}

// createInstallMessageFromCID creates the MsgInstallFunction from the given CID.
func createInstallMessageFromCID(cid string) request.InstallFunction {

	req := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		ManifestURL: fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid),
		CID:         cid,
	}

	return req
}

func deriveCIDFromURI(uri string) (string, error) {

	h := sha256.New()
	_, err := h.Write([]byte(uri))
	if err != nil {
		return "", fmt.Errorf("could not calculate hash: %w", err)
	}
	cid := fmt.Sprintf("%x", h.Sum(nil))

	return cid, nil
}
