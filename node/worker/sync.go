package worker

import (
	"context"
	"time"
)

func (w *Worker) runSyncLoop(ctx context.Context) {

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.fstore.Sync(ctx, false)
			if err != nil {
				w.Log().Error().Err(err).Msg("function sync unsuccessful")
			} else {
				w.Log().Debug().Msg("function sync ok")
			}

		case <-ctx.Done():
			return
		}
	}
}
