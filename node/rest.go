package node

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
)

// ExecuteFunction can be used to start function execution. At the moment this is used by the API server to start execution on the head node.
func (n *Node) ExecuteFunction(ctx context.Context, req execute.Request) (codes.Code, execute.Result, execute.Cluster, error) {

	if !n.isHead() {
		return codes.NotAvailable, execute.Result{}, execute.Cluster{}, fmt.Errorf("action not supported on this node type")
	}

	requestID, err := newRequestID()
	if err != nil {
		return codes.Error, execute.Result{}, execute.Cluster{}, fmt.Errorf("could not generate request ID: %w", err)
	}

	return n.headExecute(ctx, requestID, req)
}

// ExecutionResult fetches the execution result from the node cache.
func (n *Node) ExecutionResult(id string) (execute.Result, bool) {
	res, ok := n.executeResponses.Get(id)
	return res.(execute.Result), ok
}

// PublishFunctionInstall publishes a function install message.
func (n *Node) PublishFunctionInstall(ctx context.Context, uri string, cid string) error {

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
		ManifestURL: manifestURLFromCID(cid),
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
