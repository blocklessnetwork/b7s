package executor

import (
	"context"
	"encoding/json"
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

func prepExecutionManifest(ctx context.Context, requestId string, request models.RequestExecute, manifest models.FunctionManifest) (string, error) {
	config := ctx.Value("config").(models.Config)

	functionPath := filepath.Join(config.Node.WorkspaceRoot, manifest.Function.ID, request.Method)
	manifestPath := filepath.Join(config.Node.WorkspaceRoot, "t", requestId, "runtime-manifest.json")
	tempFS := filepath.Join(config.Node.WorkspaceRoot, "t", requestId, "fs")

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

// executes a shell command to execute a wasm file
func Execute(ctx context.Context, request models.RequestExecute, functionManifest models.FunctionManifest) (models.ExecutorResponse, error) {
	requestId, _ := uuid.NewRandom()
	config := ctx.Value("config").(models.Config)
	tempFSPath := filepath.Join(config.Node.WorkspaceRoot, "t", requestId.String(), "fs")
	os.MkdirAll(tempFSPath, os.ModePerm)

	// check to see if runtime is available
	err := queryRuntime(config.Node.RuntimePath)

	if err != nil {
		return models.ExecutorResponse{
			Code:      enums.ResponseCodeError,
			RequestId: requestId.String(),
		}, err
	}

	runtimeManifestPath, err := prepExecutionManifest(ctx, requestId.String(), request, functionManifest)

	var executorResponse models.ExecutorResponse

	// check to see if there is any input to pass to the runtime
	var input string = ""
	var envVars []models.RequestExecuteEnvVars = request.Config.EnvVars
	var envVarString string = ""
	var envVarKeys string = ""

	if len(envVars) > 0 {
		for _, envVar := range envVars {
			envVarString += envVar.Name + "=\"" + envVar.Value + "\" "
			envVarKeys += envVar.Name + ";"
		}
		envVarString = "env " + envVarString
		envVarKeys = envVarKeys[:len(envVarKeys)-1]
	}

	if request.Config.Stdin != nil {
		input = *request.Config.Stdin
	}

	cmd := "echo \"" + input + "\" | " + envVarString + " BLS_LIST_VARS=\"" + envVarKeys + "\" " + config.Node.RuntimePath + "/blockless-cli " + runtimeManifestPath
	run := exec.Command("bash", "-c", cmd)
	run.Dir = tempFSPath
	out, err := run.Output()

	if err != nil {

		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to execute request")

		return executorResponse, err
	}

	executionResponseMemStore := ctx.Value("executionResponseMemStore").(memstore.ReqRespStore)
	err = executionResponseMemStore.Set(requestId.String(), &models.MsgExecuteResponse{
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
		"requestId": requestId,
	}).Info("function executed")

	executorResponse = models.ExecutorResponse{
		RequestId: requestId.String(),
		Code:      enums.ResponseCodeOk,
		Result:    string(out),
	}

	return executorResponse, nil
}
