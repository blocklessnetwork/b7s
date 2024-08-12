package fstore

import (
	"time"
)

const (
	defaultTimeout   = 10 * time.Second
	defaultUserAgent = "b7s"

	tracerName = "b7s.Fstore"
)

// Tracing span names.
const (
	spanInstall     = "FunctionInstall"
	spanIsInstalled = "IsFunctionInstalled"
	spanSync        = "FunctionSync"
)

var (
	functionsInstalledMetric = []string{"functions", "installed"}
)
