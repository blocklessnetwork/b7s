package config

type configGroup uint

const (
	rootGroup = iota + 1
	logGroup
	connectivityGroup
	headGroup
	workerGroup
)

func (g configGroup) Name() string {

	switch g {
	case rootGroup:
		return ""
	case logGroup:
		return "log"
	case connectivityGroup:
		return "connectivity"
	case headGroup:
		return "head"
	case workerGroup:
		return "worker"
	default:
		return ""
	}
}
