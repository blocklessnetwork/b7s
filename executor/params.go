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
	moduleName                = "executor"
	functionExecutionsMetric  = []string{moduleName, "function", "executions"}
	functionDurationMetric    = append(functionExecutionsMetric, "milliseconds")
	functionCPUUserTimeMetric = append(functionExecutionsMetric, "cpu", "user", "time", "milliseconds")
	functionCPUSysTimeMetric  = append(functionExecutionsMetric, "cpu", "sys", "time", "milliseconds")
	functionOkMetric          = append(functionExecutionsMetric, "ok")
	functionErrMetric         = append(functionExecutionsMetric, "err")
)
