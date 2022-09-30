package rest

import (
	"encoding/json"
	"net/http"

	"github.com/blocklessnetworking/b7s/src/controller"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/gorilla/mux"
)

func handleRequestExecute(w http.ResponseWriter, r *http.Request) {
	// params
	vars := mux.Vars(r)
	id := vars["id"]

	// execute the function
	out := controller.ExecuteFunction(r.Context())

	response := models.ResponseExecute{
		Code:   enums.ResponseCodeOk,
		Type:   enums.ResponseExecute,
		Id:     id,
		Result: out,
	}

	json.NewEncoder(w).Encode(response)
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
