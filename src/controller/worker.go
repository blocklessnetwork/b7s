package controller

import (
	"context"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/executor"
	"github.com/blocklessnetworking/b7s/src/models"
)

func WorkerExecuteFunction(ctx context.Context, request models.RequestExecute) (models.ExecutorResponse, error) {
	functionManifest, err := IsFunctionInstalled(ctx, request.FunctionId)

	// return if the function isn't installed
	// maybe install it?
	if err != nil {
		out := models.ExecutorResponse{
			Code: enums.ResponseCodeNotFound,
		}
		return out, err
	}

	out, err := executor.Execute(ctx, request, functionManifest)

	if err != nil {
		return out, err
	}

	return out, nil
}
