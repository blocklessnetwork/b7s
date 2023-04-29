//go:build windows
// +build windows

package limits

import (
	"fmt"
	"strings"
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

	// errMoreData is returned by QueryInformationJobObject to notify us on memory needed to store the response.
	// Unfortunately there's no Go error defined for it.
	errMoreDataStr = "More data is available."
)

type jobObjectCPURateControlInformation struct {
	ControlFlags uint32
	CPURate      uint32
}

type jobObjectExtendedLimitInformation struct {
	BasicLimitInformation windows.JOBOBJECT_BASIC_LIMIT_INFORMATION
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

type jobObjectBasicProcessIdList struct {
	NumberOfAssignedProcesses uint32
	NumberOfProcessIDsInList  uint32
	ProcessIDList             [1]uintptr
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

func setMemLimit(h windows.Handle, memoryKB int64) error {

	info := &jobObjectExtendedLimitInformation{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_JOB_MEMORY,
		},
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

func getJobObjectPids(h windows.Handle) ([]int, error) {

	var info jobObjectBasicProcessIdList

	err := windows.QueryInformationJobObject(
		h,
		jobObjectBasicProcessIdListInformationClass,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
		nil,
	)
	if err == nil {
		if info.NumberOfProcessIDsInList == 1 {

			pids := []int{
				int(info.ProcessIDList[0]),
			}

			return pids, nil
		}

		return []int{}, nil
	}
	if err != nil && !errIsMoreData(err) {
		return nil, fmt.Errorf("could not list job object processes: %w", err)
	}

	bufSize := unsafe.Sizeof(info) + (unsafe.Sizeof(info.ProcessIDList[0]) * uintptr(info.NumberOfAssignedProcesses-1))
	buf := make([]byte, bufSize)

	err = windows.QueryInformationJobObject(
		h,
		jobObjectBasicProcessIdListInformationClass,
		uintptr(unsafe.Pointer(&buf[0])),
		uint32(len(buf)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not list job object processes: %w", err)
	}

	bufInfo := (*jobObjectBasicProcessIdList)(unsafe.Pointer(&buf[0]))

	// Some dark sorcery ported from the MS `Host Compute Service Shim` library.
	// => See `AllPids` method over at https://github.com/microsoft/hcsshim/blob/main/internal/winapi/jobobject.go#L101
	pidList := (*[(1 << 27) - 1]uintptr)(unsafe.Pointer(&bufInfo.ProcessIDList[0]))[:bufInfo.NumberOfProcessIDsInList:bufInfo.NumberOfProcessIDsInList]

	out := make([]int, 0, bufInfo.NumberOfProcessIDsInList)
	for _, pid := range pidList {
		out = append(out, int(pid))
	}

	return out, nil
}

func errIsMoreData(err error) bool {
	return strings.Contains(err.Error(), errMoreDataStr)
}
