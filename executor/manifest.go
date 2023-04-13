package executor

import (
	"encoding/json"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// writeExecutionManifest will write a predefined execution manifest to disk.
func (e *Executor) writeExecutionManifest(req execute.Request, paths requestPaths) error {

	manifest := struct {
		FSRootPath    string   `json:"fs_root_path,omitempty"`
		Entry         string   `json:"entry,omitempty"`
		Permissions   []string `json:"permissions,omitempty"`
	}{
		FSRootPath:    paths.fsRoot,
		Entry:         paths.entry,
		Permissions:   req.Config.Permissions,
	}

	// Serialize manifest.
	encoded, err := json.MarshalIndent(manifest, "", "\t")
	if err != nil {
		return fmt.Errorf("could not marshal function manifest: %w", err)
	}

	// Write manifest to disk.
	err = e.writeFile(paths.manifest, encoded)
	if err != nil {
		return fmt.Errorf("could not write manifest to disk: %w", err)
	}

	return nil
}

// writeFile is a helper function wrapping the three OS-level calls - os.Create, file Write() and file Close().
func (e *Executor) writeFile(name string, data []byte) error {

	f, err := e.cfg.FS.Create(name)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}
