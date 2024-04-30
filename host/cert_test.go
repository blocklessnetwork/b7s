package host

import (
	"crypto/rand"
	"testing"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"
)

func TestConvertLibp2pPrivKeyToCryptoPrivKey(t *testing.T) {
	// Generate a libp2p ECDSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateECDSAKeyPair(rand.Reader)
	assert.NoError(t, err, "failed to generate libp2p ECDSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	assert.NoError(t, err, "failed to convert libp2p private key to crypto private key")
	assert.NotNil(t, cryptoPriv, "converted crypto private key should not be nil")
}

func TestGenerateX509Certificate(t *testing.T) {
	// Generate a libp2p ECDSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateECDSAKeyPair(rand.Reader)
	assert.NoError(t, err, "failed to generate libp2p ECDSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	assert.NoError(t, err, "failed to convert libp2p private key")

	// Generate an X.509 certificate
	cert, err := generateX509Certificate(cryptoPriv)
	assert.NoError(t, err, "failed to generate X.509 certificate")
	assert.NotEmpty(t, cert.Certificate, "certificate should contain at least one DER encoded block")
}
