package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type ProcessInfo struct {
	Pid     int
	User    string
	Cmdline string
}

func getProcessCommand(processName string) *exec.Cmd {
	switch runtime.GOOS {
	case "linux", "darwin":
		return exec.Command("pgrep", "-fl", processName)
	case "windows":
		return exec.Command("tasklist", "/FI", fmt.Sprintf("imagename eq %s.exe", processName))
	default:
		return nil
	}
}

func parseProcessOutput(outputStr string) (*ProcessInfo, error) {
	if outputStr == "" {
		return nil, nil
	}

	var info ProcessInfo
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Sscanf(outputStr, "%d %s", &info.Pid, &info.Cmdline)
	case "windows":
		lines := strings.Split(outputStr, "\n")
		if len(lines) > 2 {
			fmt.Sscanf(lines[2], "b7s.exe %d", &info.Pid)
			info.Cmdline = "b7s.exe"
		}
	}

	return &info, nil
}

func CheckB7sRunning(processName string) (*ProcessInfo, error) {
	cmd := getProcessCommand(processName)
	if cmd == nil {
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outputStr := strings.TrimSpace(string(output))
	return parseProcessOutput(outputStr)
}
