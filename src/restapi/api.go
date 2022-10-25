package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/blocklessnetworking/b7s/src/controller"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
)

func handleRequestExecute(w http.ResponseWriter, r *http.Request) {
	// body decode
	request := models.RequestExecute{}
	json.NewDecoder(r.Body).Decode(&request)

	// execute the function
	out, err := controller.ExecuteFunction(r.Context(), request)

	if err != nil {
		response := models.ResponseExecute{
			Code: enums.ResponseCodeError,
			Id:   out.RequestId,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResponseExecute{
		Code:   out.Code,
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

	if request.Uri == "" && request.Cid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// install the function
	// err := controller.InstallFunction(r.Context(), request.Uri)
	controller.MsgInstallFunction(r.Context(), request)

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
