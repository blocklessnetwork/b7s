package crypto

import (
	"crypto/rand"
	"testing"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/require"
)

func TestConvertLibp2pPrivKeyToCryptoPrivKey(t *testing.T) {
	// Generate a libp2p ECDSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateECDSAKeyPair(rand.Reader)
	require.NoError(t, err, "failed to generate libp2p ECDSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key to crypto private key")
	require.NotNil(t, cryptoPriv, "converted crypto private key should not be nil")
}

func TestConvertLibp2pRSAPrivKeyToCryptoPrivKey(t *testing.T) {
	// Generate a libp2p RSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateRSAKeyPair(2048, rand.Reader)
	require.NoError(t, err, "failed to generate libp2p RSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key to crypto private key")
	require.NotNil(t, cryptoPriv, "converted crypto private key should not be nil")
}

func TestConvertLibp2pEd25519PrivKeyToCryptoPrivKey(t *testing.T) {
	// Generate a libp2p Ed25519 key pair for testing
	priv, _, err := libp2pcrypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err, "failed to generate libp2p Ed25519 key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key to crypto private key")
	require.NotNil(t, cryptoPriv, "converted crypto private key should not be nil")
}

func TestGenerateX509Certificate(t *testing.T) {
	// Generate a libp2p ECDSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateECDSAKeyPair(rand.Reader)
	require.NoError(t, err, "failed to generate libp2p ECDSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key")

	// Generate an X.509 certificate
	cert, err := generateX509Certificate(cryptoPriv)
	require.NoError(t, err, "failed to generate X.509 certificate")
	require.NotEmpty(t, cert.Certificate, "certificate should contain at least one DER encoded block")
}

func TestGenerateX509CertificateRSA(t *testing.T) {
	// Generate a libp2p RSA key pair for testing
	priv, _, err := libp2pcrypto.GenerateRSAKeyPair(2048, rand.Reader)
	require.NoError(t, err, "failed to generate libp2p RSA key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key")

	// Generate an X.509 certificate
	cert, err := generateX509Certificate(cryptoPriv)
	require.NoError(t, err, "failed to generate X.509 certificate")
	require.NotEmpty(t, cert.Certificate, "certificate should contain at least one DER encoded block")
}

func TestGenerateX509CertificateEd25519(t *testing.T) {
	// Generate a libp2p Ed25519 key pair for testing
	priv, _, err := libp2pcrypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err, "failed to generate libp2p Ed25519 key pair")

	// Convert the libp2p private key to a crypto.PrivateKey
	cryptoPriv, err := convertLibp2pPrivKeyToCryptoPrivKey(priv)
	require.NoError(t, err, "failed to convert libp2p private key")

	// Generate an X.509 certificate
	cert, err := generateX509Certificate(cryptoPriv)
	require.NoError(t, err, "failed to generate X.509 certificate")
	require.NotEmpty(t, cert.Certificate, "certificate should contain at least one DER encoded block")
}
