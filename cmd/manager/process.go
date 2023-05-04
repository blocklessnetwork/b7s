package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type ProcessInfo struct {
	Pid      int
	User     string
	Cmdline  string
}

func CheckB7sRunning() (*ProcessInfo, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("pgrep", "-fl", "b7s")
	case "windows":
		cmd = exec.Command("tasklist", "/FI", "imagename eq b7s.exe")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outputStr := strings.TrimSpace(string(output))
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
