package telemetry

import (
	"context"
)

type ShutdownFunc func(context.Context) error
