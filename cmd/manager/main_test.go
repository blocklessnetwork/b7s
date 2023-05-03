package main

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/stretchr/testify/assert"
)

func TestMakeBasicHost(t *testing.T) {
	listenPort := 0
	insecure := false
	randseed := int64(0)

	h, err := makeBasicHost(listenPort, insecure, randseed)
	assert.NoError(t, err)
	assert.NotNil(t, h)
	assert.Implements(t, (*host.Host)(nil), h)
}

func TestGetHostAddress(t *testing.T) {
	listenPort := 0
	insecure := false
	randseed := int64(0)

	host, _ := makeBasicHost(listenPort, insecure, randseed)
	address := getHostAddress(host)
	assert.Contains(t, address, "/p2p/")
}
