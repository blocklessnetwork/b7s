package node

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"sync"
	"testing"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_SendMessage(t *testing.T) {

	const (
		clientAddress = "127.0.0.1"
		clientPort    = 0
	)

	type record struct {
		ID          string `json:"id"`
		Value       uint64 `json:"value"`
		Description string `json:"description"`
	}

	var (
		rec = record{
			ID:          mocks.GenericUUID.String(),
			Value:       19846,
			Description: "dummy-description",
		}
	)

	client, err := host.New(mocks.NoopLogger, clientAddress, clientPort)
	require.NoError(t, err)

	clientAddresses := client.Addresses()
	require.NotEmpty(t, clientAddresses)

	addr := clientAddresses[0]

	node := createNode(t, blockless.HeadNode)
	addPeerToPeerStore(t, node.host, addr)

	var wg sync.WaitGroup
	wg.Add(1)

	client.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer wg.Done()
		defer stream.Close()

		from := stream.Conn().RemotePeer()
		require.Equal(t, node.host.ID(), from)

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		require.ErrorIs(t, err, io.EOF)

		var received record
		err = json.Unmarshal(payload, &received)
		require.NoError(t, err)

		require.Equal(t, rec, received)
	})

	err = node.send(context.Background(), client.ID(), rec)
	require.NoError(t, err)

	wg.Wait()
}
