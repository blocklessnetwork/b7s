//go:build windows
// +build windows

package limits

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Limits struct {
	cfg Config

	jh windows.Handle
}

// New creates a new process resource limit with the given configuration.
func New(opts ...Option) (*Limits, error) {

	// Create job object to which executions will be assigned to.
	name, err := windows.UTF16PtrFromString(DefaultJobObjectName)
	if err != nil {
		return nil, fmt.Errorf("could not prepare job object name: %w", err)
	}

	h, err := windows.CreateJobObject(nil, name)
	if err != nil {
		return nil, fmt.Errorf("could not create job object: %w", err)
	}

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.CPUPercentage < 1.0 {
		err := setCPULimit(h, cfg.CPUPercentage)
		if err != nil {
			windows.CloseHandle(h)
			return nil, fmt.Errorf("could not set CPU limit: %w", err)
		}
	}

	if cfg.MemoryKB > 0 {
		err := setMemLimit(h, cfg.MemoryKB)
		if err != nil {
			windows.CloseHandle(h)
			return nil, fmt.Errorf("could not set memory limit: %w", err)
		}
	}

	l := Limits{
		cfg: cfg,
		jh:  h,
	}

	return &l, nil
}

// LimitProcess will set the resource limits for the process identified by the handle.
func (l *Limits) LimitProcess(proc windows.Handle) error {

	err := windows.AssignProcessToJobObject(l.jh, proc)
	if err != nil {
		return fmt.Errorf("could not assign job to job object: %w", err)
	}

	return nil
}

type jobObjectBasicProcessIdList struct {
	NumberOfAssignedProcesses uint32
	NumberOfProcessIDsInList  uint32
	ProcessIDList             [1]uintptr
}

func (l *Limits) ListProcesses() ([]int, error) {

	info := &jobObjectBasicProcessIdList{}

	err := windows.QueryInformationJobObject(
		l.jh,
		jobObjectBasicProcessIdListInformationClass,
		uintptr(unsafe.Pointer(info)),
		uint32(unsafe.Sizeof(*info)),
		nil,
	)
	if err == nil {
		if info.NumberOfProcessIDsInList == 1 {
			return []int{
				int(info.ProcessIDList[0]),
			}, nil
		}

		return []int{}, nil
	}

	// TODO: Check if the error here is `ERROR_MORE_DATA`.
	if err != nil {
		return nil, fmt.Errorf("could not list job object processes: %w", err)
	}

	bufSize := unsafe.Sizeof(info) + (unsafe.Sizeof(info.ProcessIDList[0]) * uintptr(info.NumberOfAssignedProcesses-1))
	buf := make([]byte, bufSize)

	err = windows.QueryInformationJobObject(
		l.jh,
		jobObjectBasicProcessIdListInformationClass,
		uintptr(unsafe.Pointer(&buf[0])),
		uint32(unsafe.Sizeof(len(buf))),
		nil,
	)

	bufInfo := (*jobObjectBasicProcessIdList)(unsafe.Pointer(&buf[0]))
	pids := make([]int, 0, bufInfo.NumberOfProcessIDsInList)
	for _, pid := range bufInfo.ProcessIDList {
		p := int(pid)
		pids = append(pids, p)
	}

	return pids, nil
}

// Close will close the limiter.
func (l *Limits) Close() error {
	return windows.CloseHandle(l.jh)
}
