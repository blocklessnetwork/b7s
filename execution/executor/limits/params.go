package limits

const (
	DefaultCgroup        = "/blockless"
	DefaultMountpoint    = "/sys/fs/cgroup"
	DefaultJobObjectName = "blockless"

	// Default percentage of the CPU allowed. By default we run unlimited.
	DefaultCPUPercentage = 1.0
)
