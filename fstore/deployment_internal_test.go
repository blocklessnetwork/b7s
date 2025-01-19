package fstore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/models/bls"
)

func TestFunction_UpdateDeploymentInfo(t *testing.T) {

	t.Run("noop - deployment info present", func(t *testing.T) {
		t.Parallel()

		const (
			runtimeURL  = "https://example.com/runtime-address"
			checksum    = "123456789"
			manifestURL = "https://example.com"
		)

		manifest := bls.FunctionManifest{
			Runtime: bls.Runtime{
				URL:      runtimeURL,
				Checksum: checksum,
			},
			Deployment: bls.Deployment{
				URI:      "",
				Checksum: "",
			},
		}

		err := updateDeploymentInfo(&manifest, manifestURL)
		require.NoError(t, err)

		require.Equal(t, runtimeURL, manifest.Deployment.URI)
		require.Equal(t, checksum, manifest.Deployment.Checksum)
	})
	t.Run("fills in missing host or scheme info", func(t *testing.T) {
		t.Parallel()

		const (
			runtimeURL  = "runtime-address"
			checksum    = "123456789"
			manifestURL = "https://example.com/manifest-address"
		)

		manifest := bls.FunctionManifest{
			Runtime: bls.Runtime{
				URL:      runtimeURL,
				Checksum: checksum,
			},
			Deployment: bls.Deployment{
				URI:      "",
				Checksum: "",
			},
		}

		err := updateDeploymentInfo(&manifest, manifestURL)
		require.NoError(t, err)

		manifestAddress, err := url.Parse(manifestURL)
		require.NoError(t, err)

		deploymentURL := url.URL{
			Host:   manifestAddress.Host,
			Scheme: manifestAddress.Scheme,
			Path:   runtimeURL,
		}

		require.Equal(t, deploymentURL.String(), manifest.Deployment.URI)
		require.Equal(t, checksum, manifest.Deployment.Checksum)
	})
	t.Run("handles malformed runtime URL", func(t *testing.T) {
		t.Parallel()

		const (
			runtimeURL  = "http://example.com/runtime-address\n"
			checksum    = "123456789"
			manifestURL = "https://example.com/manifest-address"
		)

		manifest := bls.FunctionManifest{
			Runtime: bls.Runtime{
				URL:      runtimeURL,
				Checksum: checksum,
			},
			Deployment: bls.Deployment{
				URI:      "",
				Checksum: "",
			},
		}

		err := updateDeploymentInfo(&manifest, manifestURL)
		require.Error(t, err)
	})
	t.Run("handles malformed manifest URL", func(t *testing.T) {
		t.Parallel()

		const (
			runtimeURL  = "http://example.com/runtime-address"
			checksum    = "123456789"
			manifestURL = "https://example.com/manifest-address\r"
		)

		manifest := bls.FunctionManifest{
			Runtime: bls.Runtime{
				URL:      runtimeURL,
				Checksum: checksum,
			},
			Deployment: bls.Deployment{
				URI:      "",
				Checksum: "",
			},
		}

		err := updateDeploymentInfo(&manifest, manifestURL)
		require.Error(t, err)
	})
}
