package overseer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
)

func TestOverseer_Kill(t *testing.T) {

	var (
		duration = 2 * time.Second
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
	// Start job.
	id, err := ov.Start(execJob)
	require.NoError(t, err)

	// Kill job.
	out, err := ov.Kill(id)
	require.NoError(t, err)
	end := time.Now()

	execTime := end.Sub(start)
	// Verify job status and verify that the execution time
	// was shorter than expected (command did not run until completion).
	require.Equal(t, job.StatusKilled, out.Status)
	require.Less(t, execTime, duration)
}
