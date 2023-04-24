package executor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/execute"
)

func TestRuntimeFlags(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		t.Parallel()

		cfg := execute.RuntimeConfig{}
		flags := runtimeFlags(cfg)
		require.Len(t, flags, 0)
	})
	t.Run("all flags set", func(t *testing.T) {
		t.Parallel()

		const (
			entry         = "something"
			input         = "whatever.wasm"
			executionTime = 123456
			debugInfo     = true
			fsRoot        = "/var/tmp/request"
			fuel          = 987
			memory        = 256
			logger        = "runtime.log"

			// Expect six key-value pairs and a single boolean.
			flagCount = 13
		)

		cfg := execute.RuntimeConfig{
			Entry:         entry,
			Input:         input,
			ExecutionTime: executionTime,
			DebugInfo:     debugInfo,
			FSRoot:        fsRoot,
			Fuel:          fuel,
			Memory:        memory,
			Logger:        logger,
		}

		flags := runtimeFlags(cfg)

		require.Len(t, flags, flagCount)

		require.Equal(t, "--"+execute.RuntimeFlagEntry, flags[0])
		require.Equal(t, entry, flags[1])

		require.Equal(t, "--"+execute.RuntimeFlagExecutionTime, flags[2])
		require.Equal(t, fmt.Sprint(executionTime), flags[3])

		require.Equal(t, "--"+execute.RuntimeFlagDebug, flags[4])

		require.Equal(t, "--"+execute.RuntimeFlagFSRoot, flags[5])
		require.Equal(t, fsRoot, flags[6])

		require.Equal(t, "--"+execute.RuntimeFlagFuel, flags[7])
		require.Equal(t, fmt.Sprint(fuel), flags[8])

		require.Equal(t, "--"+execute.RuntimeFlagMemory, flags[9])
		require.Equal(t, fmt.Sprint(memory), flags[10])

		require.Equal(t, "--"+execute.RuntimeFlagLogger, flags[11])
		require.Equal(t, logger, flags[12])
	})
	t.Run("some fields set", func(t *testing.T) {
		t.Parallel()

		const (
			entry  = "something"
			memory = 256
		)

		cfg := execute.RuntimeConfig{
			Entry:  entry,
			Memory: memory,
		}

		flags := runtimeFlags(cfg)

		require.Len(t, flags, 4)

		require.Equal(t, "--"+execute.RuntimeFlagEntry, flags[0])
		require.Equal(t, entry, flags[1])

		require.Equal(t, "--"+execute.RuntimeFlagMemory, flags[2])
		require.Equal(t, fmt.Sprint(memory), flags[3])
	})
}
