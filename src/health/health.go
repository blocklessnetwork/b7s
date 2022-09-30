package health

import (
	"context"
	"time"

	"github.com/blocklessnetworking/b7s/src/controller"
)

func StartPing(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			controller.HealthStatus(ctx)
		}
	}
}
