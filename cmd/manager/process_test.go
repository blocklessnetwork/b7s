package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// startDummyB7s starts a dummy b7s process that does nothing but sleep for a specified duration.
func startDummyB7s(sleepDuration time.Duration) (*os.Process, error) {
	cmd := exec.Command("sleep", sleepDuration.String())
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

func TestCheckB7sRunning(t *testing.T) {
	// Start a dummy b7s process.
	dummyB7s, err := startDummyB7s(5 * time.Second)
	if err != nil {
		t.Fatalf("Failed to start dummy b7s process: %v", err)
	}
	defer dummyB7s.Kill()

	// Test CheckB7sRunning function.
	processInfo, err := CheckB7sRunning()
	if err != nil {
		t.Fatalf("Error checking b7s process: %v", err)
	}

	if processInfo == nil {
		t.Fatal("Expected to find the dummy b7s process, but it was not found")
	}

	if processInfo.Pid != dummyB7s.Pid {
		t.Fatalf("Expected process PID to be %d, but got %d", dummyB7s.Pid, processInfo.Pid)
	}
}
