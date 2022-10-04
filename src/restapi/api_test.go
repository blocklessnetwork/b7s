package restapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
)

func TestHandleInstallFunction(t *testing.T) {
	ctx := context.Background()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)

	config := models.Config{}
	config.Node.WorkSpaceRoot = "/tmp/b7s_test"

	ctx = context.WithValue(ctx, "config", config)

	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	installFunctionReq := models.RequestFunctionInstall{
		Uri: "https://bafybeibyniiukxqmb7qae7ljif6atvo7ipg6wnpwvtqb4stf4ubjjterha.ipfs.w3s.link/manifest.json",
	}

	data, _ := json.Marshal(installFunctionReq)
	stringReader := strings.NewReader(string(data))
	stringReadCloser := io.NopCloser(stringReader)

	req.Body = stringReadCloser

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleInstallFunction)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	returned := models.ResponseInstall{}
	json.Unmarshal(rr.Body.Bytes(), &returned)

	if returned.Code != enums.ResponseCodeOk {
		t.Errorf("handler returned unexpected body: got %v want %v",
			returned.Code, enums.ResponseCodeOk)
	}
}
