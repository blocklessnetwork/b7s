package selection

import (
	"context"
	"testing"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
)

type mockNetwork struct {
	conns map[peer.ID]bool
}

func (n *mockNetwork) ConnsToPeer(id peer.ID) map[peer.ID]bool {
	return map[peer.ID]bool{id: n.conns[id]}
}

type mockHost struct {
	network Network
}

func (h *mockHost) Network() Network {
	return h.network
}

func TestSelectWorkerFromRollCall(t *testing.T) {
	network := &mockNetwork{
		conns: map[peer.ID]bool{
			"peer1": true,
		},
	}

	host := &mockHost{
		network: network,
	}

	rollcallResponse := &models.MsgRollCallResponse{
		From:       "peer1",
		Code:       enums.ResponseCodeAccepted,
		FunctionId: "function1",
		RequestId:  "request1",
	}

	request := &models.RequestExecute{
		FunctionId: "function1",
	}

	rollcallRequest := &models.MsgRollCall{
		RequestId: "request1",
	}

	ctx := context.WithValue(context.Background(), "host", host)

	result := SelectWorkerFromRollCall(ctx, *rollcallResponse, *request, *rollcallRequest)
	assert.Equal(t, false, result)

	rollcallResponse.FunctionId = "function2"
	result = SelectWorkerFromRollCall(ctx, *rollcallResponse, *request, *rollcallRequest)
	assert.Equal(t, true, result)

	rollcallResponse.Code = enums.ResponseCodeError
	result = SelectWorkerFromRollCall(ctx, *rollcallResponse, *request, *rollcallRequest)
	assert.Equal(t, true, result)
}
