//go:build windows
// +build windows

package limits

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	jobObjectBasicProcessIdListInformationClass = 3
)

const (
	JOB_OBJECT_CPU_RATE_CONTROL_ENABLE = 0x1
	// The job's CPU rate is a hard limit. After the job reaches its CPU cycle limit for the current scheduling interval,
	// no threads associated with the job will run until the next interval.
	// => See https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_cpu_rate_control_information
	JOB_OBJECT_CPU_RATE_CONTROL_HARD_CAP = 0x4
)

type jobObjectCPURateControlInformation struct {
	ControlFlags uint32
	CPURate      uint32
}

func setCPULimit(h windows.Handle, cpuRate float64) error {

	// Specifies the portion of processor cycles that the threads in a job object can use during each scheduling interval, as the number of cycles per 10,000 cycles.
	// Set CpuRate to a percentage times 100. For example, to let the job use 20% of the CPU, set CpuRate to 20 times 100, or 2,000.
	// => See https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_cpu_rate_control_information

	info := &jobObjectCPURateControlInformation{
		ControlFlags: JOB_OBJECT_CPU_RATE_CONTROL_ENABLE | JOB_OBJECT_CPU_RATE_CONTROL_HARD_CAP,

		// Convert rate from e.g. 0.8 to 80(%) * 100.
		CPURate: uint32((100 * cpuRate) * 100),
	}

	_, err := windows.SetInformationJobObject(
		h,
		windows.JobObjectCpuRateControlInformation,
		uintptr(unsafe.Pointer(info)),
		uint32(unsafe.Sizeof(*info)),
	)
	if err != nil {
		return fmt.Errorf("could not set CPU limit for job: %w", err)
	}

	return nil
}

type jobObjectExtendedLimitInformation struct {
	BasicLimitInformation windows.MemoryBasicInformation
	IoInfo                ioCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

type ioCounters struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

func setMemLimit(h windows.Handle, memoryKB int64) error {

	info := &jobObjectExtendedLimitInformation{
		JobMemoryLimit: uintptr(memoryKB * 1000),
	}

	_, err := windows.SetInformationJobObject(
		h,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(info)),
		uint32(unsafe.Sizeof(*info)),
	)
	if err != nil {
		return fmt.Errorf("could not set memory limit for job: %w", err)
	}

	return nil
}
