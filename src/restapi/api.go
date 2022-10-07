package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/blocklessnetworking/b7s/src/controller"
	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/cockroachdb/pebble"
)

func isFunctionInstalled(ctx context.Context, functionId string) (models.FunctionManifest, error) {
	appDb := ctx.Value("appDb").(*pebble.DB)
	functionManifestString, err := db.Value(appDb, functionId)
	functionManifest := models.FunctionManifest{}

	json.Unmarshal([]byte(functionManifestString), &functionManifest)

	if err != nil {
		if err.Error() == "pebble: not found" {
			return functionManifest, errors.New("function not installed")
		} else {
			return functionManifest, err
		}
	}

	return functionManifest, nil
}

func handleRequestExecute(w http.ResponseWriter, r *http.Request) {
	// body decode
	request := models.RequestExecute{}
	json.NewDecoder(r.Body).Decode(&request)

	functionManifest, err := isFunctionInstalled(r.Context(), request.FunctionId)

	// return if the function isn't installed
	// maybe install it?
	if err != nil {

		response := models.ResponseExecute{
			Code: enums.ResponseCodeNotFound,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	controller.RollCall(r.Context())

	// execute the function
	out, err := controller.ExecuteFunction(r.Context(), request, functionManifest)

	if err != nil {

		response := models.ResponseExecute{
			Code: enums.ResponseCodeError,
			Id:   out.RequestId,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResponseExecute{
		Code:   enums.ResponseCodeOk,
		Id:     out.RequestId,
		Result: out.Result,
	}

	json.NewEncoder(w).Encode(response)
}

func handleInstallFunction(w http.ResponseWriter, r *http.Request) {
	// body decode
	request := models.RequestFunctionInstall{}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&request)

	if request.Uri == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// install the function
	// err := controller.InstallFunction(r.Context(), request.Uri)
	controller.MsgInstallFunction(r.Context(), request.Uri)

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	response := models.ResponseInstall{
		Code: enums.ResponseCodeOk,
	}

	json.NewEncoder(w).Encode(response)
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func handleGetExecuteResponse(w http.ResponseWriter, r *http.Request) {
	// body decode
	request := models.RequestFunctionResponse{}
	json.NewDecoder(r.Body).Decode(&request)

	// get the response
	response := controller.GetExecutionResponse(r.Context(), request.Id)
	json.NewEncoder(w).Encode(response)
}
