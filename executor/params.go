package executor

import (
	"os"
)

const (
	defaultPermissions = os.ModePerm
	blocklessCli       = "blockless-cli" // TODO: On Windows we expect blockless-cli.exe
	blsListEnvName     = "BLS_LIST_VARS"
)
