package limits

const (
	DefaultCgroup        = "/bless"
	DefaultMountpoint    = "/sys/fs/cgroup"
	DefaultJobObjectName = "bless"

	// Default percentage of the CPU allowed. By default we run unlimited.
	DefaultCPUPercentage = 1.0
)
