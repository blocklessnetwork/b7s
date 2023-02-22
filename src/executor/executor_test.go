package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/google/uuid"
)

func TestPrepExecutionManifest(t *testing.T) {
	// setup
	config := models.Config{
		Node: models.ConfigNode{
			WorkspaceRoot: "/tmp/workspace",
		},
	}
	requestID, _ := uuid.NewRandom()
	request := models.RequestExecute{
		Method: "testMethod",
		Config: models.ExecutionRequestConfig{
			Permissions: []string{"permission1", "permission2"},
		},
	}
	functionManifest := models.FunctionManifest{
		Function: models.Function{
			ID: "testFunction",
		},
	}
	ctx := context.WithValue(context.Background(), "config", config)

	// run the test
	manifestPath, err := prepExecutionManifest(ctx, requestID.String(), request, functionManifest)

	// test the result
	if err != nil {
		t.Errorf("prepExecutionManifest returned an error: %v", err)
	}
	if manifestPath != filepath.Join(config.Node.WorkspaceRoot, "t", requestID.String(), "runtime-manifest.json") {
		t.Errorf("Unexpected manifest path, got %s", manifestPath)
	}

	// check that the file exists and can be read
	file, fileErr := os.Open(manifestPath)
	if fileErr != nil {
		t.Errorf("Error opening the manifest file: %v", fileErr)
	}
	defer file.Close()

	// cleanup
	os.RemoveAll(filepath.Join(config.Node.WorkspaceRoot, "t", requestID.String()))
}
