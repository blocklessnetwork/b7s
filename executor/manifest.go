package executor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// TODO: Check - this functionality was ported but looks pretty special cased. Is this a temporary workaround for something?
// Investigate, then make proper.
func (e *Executor) writeFunctionManifest(req execute.Request, paths requestPaths) error {

	manifest := struct {
		FSRootPath    string   `json:"fs_root_path,omitempty"`
		Entry         string   `json:"entry,omitempty"`
		LimitedFuel   int      `json:"limited_fuel,omitempty"`
		LimitedMemory int      `json:"limited_memory,omitempty"`
		Permissions   []string `json:"permissions,omitempty"`
	}{
		FSRootPath:    paths.fsRoot,
		Entry:         paths.entry,
		LimitedFuel:   100_000_000,
		LimitedMemory: 200,
		Permissions:   req.Config.Permissions,
	}

	// Serialize manifest.
	encoded, err := json.MarshalIndent(manifest, "", "\t")
	if err != nil {
		return fmt.Errorf("could not marshal function manifest: %w", err)
	}

	// Write manifest to disk.
	err = os.WriteFile(paths.manifest, encoded, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write manifest to disk: %w", err)
	}

	return nil
}
