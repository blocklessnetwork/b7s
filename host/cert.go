package host

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

// Convert a libp2p PrivKey to a crypto.PrivateKey
func convertLibp2pPrivKeyToCryptoPrivKey(privKey libp2pcrypto.PrivKey) (crypto.PrivateKey, error) {
	rawKey, err := privKey.Raw()
	if err != nil {
		return nil, err
	}

	switch privKey.Type() {
	case libp2pcrypto.RSA:
		return x509.ParsePKCS1PrivateKey(rawKey)
	case libp2pcrypto.ECDSA:
		return x509.ParseECPrivateKey(rawKey)
	case libp2pcrypto.Ed25519:
		return ed25519.PrivateKey(rawKey), nil
	default:
		return nil, fmt.Errorf("unsupported key type for X.509 conversion")
	}
}

func generateX509Certificate(privKey crypto.PrivateKey) (tls.Certificate, error) {
	// Define certificate template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"b7s"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	pubKey := publicKey(privKey)

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Encode the certificate and private key
	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privKey,
	}

	return cert, nil
}

func publicKey(priv crypto.PrivateKey) crypto.PublicKey {
	switch key := priv.(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey
	case *ecdsa.PrivateKey:
		return &key.PublicKey
	case ed25519.PrivateKey:
		return key.Public().(ed25519.PublicKey)
	default:
		panic("unsupported key type")
	}
}
