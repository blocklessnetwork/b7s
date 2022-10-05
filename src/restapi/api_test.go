package restapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
)

func TestHandleInstallFunction(t *testing.T) {
	ctx := context.Background()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)

	// set test context and test appdb
	config := models.Config{}
	config.Node.WorkSpaceRoot = "/tmp/b7s_test"
	ctx = context.WithValue(ctx, "config", config)
	appDb := db.Get("/tmp/b7s_test/api_testdb")
	ctx = context.WithValue(ctx, "appDb", appDb)

	req = req.WithContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	installFunctionReq := models.RequestFunctionInstall{
		Uri: "https://bafybeiho3scwi3njueloobzhg7ndn7yjb5rkcaydvsoxmnhmu2adv6oxzq.ipfs.w3s.link/manifest.json",
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

	db.Close(appDb)
	if returned.Code != enums.ResponseCodeOk {
		t.Errorf("handler returned unexpected body: got %v want %v",
			returned.Code, enums.ResponseCodeOk)
	}
}
