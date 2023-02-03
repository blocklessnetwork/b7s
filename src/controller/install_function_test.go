package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/src/enums"
)

func TestCreateInstallMessageFromURI(t *testing.T) {

	const (
		uri = "https://example.com/manifest.json"
	)

	req, err := createInstallMessageFromURI(uri)
	require.NoError(t, err)

	assert.Equal(t, enums.MsgInstallFunction, req.Type)
	assert.Equal(t, uri, req.ManifestUrl)

	cid, err := deriveCIDFromURI(uri)
	require.NoError(t, err)

	assert.Equal(t, cid, req.Cid)
}

func TestCreateInstallMessageFromCID(t *testing.T) {

	const (
		cid                 = "test-cid-value"
		expectedManifestURL = `https://test-cid-value.ipfs.w3s.link/manifest.json`
	)

	req := createInstallMessageFromCID(cid)

	assert.Equal(t, enums.MsgInstallFunction, req.Type)
	assert.Equal(t, cid, req.Cid)
	assert.Equal(t, expectedManifestURL, req.ManifestUrl)
}
