// Executor provides functions for running and interacting with functions.
package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// prepExecutionManifest creates the execution manifest file for the specified function.
func prepExecutionManifest(ctx context.Context, requestId string, request models.RequestExecute, manifest models.FunctionManifest) (string, error) {
	config := ctx.Value("config").(models.Config)

	functionPath := filepath.Join(config.Node.WorkspaceRoot, request.FunctionId, request.Method)
	manifestPath := filepath.Join(config.Node.WorkspaceRoot, "t", requestId, "runtime-manifest.json")
	tempFS := filepath.Join(config.Node.WorkspaceRoot, "t", requestId, "fs")

	// Create the directory
	os.MkdirAll(filepath.Dir(manifestPath), os.ModePerm)

	type Manifest struct {
		FS_ROOT_PATH   string   `json:"fs_root_path,omitempty"`
		ENTRY          string   `json:"entry,omitempty"`
		LIMITED_FUEL   int      `json:"limited_fuel,omitempty"`
		LIMITED_MEMORY int      `json:"limited_memory,omitempty"`
		PERMISSIONS    []string `json:"permissions,omitempty"`
	}

	data := Manifest{
		FS_ROOT_PATH:   tempFS,
		ENTRY:          functionPath,
		LIMITED_FUEL:   100000000,
		LIMITED_MEMORY: 200,
		PERMISSIONS:    request.Config.Permissions,
	}

	file, jsonError := json.MarshalIndent(data, "", " ")

	if jsonError != nil {
		log.WithFields(log.Fields{
			"err": jsonError,
		}).Warn("failed to marshal manifest")
		return "", jsonError
	}

	_ = ioutil.WriteFile(manifestPath, file, 0644)
	return manifestPath, nil
}

func queryRuntime(runtimePath string) error {
	cmd := exec.Command(runtimePath + "/blockless-cli")
	_, err := cmd.Output()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func Execute(ctx context.Context, request models.RequestExecute, functionManifest models.FunctionManifest) (models.ExecutorResponse, error) {
	requestID, _ := uuid.NewRandom()
	config := ctx.Value("config").(models.Config)
	tempFSPath := filepath.Join(config.Node.WorkspaceRoot, "t", requestID.String(), "fs")
	os.MkdirAll(tempFSPath, os.ModePerm)

	var execCommand func(string, ...string) *exec.Cmd
	if ctxExecCommand, ok := ctx.Value("execCommand").(func(string, ...string) *exec.Cmd); ok {
		execCommand = ctxExecCommand
	} else {
		// Check if the runtime is available.
		if err := queryRuntime(config.Node.RuntimePath); err != nil {
			return models.ExecutorResponse{
				Code:      enums.ResponseCodeError,
				RequestId: requestID.String(),
				Result:    "Runtime not available",
			}, err
		}
		execCommand = exec.Command
	}

	// Prepare the execution manifest.
	runtimeManifestPath, err := prepExecutionManifest(ctx, requestID.String(), request, functionManifest)
	if err != nil {
		return models.ExecutorResponse{
			Code:      enums.ResponseCodeError,
			RequestId: requestID.String(),
		}, err
	}

	// Build the input and environment variable strings.
	input := ""
	if request.Config.Stdin != nil {
		input = *request.Config.Stdin
	}
	envVars := request.Config.EnvVars
	envVarString, envVarKeys := "", ""
	if len(envVars) > 0 {
		for _, envVar := range envVars {
			envVarString += envVar.Name + "=\"" + envVar.Value + "\" "
			envVarKeys += envVar.Name + ";"
		}
		envVarString = "env " + envVarString
		envVarKeys = envVarKeys[:len(envVarKeys)-1]
	}

	// Build the command string.
	cmd := fmt.Sprintf("echo \"%s\" | %s BLS_LIST_VARS=\"%s\" %s/blockless-cli %s", input, envVarString, envVarKeys, config.Node.RuntimePath, runtimeManifestPath)

	// Execute the command.
	run := execCommand("bash", "-c", cmd)
	run.Dir = tempFSPath
	out, err := run.Output()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("failed to execute request")
		return models.ExecutorResponse{
			Code:      enums.ResponseCodeError,
			RequestId: requestID.String(),
		}, err
	}

	// Store the result in the execution response memory store

	executionResponseMemStore := ctx.Value("executionResponseMemStore").(memstore.ReqRespStore)
	err = executionResponseMemStore.Set(requestID.String(), &models.MsgExecuteResponse{
		Type:   enums.MsgExecuteResponse,
		Code:   enums.ResponseCodeOk,
		Result: string(out),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to set execution response")
	}

	log.WithFields(log.Fields{
		"requestId": requestID,
	}).Info("function executed")

	executorResponse := models.ExecutorResponse{
		RequestId: requestID.String(),
		Code:      enums.ResponseCodeOk,
		Result:    string(out),
	}

	return executorResponse, nil
}
