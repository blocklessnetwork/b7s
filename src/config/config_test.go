package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a test config.yaml file
	configYaml := []byte(`
node:
  workspace_root: "./testdata"
  `)
	ioutil.WriteFile("/tmp/b7_test_config.yaml", configYaml, 0644)
	defer os.Remove("/tmp/b7_test_config.yaml")

	// Load config
	err := Load("/tmp/b7_test_config.yaml")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	// Check that config values are correct
	if C.Node.WorkspaceRoot != "./testdata" {
		t.Errorf("Expected workspace root to be './testdata', got %s", C.Node.WorkspaceRoot)
	}
}
