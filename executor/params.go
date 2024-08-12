package executor

import (
	"os"
)

const (
	defaultPermissions = os.ModePerm
	blsListEnvName     = "BLS_LIST_VARS"
	tracerName         = "b7s.Executor"
)

var (
	functionExecutionsMetric  = []string{"executor", "function", "executions"}
	functionDurationMetric    = []string{"executor", "function", "executions", "milliseconds"}
	functionCPUUserTimeMetric = []string{"executor", "function", "executions", "cpu", "user", "time", "milliseconds"}
	functionCPUSysTimeMetric  = []string{"executor", "function", "executions", "cpu", "sys", "time", "milliseconds"}
	functionOkMetric          = []string{"executor", "function", "executions", "ok"}
	functionErrMetric         = []string{"executor", "function", "executions", "err"}
)
