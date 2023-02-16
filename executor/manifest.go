package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// TODO: Check - this functionality was ported but looks pretty special cased. Is this a temporary workaround for something?
// Investigate, then make proper.
func (e *Executor) writeFunctionManifest(executionID string, req execute.Request, workdir string) (string, error) {

	fnpath := filepath.Join(e.workdir, req.FunctionID, req.Method)
	manifestPath := filepath.Join(e.workdir, "t", executionID, "runtime-manifest.json")

	// Create parent directory for manifest.
	parent := filepath.Dir(manifestPath)
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("could not create parent directory for manifest: %w", err)
	}

	manifest := struct {
		FSRootPath    string   `json:"fs_root_path,omitempty"`
		Entry         string   `json:"entry,omitempty"`
		LimitedFuel   int      `json:"limited_fuel,omitempty"`
		LimitedMemory int      `json:"limited_memory,omitempty"`
		Permissions   []string `json:"permissions,omitempty"`
	}{
		FSRootPath:    workdir,
		Entry:         fnpath,
		LimitedFuel:   100_000_000,
		LimitedMemory: 200,
		Permissions:   req.Config.Permissions,
	}

	// Serialize manifest.
	encoded, err := json.MarshalIndent(manifest, "", "\t")
	if err != nil {
		return "", fmt.Errorf("could not marshal function manifest: %w", err)
	}

	// Write manifest to disk.
	err = os.WriteFile(manifestPath, encoded, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("could not write manifest to disk: %w", err)
	}

	return manifestPath, nil
}
