package executor

import (
	"os"

	"github.com/blocklessnetworking/b7s/models/execute"
)

func procStateToUsage(ps *os.ProcessState) execute.Usage {

	usage := execute.Usage{
		CPUUserTime: ps.UserTime(),
		CPUSysTime:  ps.SystemTime(),
		MemoryMaxKB: getMemUsage(ps),
	}

	return usage
}
