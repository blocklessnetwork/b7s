package helpers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

type FunctionServer struct {
	*httptest.Server
}

func CreateFunctionServer(t *testing.T, manifestEndpoint string, manifest blockless.FunctionManifest, deploymentEndpoint string, archivePath string, cid string) *FunctionServer {
	t.Helper()

	// Archive to serve.
	archive, err := os.ReadFile(archivePath)
	require.NoError(t, err)

	// Checksum of the archive we serve.
	checksum := sha256.Sum256(archive)

	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			path := req.URL.Path
			switch path {
			// Manifest request.
			case manifestEndpoint:

				// Link to a URL on our own server where we'll serve the function archive.
				deploymentURL := url.URL{
					Scheme: "http",
					Host:   req.Host,
					Path:   deploymentEndpoint,
				}

				manifest.Deployment = blockless.Deployment{
					CID:      cid,
					Checksum: fmt.Sprintf("%x", checksum),
					URI:      deploymentURL.String(),
				}

				payload, err := json.Marshal(manifest)
				require.NoError(t, err)
				w.Write(payload)

			// Archive download request.
			case deploymentEndpoint:
				w.Write(archive)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

	fs := FunctionServer{
		Server: srv,
	}

	return &fs
}
