package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

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

type MsgInstallFunctionFunc func(context.Context, models.RequestFunctionInstall) error

func handleInstallFunction(w http.ResponseWriter, r *http.Request) {

	// Make sure that the request body is there.
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal request.
	var request models.RequestFunctionInstall
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Could be done using validators.
	if request.Uri == "" && request.Cid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Initialize the msgInstallFunction function - get the value from the context if set,
	// else use the default one.
	var msgInstallFunc MsgInstallFunctionFunc = controller.MsgInstallFunction

	// NOTE: At the moment, this function is no longer set on the context (for tests only).
	val := r.Context().Value("msgInstallFunc")
	if val != nil {
		// Assert that the context value is of the expected type.
		fn, ok := val.(MsgInstallFunctionFunc)
		if !ok {
			// Should never happen.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		msgInstallFunc = fn
	}

	// Add a deadline to the context.
	ctx, cancel := context.WithTimeout(r.Context(), functionInstallTimeout)
	defer cancel()

	// Start function install in a separate goroutine and signal when it's done.
	fnErr := make(chan error)
	go func() {
		err = msgInstallFunc(ctx, request)
		fnErr <- err
	}()

	// Wait until either function install finishes, or request times out.
	select {

	// Context timed out.
	case <-ctx.Done():

		status := http.StatusRequestTimeout
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		return

	// Work done.
	case err = <-fnErr:
		break
	}

	// Check if function install succeeded and handle error or return response.
	if err != nil {

		log.WithError(err).
			WithField("uri", request.Uri).
			WithField("cid", request.Cid).
			Error("failed to install function")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := models.ResponseInstall{
		Code: enums.ResponseCodeOk,
	}

	// Write response.
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
