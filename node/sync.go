package node

import (
	"context"
	"time"
)

// syncFunctions will try to redownload any functions that were removed from local disk
// but were previously installed. We do NOT abort on failure.
func (n *Node) syncFunctions() {

	cids := n.fstore.InstalledFunctions()

	for _, cid := range cids {

		err := n.fstore.Sync(cid)
		if err != nil {
			n.log.Error().Err(err).Str("cid", cid).Msg("function sync error")
			continue
		}

		n.log.Debug().Str("function", cid).Msg("function sync ok")
	}
}

func (n *Node) runSyncLoop(ctx context.Context) {

	ticker := time.NewTicker(syncInterval)

	for {
		select {
		case <-ticker.C:
			n.syncFunctions()

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
