package node

import (
	"context"
	"net/http"
	"time"

	"github.com/blessnetwork/b7s/models/response"
)

// HealthPing will run a long running loop, publishing health signal until cancelled.
func (c *core) emitHealthPing(ctx context.Context, interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:

			msg := response.Health{
				Code: http.StatusOK,
			}

			err := c.Publish(ctx, &msg)
			if err != nil {
				c.log.Warn().Err(err).Msg("could not publish health signal")
				return
			}

			c.log.Trace().Msg("emitted health ping")

		case <-ctx.Done():
			c.log.Info().Msg("stopping health ping")
			return
		}
	}
}
