package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	const (
		contentTypeJSON = "application/json"
		requestPayload  = `{ "uri": "https://example.com/manifest.json" }`
	)

	var installFn MsgInstallFunctionFunc = func(ctx context.Context, req models.RequestFunctionInstall) error {
		return nil
	}

	ctx := context.WithValue(context.Background(), "msgInstallFunc", installFn)

	// Function that processes the HTTP request.
	handler := func(w http.ResponseWriter, r *http.Request) {
		handleInstallFunction(w, r.WithContext(ctx))
	}

	// Create a test server to handle the request.
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	// Send the request to the test server.
	res, err := http.Post(srv.URL, contentTypeJSON, strings.NewReader(requestPayload))
	if err != nil {
		t.Fatalf("could not execute POST request: %s", err)
	}

	// Check the response status code.
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code (want: %v, got %v)", http.StatusOK, res.StatusCode)
	}

	// Unpack the response body.
	var response models.ResponseInstall
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("could not decode server response: %s", err)
	}

	// Verify the response
	if response.Code != enums.ResponseCodeOk {
		t.Errorf("unexpected response code (want: %v, got %v)", enums.ResponseCodeOk, response.Code)
	}
}

func TestHandleInstallFunction_HandlesErrors(t *testing.T) {

	const (
		contentTypeJSON = "application/json"

		malformedJSON = `{ "uri": "https://example.com/manifest.json" ` // JSON with missing closing brace.
		emptyRequest  = `{ "uri": "", "cid": ""} `                      // Both URI and CID are empty.
		validJSON     = `{ "uri": "https://example.com/manifest.json" }`
	)

	var (
		installFn MsgInstallFunctionFunc = func(context.Context, models.RequestFunctionInstall) error {
			return nil
		}
		failingInstallFn MsgInstallFunctionFunc = func(context.Context, models.RequestFunctionInstall) error {
			return errors.New("stop")
		}
	)

	tests := []struct {
		name string

		payload   io.Reader
		installFn MsgInstallFunctionFunc

		expectedCode int
	}{
		{
			name:         "missing response body",
			payload:      nil,
			installFn:    installFn,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "malformed JSON payload",
			payload:      strings.NewReader(malformedJSON),
			installFn:    installFn,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing URI and CID",
			payload:      strings.NewReader(emptyRequest),
			installFn:    installFn,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "install function failed",
			payload:      strings.NewReader(validJSON),
			installFn:    failingInstallFn,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {

		test := test

		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			ctx := context.WithValue(context.Background(), "msgInstallFunc", test.installFn)

			// Function that processes the HTTP request.
			handler := func(w http.ResponseWriter, r *http.Request) {
				handleInstallFunction(w, r.WithContext(ctx))
			}

			// Create a test server to handle the request.
			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()

			// Send the request to the test server.
			res, err := http.Post(srv.URL, contentTypeJSON, test.payload)
			if err != nil {
				t.Fatalf("could not execute POST request: %s", err)
			}

			// Check the response status code.
			if res.StatusCode != test.expectedCode {
				t.Fatalf("unexpected status code (want: %v, got %v)", test.expectedCode, res.StatusCode)
			}
		})
	}
}
