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
	out := controller.ExecuteFunction(r.Context())

	response := models.ResponseExecute{
		Code:   enums.ResponseCodeOk,
		Type:   enums.ResponseExecute,
		Id:     "",
		Result: out,
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

	// install the function
	out := controller.InstallFunction(r.Context(), request.Uri)

	response := models.ResponseInstall{
		Code:   enums.ResponseCodeOk,
		Type:   enums.ResponseInstall,
		Id:     "",
		Result: out,
	}

	json.NewEncoder(w).Encode(response)
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
