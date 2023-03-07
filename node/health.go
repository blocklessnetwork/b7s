package node

import (
	"context"
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/response"
)

// HealthPing will run a long running loop, publishing health signal until cancelled.
func (n *Node) HealthPing(ctx context.Context) {

	ticker := time.NewTicker(n.cfg.HealthInterval)

	for {
		select {

		case <-ticker.C:

			msg := response.Health{
				Type: blockless.MessageHealthCheck,
				Code: response.CodeOK,
			}

			err := n.publish(ctx, msg)
			if err != nil {
				n.log.Warn().Err(err).Msg("could not publish health signal")
			}

			n.log.Debug().Msg("emitted health ping")

		case <-ctx.Done():
			ticker.Stop()
			n.log.Info().Msg("stopping health ping")
			return
		}
	}
}
