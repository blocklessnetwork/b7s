package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blocklessnetworking/b7s/src/db"
	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
)

func TestJSONRepository_Get(t *testing.T) {
	// remove b7s_test folder before starting the test
	os.RemoveAll("/tmp/b7s_test")

	mockConfig := models.Config{
		Protocol: models.ConfigProtocol{
			Role: enums.RoleWorker,
		},
		Node: models.ConfigNode{
			WorkspaceRoot: "/tmp/b7s_tests",
		},
	}
	ctx := context.WithValue(context.Background(), "config", mockConfig)

	appDb := db.GetDb("/tmp/b7s_test")
	ctx = context.WithValue(ctx, "appDb", appDb)
	defer db.Close(ctx)

	mockManifest := models.FunctionManifest{
		Function: models.Function{
			ID: "testFunction",
		},
		Deployment: models.Deployment{
			Uri: "",
		},
		Runtime: models.Runtime{
			Url:      "",
			Checksum: "fb1f409f1044844020c0aed9d8fe670484ce4af98c3768a72516d62cbf6a3c02",
		},
	}

	// create a test server to simulate the function repository
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/manifest.json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockManifest)
		} else if r.URL.Path == "/testFunction.tar.gz" {
			// Open the testFunction.tar.gz file
			file, _ := os.Open("testdata/testFunction.tar.gz")
			defer file.Close()

			// Send the file to the client
			io.Copy(w, file)
		}
	}))

	mockManifest.Deployment.Uri = ts.URL + "/testFunction.tar.gz"
	mockManifest.Runtime.Url = ts.URL + "/testFunction.tar.gz"
	defer ts.Close()

	// create a new JSONRepository struct
	repo := JSONRepository{}

	mockManifestBytes, _ := json.Marshal(mockManifest)

	// insert mock manifest into the database
	db.Set(ctx, "test_key", string(mockManifestBytes))
	installMsg := &models.MsgInstallFunction{
		ManifestUrl: ts.URL + "/manifest.json",
	}

	// call the Get method with the test server URL
	manifest, err := repo.Get(ctx, *installMsg)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// check that the function ID is correct
	if manifest.Function.ID != "testFunction" {
		t.Errorf("Expected function ID to be 'testFunction', got %s", manifest.Function.ID)
	}

	// check that the deployment URI is correct
	if manifest.Deployment.Uri != mockManifest.Deployment.Uri {
		t.Errorf("Expected deployment URI to be %s, got %s", mockManifest.Deployment.Uri, manifest.Deployment.Uri)
	}
}
