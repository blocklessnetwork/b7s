package overseer

import (
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

var (
	runnableExePath = "/home/aco/code/Maelkum/overseer/tools/runnable/runnable"
)

func TestOverseer_Run(t *testing.T) {

	var (
		duration = 1 * time.Second
		stdout   = fmt.Sprintf("test-string-%v", rand.Int())
		stderr   = fmt.Sprintf("test-string-%v", rand.Int())

		execJob = job.Job{
			Exec: job.Command{
				Path: runnableExePath,
				Args: []string{

					"--stdout", stdout,
					"--stderr", stderr,
					"--duration", duration.String(),
				},
			},
		}
	)

	ov := createOverseer(t, runnableExePath)

	start := time.Now()

	out, err := ov.Run(execJob)
	require.NoError(t, err)

	end := time.Now()

	require.Equal(t, job.StatusDone, out.Status)
	require.Equal(t, stdout, out.Stdout)
	require.Equal(t, stderr, out.Stderr)

	require.NotNil(t, out.EndTime)
	require.Greater(t, *out.EndTime, out.StartTime)

	require.GreaterOrEqual(t, out.StartTime, start)
	require.LessOrEqual(t, *out.EndTime, end)

	// Not verifying CPU times.

	require.NotZero(t, out.ResourceUsage.MemoryMaxKB)
}

func createOverseer(t *testing.T, exe string) *Overseer {
	t.Helper()

	ov, err := New(
		zerolog.New(io.Discard),
		WithAllowlist(exe),
	)
	require.NoError(t, err)

	return ov
}
