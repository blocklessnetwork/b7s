package node

import (
	"context"
	"net/http"
	"time"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/response"
)

// HealthPing will run a long running loop, publishing health signal until cancelled.
func (n *Node) HealthPing(ctx context.Context) {

	ticker := time.NewTicker(n.cfg.HealthInterval)

	for {
		select {

		case <-ticker.C:

			msg := response.Health{
				Type: blockless.MessageHealthCheck,
				Code: http.StatusOK,
			}

			err := n.publish(ctx, msg)
			if err != nil {
				n.log.Warn().Err(err).Msg("could not publish health signal")
			}

			n.log.Trace().Msg("emitted health ping")

		case <-ctx.Done():
			ticker.Stop()
			n.log.Info().Msg("stopping health ping")
			return
		}
	}
}
