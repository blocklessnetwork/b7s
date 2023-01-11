package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/models"
)

func TestIsFunctionInstalled(t *testing.T) {

	// set up a mock function manifest to store in the database
	mockManifest := models.FunctionManifest{
		Function: models.Function{
			ID:      "test-function",
			Name:    "Test Function",
			Version: "1.0.0",
			Runtime: "go",
		},
		Deployment: models.Deployment{
			Cid:      "Qmabcdef",
			Checksum: "123456789",
			Uri:      "https://ipfs.io/ipfs/Qmabcdef",
			Methods: []models.Methods{
				{
					Name:  "TestMethod",
					Entry: "main.TestMethod",
				},
			},
		},
		Runtime: models.Runtime{
			Checksum: "987654321",
			Url:      "https://ipfs.io/ipfs/Qmzyxwvu",
		},
	}
	mockManifestBytes, _ := json.Marshal(mockManifest)

	appDb := db.GetDb("/tmp/b7s")
	defer db.Close(appDb)
	ctx := context.WithValue(context.Background(), "appDb", appDb)

	// Insert a test value into the database
	db.Set(ctx, "test_key", string(mockManifestBytes))

	// Call IsFunctionInstalled
	functionManifest, err := IsFunctionInstalled(ctx, "test_key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Compare the function manifests as strings to account for potential encoding issues
	functionManifestBytes, _ := json.Marshal(functionManifest)
	expectedManifestBytes, _ := json.Marshal(mockManifest)
	if string(functionManifestBytes) != string(expectedManifestBytes) {
		t.Errorf("Unexpected function manifest. Got %v, expected %v", functionManifest, mockManifest)
	}
}
func TestExecuteFunction(t *testing.T) {
	// Create a mock Config value to pass to the context
	mockConfig := models.Config{
		Protocol: models.ConfigProtocol{
			Role: enums.RoleWorker,
		},
		Node: models.ConfigNode{
			WorkspaceRoot: "/tmp/b7s_tests",
		},
	}
	ctx := context.WithValue(context.Background(), "config", mockConfig)
	testStringValue := "foo"
	testString := fmt.Sprintf("echo %s", testStringValue)
	// Inject a mock execCommand function
	mockExecCommand := func(command string, args ...string) *exec.Cmd {
		cs := []string{"-c", testString}
		cmd := exec.Command("bash", cs...)
		return cmd
	}
	ctx = context.WithValue(ctx, "execCommand", mockExecCommand)

	// Set up a mock function manifest to store in the database
	mockManifest := models.FunctionManifest{
		Function: models.Function{
			ID:      "test-function",
			Name:    "Test Function",
			Version: "1.0.0",
			Runtime: "go",
		},
		Deployment: models.Deployment{
			Cid:      "Qmabcdef",
			Checksum: "123456789",
			Uri:      "https://ipfs.io/ipfs/Qmabcdef",
			Methods: []models.Methods{
				{
					Name:  "TestMethod",
					Entry: "main.TestMethod",
				},
			},
		},
		Runtime: models.Runtime{
			Checksum: "987654321",
			Url:      "https://ipfs.io/ipfs/Qmzyxwvu",
		},
	}
	mockManifestBytes, _ := json.Marshal(mockManifest)

	appDb := db.GetDb("/tmp/b7s")
	defer db.Close(appDb)
	ctx = context.WithValue(ctx, "appDb", appDb)

	// response memstore
	executionResponseMemStore := memstore.NewReqRespStore()
	ctx = context.WithValue(ctx, "executionResponseMemStore", executionResponseMemStore)

	// Insert the mock function manifest into the database
	db.Set(ctx, "test-function", string(mockManifestBytes))

	// Create a mock RequestExecute value to pass to ExecuteFunction
	mockRequest := models.RequestExecute{
		FunctionId: "test-function",
		Method:     "TestMethod",
		Parameters: []models.RequestExecuteParameters{
			{
				Name:  "param1",
				Value: "value1",
			},
		},
		Config: models.ExecutionRequestConfig{},
	}

	// Call ExecuteFunction with the mock context and request
	response, err := ExecuteFunction(ctx, mockRequest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the correct function was called (Worker`Exec`uteFunction in this case)
	if strings.Trim(response.Result, "\n") != testStringValue {
		t.Errorf("Unexpected response. Got %v, expected %v", response.Result, testStringValue)
	}
}
