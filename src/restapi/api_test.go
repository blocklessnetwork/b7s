package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/memstore"
	"github.com/blocklessnetworking/b7s/src/models"
)

func TestHandleRequestExecute(t *testing.T) {
	// setup
	req, err := http.NewRequest("POST", "/function/request", bytes.NewBuffer([]byte(`{"functionId": "test_function", "input": "test_input"}`)))
	if err != nil {
		t.Fatal(err)
	}

	mockConfig := models.Config{
		Protocol: models.ConfigProtocol{
			Role: enums.RoleWorker,
		},
		Node: models.ConfigNode{
			WorkspaceRoot: "/tmp/b7s_tests",
		},
	}

	appDb := db.GetDb("/tmp/b7s")

	// mock the context
	ctx := context.WithValue(req.Context(), "config", models.Config{})
	ctx = context.WithValue(ctx, "config", mockConfig)
	ctx = context.WithValue(ctx, "appDb", appDb)
	req = req.WithContext(ctx)
	defer db.Close(ctx)

	// mock the response writer
	rr := httptest.NewRecorder()

	// call the function
	handleRequestExecute(rr, req)

	// check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "{\"code\":\"500\",\"id\":\"\",\"result\":\"\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleRootRequest(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRootRequest)

	// Serve the request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestHandleGetExecuteResponse(t *testing.T) {
	// setup
	store := memstore.NewReqRespStore()
	requestID := "123"
	response := &models.MsgExecuteResponse{
		Code:   enums.ResponseCodeOk,
		Result: "hello world",
	}
	store.Set(requestID, response)

	// create a request with the requestID
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"id":"`+requestID+`"}`)))
	if err != nil {
		t.Fatal(err)
	}

	// create a response recorder
	rr := httptest.NewRecorder()

	// create a context with the store
	ctx := context.WithValue(req.Context(), "executionResponseMemStore", store)

	// attach the context to the request
	req = req.WithContext(ctx)

	// call the handler
	handleGetExecuteResponse(rr, req)

	// check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the response body
	expected := "{\"code\":\"200\",\"result\":\"hello world\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestHandleInstallFunction(t *testing.T) {
	mockInstallFunction := func(ctx context.Context, request models.RequestFunctionInstall) {
		// Perform any necessary checks or assertions here
	}

	ctx := context.WithValue(context.Background(), "msgInstallFunc", mockInstallFunction)

	// Create a new function that wraps the handleInstallFunction
	// and calls it with the desired context
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		handleInstallFunction(w, r.WithContext(ctx))
	}

	// Create a test server to handle the request
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	// Prepare the request body
	request := models.RequestFunctionInstall{
		Uri: "https://example.com/manifest.json",
	}
	requestBytes, _ := json.Marshal(request)

	// Send the request to the test server
	resp, err := http.Post(testServer.URL, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	var response models.ResponseInstall
	json.NewDecoder(resp.Body).Decode(&response)
	if response.Code != enums.ResponseCodeOk {
		t.Errorf("Expected response code %v, got %v", enums.ResponseCodeOk, response.Code)
	}
}
