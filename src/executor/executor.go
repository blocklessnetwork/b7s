package executor

import (
	"context"
	"os/exec"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// executes a shell command to execute a wasm file
func Execute(ctx context.Context) (models.ExecutorResponse, error) {
	var executorResponse models.ExecutorResponse
	requestId, _ := uuid.NewRandom()
	cmd := "echo \"hello world\""
	run := exec.Command("bash", "-c", cmd)

	run.Dir = "/tmp"
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
