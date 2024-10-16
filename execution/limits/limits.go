package limits

var DefaultLimits = Limits{
	CPUPercentage: 1.0,
	MemoryKB:      -1,
	ProcLimit:     0,
}

type Limits struct {
	CPUPercentage float64
	MemoryKB      int64
	ProcLimit     uint
}

type LimitOption func(*Limits)

func WithCPUPercentage(p float64) LimitOption {
	return func(l *Limits) {
		l.CPUPercentage = p
	}
}

func WithMemoryKB(m int64) LimitOption {
	return func(l *Limits) {
		l.MemoryKB = m
	}
}

func WithProcLimit(n uint) LimitOption {
	return func(l *Limits) {
		l.ProcLimit = n
	}
}
