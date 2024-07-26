package node

import (
	"context"
	"time"
)

func (n *Node) runSyncLoop(ctx context.Context) {

	ticker := time.NewTicker(syncInterval)

	for {
		select {
		case <-ticker.C:
			err := n.fstore.Sync(ctx, false)
			if err != nil {
				n.log.Error().Err(err).Msg("function sync unsuccessful")
			} else {
				n.log.Debug().Msg("function sync ok")
			}

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
